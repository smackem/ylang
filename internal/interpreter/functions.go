package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"image"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strings"
)

type functionDecl struct {
	body   func(ir *interpreter, values []Value) (Value, error)
	params []reflect.Type
}

var numberType = reflect.TypeOf(number(0))
var pointType = reflect.TypeOf(point{})
var kernelType = reflect.TypeOf(kernel{})
var listType = reflect.TypeOf(list{})
var colorType = reflect.TypeOf(lang.Color{})
var functionType = reflect.TypeOf(function{})
var lineType = reflect.TypeOf(line{})
var rectType = reflect.TypeOf(rect{})
var polygonType = reflect.TypeOf(polygon{})
var numberSliceType = reflect.TypeOf([]number{})
var pointSliceType = reflect.TypeOf([]point{})
var circleType = reflect.TypeOf(circle{})
var hsvType = reflect.TypeOf(colorHsv{})
var valueType = reflect.TypeOf((*Value)(nil)).Elem()

var functions map[string][]functionDecl

func initFunctions() {
	if functions != nil {
		return
	}
	functions = map[string][]functionDecl{
		"rgb": {
			{
				body:   invokeRgb,
				params: []reflect.Type{numberType, numberType, numberType},
			},
			{
				body:   invokeHsv2Rgb,
				params: []reflect.Type{hsvType},
			},
			{
				body:   invokeGrey,
				params: []reflect.Type{numberType},
			},
		},
		"rgb01": {
			{
				body:   invokeSrgb,
				params: []reflect.Type{numberType, numberType, numberType},
			},
			{
				body:   invokeSgrey,
				params: []reflect.Type{numberType},
			},
		},
		"rgba": {
			{
				body:   invokeRgba,
				params: []reflect.Type{numberType, numberType, numberType, numberType},
			},
			{
				body:   invokeRgb2Rgba,
				params: []reflect.Type{colorType, numberType},
			},
			{
				body:   invokeHsv2Rgba,
				params: []reflect.Type{hsvType, numberType},
			},
		},
		"rgba01": {
			{
				body:   invokeSrgba,
				params: []reflect.Type{numberType, numberType, numberType, numberType},
			},
		},
		"rect": {
			{
				body:   invokeRect,
				params: []reflect.Type{numberType, numberType, numberType, numberType},
			},
		},
		"convolute": {
			{
				body:   invokeConvolute,
				params: []reflect.Type{pointType, kernelType},
			},
		},
		"blt": {
			{
				body:   invokeBlt,
				params: []reflect.Type{rectType},
			},
		},
		"flip": {
			{
				body:   invokeFlip,
				params: []reflect.Type{},
			},
		},
		"recall": {
			{
				body:   invokeRecall,
				params: []reflect.Type{numberType},
			},
		},
		"sort": {
			{
				body:   invokeSortKernel,
				params: []reflect.Type{kernelType},
			},
			{
				body:   invokeSortList,
				params: []reflect.Type{listType},
			},
			{
				body:   invokeSortListFn,
				params: []reflect.Type{listType, functionType},
			},
		},
		"fetchRed": {
			{
				body:   invokeFetchRed,
				params: []reflect.Type{pointType, kernelType},
			},
		},
		"fetchGreen": {
			{
				body:   invokeFetchGreen,
				params: []reflect.Type{pointType, kernelType},
			},
		},
		"fetchBlue": {
			{
				body:   invokeFetchBlue,
				params: []reflect.Type{pointType, kernelType},
			},
		},
		"fetchAlpha": {
			{
				body:   invokeFetchAlpha,
				params: []reflect.Type{pointType, kernelType},
			},
		},
		"sin": {
			{
				body:   invokeSin,
				params: []reflect.Type{numberType},
			},
		},
		"cos": {
			{
				body:   invokeCos,
				params: []reflect.Type{numberType},
			},
		},
		"tan": {
			{
				body:   invokeTan,
				params: []reflect.Type{numberType},
			},
		},
		"asin": {
			{
				body:   invokeAsin,
				params: []reflect.Type{numberType},
			},
		},
		"acos": {
			{
				body:   invokeAcos,
				params: []reflect.Type{numberType},
			},
		},
		"atan": {
			{
				body:   invokeAtan,
				params: []reflect.Type{numberType},
			},
		},
		"atan2": {
			{
				body:   invokeAtan2,
				params: []reflect.Type{numberType, numberType},
			},
		},
		"sqrt": {
			{
				body:   invokeSqrt,
				params: []reflect.Type{numberType},
			},
		},
		"pow": {
			{
				body:   invokePow,
				params: []reflect.Type{numberType, numberType},
			},
		},
		"abs": {
			{
				body:   invokeAbs,
				params: []reflect.Type{numberType},
			},
		},
		"round": {
			{
				body:   invokeRound,
				params: []reflect.Type{numberType},
			},
		},
		"floor": {
			{
				body:   invokeFloor,
				params: []reflect.Type{numberType},
			},
		},
		"ceil": {
			{
				body:   invokeCeil,
				params: []reflect.Type{numberType},
			},
		},
		"hypot": {
			{
				body:   invokeHypot,
				params: []reflect.Type{numberType, numberType},
			},
			{
				body:   invokeHypotRgb,
				params: []reflect.Type{colorType, colorType},
			},
			{
				body:   invokeHypotPoint,
				params: []reflect.Type{pointType},
			},
		},
		"random": {
			{
				body:   invokeRandom,
				params: []reflect.Type{numberType, numberType},
			},
			{
				body:   invokeRandom01,
				params: []reflect.Type{},
			},
		},
		"min": {
			{
				body:   invokeMin,
				params: []reflect.Type{numberSliceType},
			},
			{
				body:   invokeMinKernel,
				params: []reflect.Type{kernelType},
			},
			{
				body:   invokeMinList,
				params: []reflect.Type{listType},
			},
		},
		"max": {
			{
				body:   invokeMax,
				params: []reflect.Type{numberSliceType},
			},
			{
				body:   invokeMaxKernel,
				params: []reflect.Type{kernelType},
			},
			{
				body:   invokeMaxList,
				params: []reflect.Type{listType},
			},
		},
		"list": {
			{
				body:   invokeList,
				params: []reflect.Type{numberType, numberType},
			},
			{
				body:   invokeListFn,
				params: []reflect.Type{numberType, functionType},
			},
		},
		"kernel": {
			{
				body:   invokeKernel,
				params: []reflect.Type{numberType, numberType, numberType},
			},
			{
				body:   invokeKernelFn,
				params: []reflect.Type{numberType, numberType, functionType},
			},
		},
		"gauss": {
			{
				body:   invokeGauss,
				params: []reflect.Type{numberType},
			},
		},
		"resize": {
			{
				body:   invokeResize,
				params: []reflect.Type{numberType, numberType},
			},
		},
		"line": {
			{
				body:   invokeLine,
				params: []reflect.Type{pointType, pointType},
			},
		},
		"polygon": {
			{
				body:   invokePolygon,
				params: []reflect.Type{pointSliceType},
			},
			{
				body:   invokePolygonList,
				params: []reflect.Type{listType},
			},
		},
		"circle": {
			{
				body:   invokeCircle,
				params: []reflect.Type{pointType, numberType},
			},
		},
		"intersect": {
			{
				body:   invokeIntersect,
				params: []reflect.Type{lineType, lineType},
			},
		},
		"translate": {
			{
				body:   invokeTranslateLine,
				params: []reflect.Type{lineType, pointType},
			},
			{
				body:   invokeTranslateRect,
				params: []reflect.Type{rectType, pointType},
			},
			{
				body:   invokeTranslatePolygon,
				params: []reflect.Type{polygonType, pointType},
			},
			{
				body:   invokeTranslateCircle,
				params: []reflect.Type{circleType, pointType},
			},
		},
		"clamp": {
			{
				body:   invokeClamp,
				params: []reflect.Type{numberType, numberType, numberType},
			},
			{
				body:   invokeClampRgb,
				params: []reflect.Type{colorType},
			},
		},
		"compose": {
			{
				body:   invokeCompose,
				params: []reflect.Type{colorType, colorType},
			},
		},
		"sum": {
			{
				body:   invokeSumList,
				params: []reflect.Type{listType},
			},
			{
				body:   invokeSumKernel,
				params: []reflect.Type{kernelType},
			},
		},
		"outline": {
			{
				body:   invokeOutlineRect,
				params: []reflect.Type{rectType},
			},
			{
				body:   invokeOutlinePolygon,
				params: []reflect.Type{polygonType},
			},
			{
				body:   invokeOutlineCircle,
				params: []reflect.Type{circleType},
			},
		},
		"hsv": {
			{
				body:   invokeHsv,
				params: []reflect.Type{numberType, numberType, numberType},
			},
			{
				body:   invokeRgb2Hsv,
				params: []reflect.Type{colorType},
			},
		},
		"compare": {
			{
				body:   invokeCompare,
				params: []reflect.Type{valueType, valueType},
			},
		},
		"plot": {
			{
				body:   invokePlot,
				params: []reflect.Type{valueType, colorType},
			},
		},
	}
}

func signature(name string, f functionDecl) string {
	paramTypeNames := make([]string, len(f.params))
	for i, param := range f.params {
		if param.Kind() == reflect.Slice {
			paramTypeNames[i] = param.Elem().Name() + "..."
		} else {
			paramTypeNames[i] = param.Name()
		}
	}
	return fmt.Sprintf("fn %s(%s)", name, strings.Join(paramTypeNames, ", "))
}

func invokeRgb(ir *interpreter, args []Value) (Value, error) {
	return color(lang.NewRgba(
		lang.Number(args[0].(number)),
		lang.Number(args[1].(number)),
		lang.Number(args[2].(number)), 255)), nil
}

func invokeSrgb(ir *interpreter, args []Value) (Value, error) {
	return color(lang.NewSrgba(
		lang.Number(args[0].(number)),
		lang.Number(args[1].(number)),
		lang.Number(args[2].(number)), 1.0)), nil
}

func invokeRgba(ir *interpreter, args []Value) (Value, error) {
	return color(lang.NewRgba(
		lang.Number(args[0].(number)),
		lang.Number(args[1].(number)),
		lang.Number(args[2].(number)),
		lang.Number(args[3].(number)))), nil
}

func invokeRgb2Rgba(ir *interpreter, args []Value) (Value, error) {
	rgb := args[0].(color)
	a := args[1].(number)
	return color(lang.NewRgba(rgb.R, rgb.G, rgb.B, lang.Number(a))), nil
}

func invokeHsv2Rgba(ir *interpreter, args []Value) (Value, error) {
	hsv := args[0].(colorHsv)
	rgb := hsv.rgb()
	a := args[1].(number)
	return color(lang.NewRgba(rgb.R, rgb.G, rgb.B, lang.Number(a))), nil
}

func invokeSrgba(ir *interpreter, args []Value) (Value, error) {
	return color(lang.NewSrgba(
		lang.Number(args[0].(number)),
		lang.Number(args[1].(number)),
		lang.Number(args[2].(number)),
		lang.Number(args[3].(number)))), nil
}

func invokeGrey(ir *interpreter, args []Value) (Value, error) {
	v := lang.Number(args[0].(number))
	return color(lang.NewRgba(v, v, v, 255)), nil
}

func invokeSgrey(ir *interpreter, args []Value) (Value, error) {
	v := lang.Number(args[0].(number))
	return color(lang.NewSrgba(v, v, v, 1.0)), nil
}

func invokeRect(ir *interpreter, args []Value) (Value, error) {
	x, y := int(args[0].(number)), int(args[1].(number))
	return rect{
		Min: image.Point{x, y},
		Max: image.Point{x + int(args[2].(number)+0.5), y + int(args[3].(number)+0.5)},
	}, nil
}

func invokeConvolute(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	return color(ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)), nil
}

func invokeBlt(ir *interpreter, args []Value) (Value, error) {
	rect := args[0].(rect)
	ir.bitmap.Blt(rect.Min.X, rect.Min.Y, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
	return nil, nil
}

func invokeFlip(ir *interpreter, args []Value) (Value, error) {
	imageID := ir.bitmap.Flip()
	ir.assignBounds(false)
	return number(imageID), nil
}

func invokeRecall(ir *interpreter, args []Value) (Value, error) {
	imageID := args[0].(number)
	if err := ir.bitmap.Recall(int(imageID)); err != nil {
		return nil, err
	}
	ir.assignBounds(false)
	return nil, nil
}

func invokeSortKernel(ir *interpreter, args []Value) (Value, error) {
	kernelVal := args[0].(kernel)
	numbers := convertLangNumbersToNumbers(kernelVal.values)
	sort.Sort(numberSlice(numbers))
	result := kernelVal
	result.values = convertNumbersToLangNumbers(numbers)
	return result, nil
}

func invokeSortList(ir *interpreter, args []Value) (Value, error) {
	listVal := args[0].(list)
	result := listVal
	result.elements = append([]Value(nil), result.elements...) // clone elements
	sort.Sort(valueSlice(result.elements))
	return result, nil
}

func invokeSortListFn(ir *interpreter, args []Value) (Value, error) {
	listVal := args[0].(list)
	fn := args[1].(function)

	fnArgs := make([]Value, 2)
	result := listVal
	result.elements = append([]Value(nil), result.elements...) // clone elements

	sort.Slice(result.elements, func(i, j int) bool {
		fnArgs[0] = result.elements[i]
		fnArgs[1] = result.elements[j]
		retVal, err := ir.invokeFunctionExpr("<sort_fn>", fn, fnArgs)
		if err != nil {
			return false
		}
		retNum, ok := retVal.(number)
		return ok && retNum < 0
	})
	return result, nil
}

func invokeFetchRed(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapRed(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeFetchGreen(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapGreen(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeFetchBlue(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapBlue(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeFetchAlpha(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapAlpha(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeSin(ir *interpreter, args []Value) (Value, error) {
	return number(math.Sin(float64(args[0].(number)))), nil
}

func invokeCos(ir *interpreter, args []Value) (Value, error) {
	return number(math.Cos(float64(args[0].(number)))), nil
}

func invokeTan(ir *interpreter, args []Value) (Value, error) {
	return number(math.Tan(float64(args[0].(number)))), nil
}

func invokeAsin(ir *interpreter, args []Value) (Value, error) {
	return number(math.Asin(float64(args[0].(number)))), nil
}

func invokeAcos(ir *interpreter, args []Value) (Value, error) {
	return number(math.Acos(float64(args[0].(number)))), nil
}

func invokeAtan(ir *interpreter, args []Value) (Value, error) {
	return number(math.Atan(float64(args[0].(number)))), nil
}

func invokeAtan2(ir *interpreter, args []Value) (Value, error) {
	return number(math.Atan2(float64(args[0].(number)), float64(args[1].(number)))), nil
}

func invokeSqrt(ir *interpreter, args []Value) (Value, error) {
	return number(math.Sqrt(float64(args[0].(number)))), nil
}

func invokePow(ir *interpreter, args []Value) (Value, error) {
	return number(math.Pow(float64(args[0].(number)), float64(args[1].(number)))), nil
}

func invokeAbs(ir *interpreter, args []Value) (Value, error) {
	return number(math.Abs(float64(args[0].(number)))), nil
}

func invokeRound(ir *interpreter, args []Value) (Value, error) {
	return number(math.Round(float64(args[0].(number)))), nil
}

func invokeFloor(ir *interpreter, args []Value) (Value, error) {
	return number(math.Floor(float64(args[0].(number)))), nil
}

func invokeCeil(ir *interpreter, args []Value) (Value, error) {
	return number(math.Ceil(float64(args[0].(number)))), nil
}

func invokeHypot(ir *interpreter, args []Value) (Value, error) {
	return number(math.Hypot(float64(args[0].(number)), float64(args[1].(number)))), nil
}

func invokeHypotRgb(ir *interpreter, args []Value) (Value, error) {
	a, b := args[0].(color), args[1].(color)
	return color{
		R: lang.Number(math.Hypot(float64(a.R), float64(b.R))),
		G: lang.Number(math.Hypot(float64(a.G), float64(b.G))),
		B: lang.Number(math.Hypot(float64(a.B), float64(b.B))),
		A: a.A,
	}, nil
}

func invokeHypotPoint(ir *interpreter, args []Value) (Value, error) {
	p := args[0].(point)
	return number(math.Hypot(float64(p.X), float64(p.Y))), nil
}

func invokeRandom(ir *interpreter, args []Value) (Value, error) {
	min := args[0].(number)
	max := args[1].(number)
	if min < 0 || max-min <= 0 {
		return nil, fmt.Errorf("invalid range [%g - %g] for random", min, max)
	}
	return number(int(min) + rand.Intn(int(max-min))), nil
}

func invokeRandom01(ir *interpreter, args []Value) (Value, error) {
	return number(rand.Float32()), nil
}

func invokeMax(it *interpreter, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() arguments must not be empty")
	}
	max := number(lang.MinNumber)
	for _, v := range args {
		n := v.(number)
		if n > max {
			max = n
		}
	}
	return max, nil
}

func invokeMaxKernel(it *interpreter, args []Value) (Value, error) {
	kernelVal := args[0].(kernel)
	max := number(lang.MinNumber)
	for _, n := range kernelVal.values {
		nn := number(n)
		if nn > max {
			max = nn
		}
	}
	return max, nil
}

func invokeMaxList(it *interpreter, args []Value) (Value, error) {
	listVal := args[0].(list)
	if len(listVal.elements) == 0 {
		return nil, fmt.Errorf("max() arguments must not be empty")
	}
	var max Value = number(lang.MinNumber)
	for _, v := range listVal.elements {
		cmp, err := v.compare(max)
		if err != nil {
			return nil, err
		}
		if n, ok := cmp.(number); ok && n > 0 {
			max = v
		}
	}
	return max, nil
}

func invokeMin(it *interpreter, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() arguments must not be empty")
	}
	min := number(lang.MaxNumber)
	for _, v := range args {
		n := v.(number)
		if n < min {
			min = n
		}
	}
	return min, nil
}

func invokeMinKernel(it *interpreter, args []Value) (Value, error) {
	kernelVal := args[0].(kernel)
	min := number(lang.MaxNumber)
	for _, n := range kernelVal.values {
		nn := number(n)
		if nn < min {
			min = nn
		}
	}
	return number(min), nil
}

func invokeMinList(it *interpreter, args []Value) (Value, error) {
	listVal := args[0].(list)
	if len(listVal.elements) == 0 {
		return nil, fmt.Errorf("min() arguments must not be empty")
	}
	var min Value = number(lang.MaxNumber)
	for _, v := range listVal.elements {
		cmp, err := v.compare(min)
		if err != nil {
			return nil, err
		}
		if n, ok := cmp.(number); ok && n < 0 {
			min = v
		}
	}
	return min, nil
}

func invokeList(it *interpreter, args []Value) (Value, error) {
	count := args[0].(number)
	val := args[1].(number)
	values := make([]Value, int(count))
	for i := range values {
		values[i] = val
	}
	return list{elements: values}, nil
}

func invokeListFn(it *interpreter, args []Value) (Value, error) {
	count := args[0].(number)
	fn := args[1].(function)
	values := make([]Value, int(count))
	fnArgs := make([]Value, 1)
	for i := range values {
		fnArgs[0] = number(i)
		retVal, err := it.invokeFunctionExpr("<list_fn>", fn, fnArgs)
		if err != nil {
			return nil, err
		}
		values[i] = retVal
	}
	return list{elements: values}, nil
}

func invokeKernel(ir *interpreter, args []Value) (Value, error) {
	width := args[0].(number)
	height := args[1].(number)
	val := args[2].(number)

	values := make([]number, int(width*height))
	for i := range values {
		values[i] = val
	}

	return kernel{width: int(width), height: int(height), values: convertNumbersToLangNumbers(values)}, nil
}

func invokeKernelFn(ir *interpreter, args []Value) (Value, error) {
	width := int(args[0].(number))
	height := int(args[1].(number))
	fn := args[2].(function)
	values := make([]number, width*height)
	fnArgs := make([]Value, 2)

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fnArgs[0] = number(x)
			fnArgs[1] = number(y)
			retVal, err := ir.invokeFunctionExpr("<kernel_fn>", fn, fnArgs)
			if err != nil {
				return nil, err
			}
			retNum, ok := retVal.(number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: function passed to kernel_fn must return number, not %s", reflect.TypeOf(retVal))
			}
			values[i] = retNum
			i++
		}
	}

	return kernel{width: int(width), height: int(height), values: convertNumbersToLangNumbers(values)}, nil
}

func invokeGauss(ir *interpreter, args []Value) (Value, error) {
	radius := int(args[0].(number))
	length := int(radius*2 + 1)

	values := make([]number, int(length*length))
	i := 0
	for y := 0; y < length; y++ {
		var base int
		if y > radius {
			base = length - y - 1
		} else {
			base = y
		}
		for x := 0; x < length; x++ {
			var val int
			if x > radius {
				val = length - x - 1
			} else {
				val = x
			}
			intVal := 1 << uint(base+val)
			values[i] = number(intVal)
			i++
		}
	}

	return kernel{width: int(length), height: int(length), values: convertNumbersToLangNumbers(values)}, nil
}

func invokeResize(ir *interpreter, args []Value) (Value, error) {
	width := args[0].(number)
	height := args[1].(number)

	ir.bitmap.ResizeTarget(int(width), int(height))

	return rect{
		Max: image.Point{int(width), int(height)},
	}, nil
}

func invokeLine(ir *interpreter, args []Value) (Value, error) {
	point1, point2 := args[0].(point), args[1].(point)
	return line{point1, point2}, nil
}

func invokePolygon(ir *interpreter, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("polygon must not be empty")
	}

	points := make([]point, len(args))
	for i, param := range args {
		points[i] = param.(point)
	}

	if points[0] == points[len(points)-1] {
		// first and last are equal -> remove last
		points = points[:len(points)-1]
	}

	return polygon{
		vertices: points,
	}, nil
}

func invokePolygonList(ir *interpreter, args []Value) (Value, error) {
	list := args[0].(list)

	points := make([]point, len(list.elements))
	for i, param := range list.elements {
		pt, ok := param.(point)
		if !ok {
			return nil, fmt.Errorf("type mismatch: polygon_list expects a list of points but found a %s", reflect.TypeOf(param))
		}
		points[i] = pt
	}

	if points[0] == points[len(points)-1] {
		// first and last are equal -> remove last
		points = points[:len(points)-1]
	}

	return polygon{
		vertices: points,
	}, nil
}

func invokeCircle(ir *interpreter, args []Value) (Value, error) {
	center, radius := args[0].(point), args[1].(number)
	return circle{
		center: center,
		radius: radius,
	}, nil
}

func invokeIntersect(ir *interpreter, args []Value) (Value, error) {
	line1, line2 := args[0].(line), args[1].(line)
	p1, p2 := line1.point1, line1.point2
	p3, p4 := line2.point1, line2.point2
	x1, y1, x2, y2 := float32(p1.X), float32(p1.Y), float32(p2.X), float32(p2.Y)
	x3, y3, x4, y4 := float32(p3.X), float32(p3.Y), float32(p4.X), float32(p4.Y)

	if x1 == x2 {
		return intersectVertical(p3, p4, p1.X), nil
	}
	if x3 == x4 {
		return intersectVertical(p1, p2, p3.X), nil
	}

	m1 := (y2 - y1) / (x2 - x1)
	m2 := (y4 - y3) / (x4 - x3)

	if m1 == m2 {
		// the lines are parallel
		return nil, nil
	}

	x := (m1*x1 - m2*x3 + y3 - y1) / (m1 - m2)
	y := (x-x1)*m1 + y1
	return point{int(x + 0.5), int(y + 0.5)}, nil
}

func intersectVertical(p1 point, p2 point, x int) Value {
	if p1.X == p2.X {
		// line is parallel to y axis
		return nil
	}
	return point{
		x,
		((x-p1.X)*(p2.Y-p1.Y)/(p2.X-p1.X) + p1.Y),
	}
}

func invokeTranslateLine(ir *interpreter, args []Value) (Value, error) {
	ln, pt := args[0].(line), args[1].(point)
	return line{
		point1: point{ln.point1.X + pt.X, ln.point1.Y + pt.Y},
		point2: point{ln.point2.X + pt.X, ln.point2.Y + pt.Y},
	}, nil
}

func invokeTranslateRect(ir *interpreter, args []Value) (Value, error) {
	rc, pt := args[0].(rect), args[1].(point)
	return rect(image.Rectangle(rc).Add(image.Point(pt))), nil
}

func invokeTranslatePolygon(ir *interpreter, args []Value) (Value, error) {
	poly, pt := args[0].(polygon), args[1].(point)
	newVertices := make([]point, len(poly.vertices))
	for i, v := range poly.vertices {
		newVertices[i] = point{v.X + pt.X, v.Y + pt.Y}
	}
	return polygon{newVertices}, nil
}

func invokeTranslateCircle(ir *interpreter, args []Value) (Value, error) {
	cir, pt := args[0].(circle), args[1].(point)
	return circle{
		center: point{
			X: cir.center.X + pt.X,
			Y: cir.center.Y + pt.Y,
		},
	}, nil
}

func invokeClamp(ir *interpreter, args []Value) (Value, error) {
	n, min, max := args[0].(number), args[1].(number), args[2].(number)
	if n < min {
		n = min
	}
	if n > max {
		n = max
	}
	return n, nil
}

func invokeClampRgb(ir *interpreter, args []Value) (Value, error) {
	col := args[0].(color)
	return color(lang.Color(col).Clamp()), nil
}

func invokeCompose(ir *interpreter, args []Value) (Value, error) {
	lower, upper := args[0].(color), args[1].(color)
	lowerA, upperA := lang.Color(lower).ScA(), lang.Color(upper).ScA()
	inverseUpperA := 1.0 - upperA
	a := lowerA + (1.0-lowerA)*upperA

	return color(lang.NewRgba(
		((upper.R*upperA + lower.R*lowerA*inverseUpperA) / a).Clamp(),
		((upper.G*upperA + lower.G*lowerA*inverseUpperA) / a).Clamp(),
		((upper.B*upperA + lower.B*lowerA*inverseUpperA) / a).Clamp(),
		255.0*a)), nil
}

func invokeSumList(ir *interpreter, args []Value) (Value, error) {
	list := args[0].(list)
	var sum Value
	var err error
	for _, v := range list.elements {
		if sum == nil {
			sum = v
		} else {
			sum, err = sum.add(v)
		}
		if err != nil {
			return nil, err
		}
	}
	return sum, nil
}

func invokeSumKernel(ir *interpreter, args []Value) (Value, error) {
	kernel := args[0].(kernel)
	var sum number
	for _, v := range kernel.values {
		sum += number(v)
	}
	return sum, nil
}

func invokeOutlineRect(ir *interpreter, args []Value) (Value, error) {
	rc := args[0].(rect)
	lines := make([]Value, 4)
	lines[0] = line{
		point1: point(rc.Min),
		point2: point{rc.Max.X, rc.Min.Y},
	}
	lines[1] = line{
		point1: point{rc.Max.X, rc.Min.Y},
		point2: point(rc.Max),
	}
	lines[2] = line{
		point1: point(rc.Max),
		point2: point{rc.Min.X, rc.Max.Y},
	}
	lines[3] = line{
		point1: point{rc.Min.X, rc.Max.Y},
		point2: point(rc.Min),
	}
	return list{lines}, nil
}

func invokeOutlinePolygon(ir *interpreter, args []Value) (Value, error) {
	poly := args[0].(polygon)
	lines := make([]Value, len(poly.vertices))

	for i, vertex := range poly.vertices {
		inext := i + 1
		if inext >= len(poly.vertices) {
			inext = 0
		}
		nextVertex := poly.vertices[inext]
		lines[i] = line{
			point1: vertex,
			point2: nextVertex,
		}
	}

	return list{lines}, nil
}

func invokeOutlineCircle(ir *interpreter, args []Value) (Value, error) {
	cir := args[0].(circle)
	deg2rad := math.Pi / 180
	var lines []Value
	var prevPt point

	for i := 0; i <= 360; i++ {
		angle := i % 360
		pt := point{
			X: int(float64(cir.center.X) + math.Sin(float64(angle)*deg2rad)*float64(cir.radius)),
			Y: int(float64(cir.center.Y) + math.Cos(float64(angle)*deg2rad)*float64(cir.radius)),
		}
		if i > 0 && pt != prevPt {
			l := line{
				point1: prevPt,
				point2: pt,
			}
			lines = append(lines, l)
		}
		prevPt = pt
	}
	return list{lines}, nil
}

func invokeRgb2Hsv(ir *interpreter, args []Value) (Value, error) {
	rgb := lang.Color(args[0].(color))
	return hsvFromRgb(rgb), nil
}

func invokeHsv2Rgb(ir *interpreter, args []Value) (Value, error) {
	hsv := args[0].(colorHsv)
	return color(hsv.rgb()), nil
}

func invokeHsv(ir *interpreter, args []Value) (Value, error) {
	return colorHsv{
		h: lang.Number(args[0].(number)),
		s: lang.Number(args[1].(number)),
		v: lang.Number(args[2].(number)),
	}, nil
}

func invokeCompare(ir *interpreter, args []Value) (Value, error) {
	return args[0].compare(args[1])
}

func invokePlot(ir *interpreter, args []Value) (Value, error) {
	iterable := args[0]
	color := lang.Color(args[1].(color))

	err := iterable.iterate(func(v Value) error {
		pt, ok := v.(point)
		if !ok {
			return fmt.Errorf("type mismatch: expected point, but found %s", reflect.TypeOf(v))
		}
		ir.bitmap.SetPixel(pt.X, pt.Y, color)
		return nil
	})

	return nil, err
}

func convertNumbersToLangNumbers(numbers []number) []lang.Number {
	result := make([]lang.Number, len(numbers))
	for i, n := range numbers {
		result[i] = lang.Number(n)
	}
	return result
}

func convertLangNumbersToNumbers(numbers []lang.Number) []number {
	result := make([]number, len(numbers))
	for i, n := range numbers {
		result[i] = number(n)
	}
	return result
}
