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
var positionType = reflect.TypeOf(point{})
var kernelType = reflect.TypeOf(kernel{})

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
		params: []reflect.Type{numberType, numberType, numberType},
	},
	"resize": {
		body:   invokeResize,
		params: []reflect.Type{numberType, numberType},
	},
	"line": {
		body:   invokeLine,
		params: []reflect.Type{positionType, positionType},
	},
	"polygon": {
		body:   invokePolygon,
		params: []reflect.Type{reflect.TypeOf([]point{})},
	},
}

func invokeRgb(ir *interpreter, params []value) (value, error) {
	return NewRgba(params[0].(Number), params[1].(Number), params[2].(Number), 255), nil
}

func invokeSrgb(ir *interpreter, params []value) (value, error) {
	return NewSrgba(params[0].(Number), params[1].(Number), params[2].(Number), 1.0), nil
}

func invokeRgba(ir *interpreter, params []value) (value, error) {
	return NewRgba(params[0].(Number), params[1].(Number), params[2].(Number), params[3].(Number)), nil
}

func invokeSrgba(ir *interpreter, params []value) (value, error) {
	return NewSrgba(params[0].(Number), params[1].(Number), params[2].(Number), params[3].(Number)), nil
}

func invokeRect(ir *interpreter, params []value) (value, error) {
	x, y := int(params[0].(Number)), int(params[1].(Number))
	return rect{
		Min: image.Point{x, y},
		Max: image.Point{x + int(params[2].(Number)+0.5), y + int(params[3].(Number)+0.5)},
	}, nil
}

func invokeConvolute(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(point)
	kernelVal := params[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values), nil
}

// implement sort.Interface for Number slice

type numberSlice []Number

func (p numberSlice) Len() int           { return len(p) }
func (p numberSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p numberSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func invokeSort(ir *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	result := kernelVal
	result.values = append([]Number(nil), result.values...) // clone values
	sort.Sort(numberSlice(result.values))
	return result, nil
}

func invokeMapR(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(point)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapRed(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeMapG(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(point)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapGreen(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeMapB(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(point)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapBlue(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeMapA(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(point)
	kernelVal := params[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapAlpha(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeSin(ir *interpreter, params []value) (value, error) {
	return Number(math.Sin(float64(params[0].(Number)))), nil
}

func invokeCos(ir *interpreter, params []value) (value, error) {
	return Number(math.Cos(float64(params[0].(Number)))), nil
}

func invokeTan(ir *interpreter, params []value) (value, error) {
	return Number(math.Tan(float64(params[0].(Number)))), nil
}

func invokeAsin(ir *interpreter, params []value) (value, error) {
	return Number(math.Asin(float64(params[0].(Number)))), nil
}

func invokeAcos(ir *interpreter, params []value) (value, error) {
	return Number(math.Acos(float64(params[0].(Number)))), nil
}

func invokeAtan(ir *interpreter, params []value) (value, error) {
	return Number(math.Atan(float64(params[0].(Number)))), nil
}

func invokeAtan2(ir *interpreter, params []value) (value, error) {
	return Number(math.Atan2(float64(params[0].(Number)), float64(params[1].(Number)))), nil
}

func invokeSqrt(ir *interpreter, params []value) (value, error) {
	return Number(math.Sqrt(float64(params[0].(Number)))), nil
}

func invokeAbs(ir *interpreter, params []value) (value, error) {
	return Number(math.Abs(float64(params[0].(Number)))), nil
}

func invokeMax(it *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	max := Number(math.MinInt32)
	for _, n := range kernelVal.values {
		if n > max {
			max = n
		}
	}
	return Number(max), nil
}

func invokeMin(it *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	min := Number(math.MaxInt32)
	for _, n := range kernelVal.values {
		if n < min {
			min = n
		}
	}
	return Number(min), nil
}

func invokeList(it *interpreter, params []value) (value, error) {
	count := params[0].(Number)
	val := params[1].(Number)
	values := make([]Number, int(count))
	for i := range values {
		values[i] = val
	}
	return kernel{values: values}, nil
}

func invokeKernel(ir *interpreter, params []value) (value, error) {
	width := params[0].(Number)
	height := params[1].(Number)
	val := params[2].(Number)
	values := make([]Number, int(width*height))
	for i := range values {
		values[i] = val
	}
	return kernel{width: int(width), height: int(height), values: values}, nil
}

func invokeResize(ir *interpreter, params []value) (value, error) {
	width := params[0].(Number)
	height := params[1].(Number)
	ir.bitmap.ResizeTarget(int(width), int(height))
	return rect{
		Max: image.Point{int(width), int(height)},
	}, nil
}

func invokeLine(ir *interpreter, params []value) (value, error) {
	point1, point2 := params[0].(point), params[1].(point)
	return line{point1, point2}, nil
}

func invokePolygon(ir *interpreter, params []value) (value, error) {
	return nil, nil
}
