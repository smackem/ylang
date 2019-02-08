package lang

import (
	"image"
	"math"
	"reflect"
	"sort"
)

type functionDecl struct {
	body   func(ir *interpreter, values []value) (value, error)
	params []reflect.Type
}

var numberType = reflect.TypeOf(Number(0))
var kernelType = reflect.TypeOf(kernel{})
var positionType = reflect.TypeOf(Position{})

var functions = map[string]functionDecl{
	"rgb": {
		body:   invokeRgb,
		params: []reflect.Type{numberType, numberType, numberType},
	},
	"srgb": {
		body:   invokeSrgb,
		params: []reflect.Type{numberType, numberType, numberType},
	},
	"rgba": {
		body:   invokeRgba,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"srgba": {
		body:   invokeSrgba,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"rect": {
		body:   invokeRect,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"convolute": {
		body:   invokeConvolute,
		params: []reflect.Type{positionType, kernelType},
	},
	"sort": {
		body:   invokeSort,
		params: []reflect.Type{kernelType},
	},
	"map_r": {
		body:   invokeMapR,
		params: []reflect.Type{positionType, kernelType},
	},
	"map_g": {
		body:   invokeMapG,
		params: []reflect.Type{positionType, kernelType},
	},
	"map_b": {
		body:   invokeMapB,
		params: []reflect.Type{positionType, kernelType},
	},
	"map_a": {
		body:   invokeMapA,
		params: []reflect.Type{positionType, kernelType},
	},
	"sin": {
		body:   invokeSin,
		params: []reflect.Type{numberType},
	},
	"cos": {
		body:   invokeCos,
		params: []reflect.Type{numberType},
	},
	"tan": {
		body:   invokeTan,
		params: []reflect.Type{numberType},
	},
	"asin": {
		body:   invokeAsin,
		params: []reflect.Type{numberType},
	},
	"acos": {
		body:   invokeAcos,
		params: []reflect.Type{numberType},
	},
	"atan": {
		body:   invokeAtan,
		params: []reflect.Type{numberType},
	},
	"atan2": {
		body:   invokeAtan2,
		params: []reflect.Type{numberType, numberType},
	},
	"sqrt": {
		body:   invokeSqrt,
		params: []reflect.Type{numberType},
	},
	"abs": {
		body:   invokeAbs,
		params: []reflect.Type{numberType},
	},
	"min": {
		body:   invokeMin,
		params: []reflect.Type{kernelType},
	},
	"max": {
		body:   invokeMax,
		params: []reflect.Type{kernelType},
	},
	"list": {
		body:   invokeList,
		params: []reflect.Type{numberType, numberType},
	},
	"kernel": {
		body:   invokeKernel,
		params: []reflect.Type{numberType, numberType},
	},
}

func invokeRgb(_ *interpreter, params []value) (value, error) {
	return NewRgba(params[0].(Number), params[1].(Number), params[2].(Number), 255), nil
}

func invokeSrgb(_ *interpreter, params []value) (value, error) {
	return NewSrgba(params[0].(Number), params[1].(Number), params[2].(Number), 1.0), nil
}

func invokeRgba(_ *interpreter, params []value) (value, error) {
	return NewRgba(params[0].(Number), params[1].(Number), params[2].(Number), params[3].(Number)), nil
}

func invokeSrgba(_ *interpreter, params []value) (value, error) {
	return NewSrgba(params[0].(Number), params[1].(Number), params[2].(Number), params[3].(Number)), nil
}

func invokeRect(_ *interpreter, params []value) (value, error) {
	x, y := int(params[0].(Number)), int(params[1].(Number))
	return Rect{
		Min: image.Point{x, y},
		Max: image.Point{x + int(params[2].(Number)+0.5), y + int(params[3].(Number)+0.5)},
	}, nil
}

func invokeConvolute(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values), nil
}

// implement sort.Interface for Number slice

type numberSlice []Number

func (p numberSlice) Len() int           { return len(p) }
func (p numberSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p numberSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func invokeSort(_ *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	result := kernelVal
	result.values = append([]Number(nil), result.values...) // clone values
	sort.Sort(numberSlice(result.values))
	return result, nil
}

func invokeMapR(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapRed(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values)
	return result, nil
}

func invokeMapG(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapGreen(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values)
	return result, nil
}

func invokeMapB(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapBlue(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values)
	return result, nil
}

func invokeMapA(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapAlpha(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values)
	return result, nil
}

func invokeSin(_ *interpreter, params []value) (value, error) {
	return Number(math.Sin(float64(params[0].(Number)))), nil
}

func invokeCos(_ *interpreter, params []value) (value, error) {
	return Number(math.Cos(float64(params[0].(Number)))), nil
}

func invokeTan(_ *interpreter, params []value) (value, error) {
	return Number(math.Tan(float64(params[0].(Number)))), nil
}

func invokeAsin(_ *interpreter, params []value) (value, error) {
	return Number(math.Asin(float64(params[0].(Number)))), nil
}

func invokeAcos(_ *interpreter, params []value) (value, error) {
	return Number(math.Acos(float64(params[0].(Number)))), nil
}

func invokeAtan(_ *interpreter, params []value) (value, error) {
	return Number(math.Atan(float64(params[0].(Number)))), nil
}

func invokeAtan2(_ *interpreter, params []value) (value, error) {
	return Number(math.Atan2(float64(params[0].(Number)), float64(params[1].(Number)))), nil
}

func invokeSqrt(_ *interpreter, params []value) (value, error) {
	return Number(math.Sqrt(float64(params[0].(Number)))), nil
}

func invokeAbs(_ *interpreter, params []value) (value, error) {
	return Number(math.Abs(float64(params[0].(Number)))), nil
}

func invokeMax(_ *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	max := Number(math.MinInt32)
	for _, n := range kernelVal.values {
		if n > max {
			max = n
		}
	}
	return Number(max), nil
}

func invokeMin(_ *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	min := Number(math.MaxInt32)
	for _, n := range kernelVal.values {
		if n < min {
			min = n
		}
	}
	return Number(min), nil
}

func invokeList(_ *interpreter, params []value) (value, error) {
	count := params[0].(Number)
	val := params[1].(Number)
	result := kernel{values: make([]Number, int(count))}
	for i := range result.values {
		result.values[i] = val
	}
	return result, nil
}

func invokeKernel(_ *interpreter, params []value) (value, error) {
	width := int(params[0].(Number))
	height := int(params[1].(Number))
	count := width * height
	result := kernel{width: width, height: height, values: make([]Number, count)}
	return result, nil
}
