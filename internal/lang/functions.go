package lang

import (
	"image"
	"reflect"
	"sort"
)

type functionDecl struct {
	body   func(ir *interpreter, values []value) (value, error)
	params []reflect.Type
}

var numberType = reflect.TypeOf(Number(0))
var positionType = reflect.TypeOf(Position{})
var kernelType = reflect.TypeOf(kernel{})

var functions = map[string]functionDecl{
	"rgb": functionDecl{
		body:   invokeRgb,
		params: []reflect.Type{numberType, numberType, numberType},
	},
	"srgb": functionDecl{
		body:   invokeSrgb,
		params: []reflect.Type{numberType, numberType, numberType},
	},
	"rgba": functionDecl{
		body:   invokeRgba,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"srgba": functionDecl{
		body:   invokeSrgba,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"rect": functionDecl{
		body:   invokeRect,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"convolute": functionDecl{
		body:   invokeConvolute,
		params: []reflect.Type{positionType, kernelType},
	},
	"sort": functionDecl{
		body:   invokeSort,
		params: []reflect.Type{kernelType},
	},
	"map_r": functionDecl{
		body:   invokeMapR,
		params: []reflect.Type{positionType, kernelType},
	},
	"map_g": functionDecl{
		body:   invokeMapG,
		params: []reflect.Type{positionType, kernelType},
	},
	"map_b": functionDecl{
		body:   invokeMapB,
		params: []reflect.Type{positionType, kernelType},
	},
	"map_a": functionDecl{
		body:   invokeMapA,
		params: []reflect.Type{positionType, kernelType},
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

func invokeSort(ir *interpreter, params []value) (value, error) {
	kernelVal := params[0].(kernel)
	result := kernelVal
	result.values = append([]Number(nil), result.values...) // clone values
	sort.Sort(numberSlice(result.values))
	return result, nil
}

func invokeMapR(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values), nil
}

func invokeMapG(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values), nil
}

func invokeMapB(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values), nil
}

func invokeMapA(ir *interpreter, params []value) (value, error) {
	posVal := params[0].(Position)
	kernelVal := params[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.radius, kernelVal.width, kernelVal.values), nil
}
