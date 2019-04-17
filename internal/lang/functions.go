package lang

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strings"
)

type functionDecl struct {
	body   func(ir *interpreter, values []value) (value, error)
	params []reflect.Type
}

var numberType = reflect.TypeOf(Number(0))
var pointType = reflect.TypeOf(point{})
var kernelType = reflect.TypeOf(kernel{})
var listType = reflect.TypeOf(list{})
var colorType = reflect.TypeOf(Color{})
var functionType = reflect.TypeOf(function{})
var lineType = reflect.TypeOf(line{})
var rectType = reflect.TypeOf(rect{})
var polygonType = reflect.TypeOf(polygon{})
var numberSliceType = reflect.TypeOf([]Number{})
var pointSliceType = reflect.TypeOf([]point{})
var circleType = reflect.TypeOf(circle{})
var hsvType = reflect.TypeOf(colorHsv{})
var valueType = reflect.TypeOf((*value)(nil)).Elem()

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
		},
		"rgb01": {
			{
				body:   invokeSrgb,
				params: []reflect.Type{numberType, numberType, numberType},
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
		"grey": {
			{
				body:   invokeGrey,
				params: []reflect.Type{numberType},
			},
		},
		"grey01": {
			{
				body:   invokeSgrey,
				params: []reflect.Type{numberType},
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

func invokeRgb(ir *interpreter, args []value) (value, error) {
	return NewRgba(args[0].(Number), args[1].(Number), args[2].(Number), 255), nil
}

func invokeSrgb(ir *interpreter, args []value) (value, error) {
	return NewSrgba(args[0].(Number), args[1].(Number), args[2].(Number), 1.0), nil
}

func invokeRgba(ir *interpreter, args []value) (value, error) {
	return NewRgba(args[0].(Number), args[1].(Number), args[2].(Number), args[3].(Number)), nil
}

func invokeRgb2Rgba(ir *interpreter, args []value) (value, error) {
	rgb := args[0].(Color)
	a := args[1].(Number)
	return NewRgba(rgb.R, rgb.G, rgb.B, a), nil
}

func invokeHsv2Rgba(ir *interpreter, args []value) (value, error) {
	hsv := args[0].(colorHsv)
	rgb := hsv.rgb()
	a := args[1].(Number)
	return NewRgba(rgb.R, rgb.G, rgb.B, a), nil
}

func invokeSrgba(ir *interpreter, args []value) (value, error) {
	return NewSrgba(args[0].(Number), args[1].(Number), args[2].(Number), args[3].(Number)), nil
}

func invokeGrey(ir *interpreter, args []value) (value, error) {
	v := args[0].(Number)
	return NewRgba(v, v, v, 255), nil
}

func invokeSgrey(ir *interpreter, args []value) (value, error) {
	v := args[0].(Number)
	return NewSrgba(v, v, v, 1.0), nil
}

func invokeRect(ir *interpreter, args []value) (value, error) {
	x, y := int(args[0].(Number)), int(args[1].(Number))
	return rect{
		Min: image.Point{x, y},
		Max: image.Point{x + int(args[2].(Number)+0.5), y + int(args[3].(Number)+0.5)},
	}, nil
}

func invokeConvolute(ir *interpreter, args []value) (value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	return ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values), nil
}

func invokeBlt(ir *interpreter, args []value) (value, error) {
	rect := args[0].(rect)
	ir.bitmap.Blt(rect.Min.X, rect.Min.Y, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
	return nil, nil
}

func invokeFlip(ir *interpreter, args []value) (value, error) {
	imageID := ir.bitmap.Flip()
	ir.assignBounds(false)
	return Number(imageID), nil
}

func invokeRecall(ir *interpreter, args []value) (value, error) {
	imageID := args[0].(Number)
	if err := ir.bitmap.Recall(int(imageID)); err != nil {
		return nil, err
	}
	ir.assignBounds(false)
	return nil, nil
}

func invokeSortKernel(ir *interpreter, args []value) (value, error) {
	kernelVal := args[0].(kernel)
	result := kernelVal
	result.values = append([]Number(nil), result.values...) // clone values
	sort.Sort(numberSlice(result.values))
	return result, nil
}

func invokeSortList(ir *interpreter, args []value) (value, error) {
	listVal := args[0].(list)
	result := listVal
	result.elements = append([]value(nil), result.elements...) // clone elements
	sort.Sort(valueSlice(result.elements))
	return result, nil
}

func invokeSortListFn(ir *interpreter, args []value) (value, error) {
	listVal := args[0].(list)
	fn := args[1].(function)

	fnArgs := make([]value, 2)
	result := listVal
	result.elements = append([]value(nil), result.elements...) // clone elements

	sort.Slice(result.elements, func(i, j int) bool {
		fnArgs[0] = result.elements[i]
		fnArgs[1] = result.elements[j]
		retVal, err := ir.invokeFunctionExpr("<sort_fn>", fn, fnArgs)
		if err != nil {
			return false
		}
		retNum, ok := retVal.(Number)
		return ok && retNum < 0
	})
	return result, nil
}

func invokeFetchRed(ir *interpreter, args []value) (value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapRed(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeFetchGreen(ir *interpreter, args []value) (value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapGreen(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeFetchBlue(ir *interpreter, args []value) (value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapBlue(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeFetchAlpha(ir *interpreter, args []value) (value, error) {
	posVal := args[0].(point)
	kernelVal := args[1].(kernel)
	result := kernelVal
	result.values = ir.bitmap.MapAlpha(posVal.X, posVal.Y, kernelVal.width, kernelVal.height, kernelVal.values)
	return result, nil
}

func invokeSin(ir *interpreter, args []value) (value, error) {
	return Number(math.Sin(float64(args[0].(Number)))), nil
}

func invokeCos(ir *interpreter, args []value) (value, error) {
	return Number(math.Cos(float64(args[0].(Number)))), nil
}

func invokeTan(ir *interpreter, args []value) (value, error) {
	return Number(math.Tan(float64(args[0].(Number)))), nil
}

func invokeAsin(ir *interpreter, args []value) (value, error) {
	return Number(math.Asin(float64(args[0].(Number)))), nil
}

func invokeAcos(ir *interpreter, args []value) (value, error) {
	return Number(math.Acos(float64(args[0].(Number)))), nil
}

func invokeAtan(ir *interpreter, args []value) (value, error) {
	return Number(math.Atan(float64(args[0].(Number)))), nil
}

func invokeAtan2(ir *interpreter, args []value) (value, error) {
	return Number(math.Atan2(float64(args[0].(Number)), float64(args[1].(Number)))), nil
}

func invokeSqrt(ir *interpreter, args []value) (value, error) {
	return Number(math.Sqrt(float64(args[0].(Number)))), nil
}

func invokePow(ir *interpreter, args []value) (value, error) {
	return Number(math.Pow(float64(args[0].(Number)), float64(args[1].(Number)))), nil
}

func invokeAbs(ir *interpreter, args []value) (value, error) {
	return Number(math.Abs(float64(args[0].(Number)))), nil
}

func invokeRound(ir *interpreter, args []value) (value, error) {
	return Number(math.Round(float64(args[0].(Number)))), nil
}

func invokeFloor(ir *interpreter, args []value) (value, error) {
	return Number(math.Floor(float64(args[0].(Number)))), nil
}

func invokeCeil(ir *interpreter, args []value) (value, error) {
	return Number(math.Ceil(float64(args[0].(Number)))), nil
}

func invokeHypot(ir *interpreter, args []value) (value, error) {
	return Number(math.Hypot(float64(args[0].(Number)), float64(args[1].(Number)))), nil
}

func invokeHypotRgb(ir *interpreter, args []value) (value, error) {
	a, b := args[0].(Color), args[1].(Color)
	return Color{
		R: Number(math.Hypot(float64(a.R), float64(b.R))),
		G: Number(math.Hypot(float64(a.G), float64(b.G))),
		B: Number(math.Hypot(float64(a.B), float64(b.B))),
		A: a.A,
	}, nil
}

func invokeHypotPoint(ir *interpreter, args []value) (value, error) {
	p := args[0].(point)
	return Number(math.Hypot(float64(p.X), float64(p.Y))), nil
}

func invokeRandom(ir *interpreter, args []value) (value, error) {
	min := args[0].(Number)
	max := args[1].(Number)
	if min < 0 || max-min <= 0 {
		return nil, fmt.Errorf("invalid range [%g - %g] for random", min, max)
	}
	return Number(int(min) + rand.Intn(int(max-min))), nil
}

func invokeMax(it *interpreter, args []value) (value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() arguments must not be empty")
	}
	max := MinNumber
	for _, v := range args {
		n := v.(Number)
		if n > max {
			max = n
		}
	}
	return max, nil
}

func invokeMaxKernel(it *interpreter, args []value) (value, error) {
	kernelVal := args[0].(kernel)
	max := MinNumber
	for _, n := range kernelVal.values {
		if n > max {
			max = n
		}
	}
	return Number(max), nil
}

func invokeMaxList(it *interpreter, args []value) (value, error) {
	listVal := args[0].(list)
	if len(listVal.elements) == 0 {
		return nil, fmt.Errorf("max() arguments must not be empty")
	}
	var max value = MinNumber
	for _, v := range listVal.elements {
		cmp, err := v.compare(max)
		if err != nil {
			return nil, err
		}
		if n, ok := cmp.(Number); ok && n > 0 {
			max = v
		}
	}
	return max, nil
}

func invokeMin(it *interpreter, args []value) (value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() arguments must not be empty")
	}
	min := MaxNumber
	for _, v := range args {
		n := v.(Number)
		if n < min {
			min = n
		}
	}
	return min, nil
}

func invokeMinKernel(it *interpreter, args []value) (value, error) {
	kernelVal := args[0].(kernel)
	min := MaxNumber
	for _, n := range kernelVal.values {
		if n < min {
			min = n
		}
	}
	return Number(min), nil
}

func invokeMinList(it *interpreter, args []value) (value, error) {
	listVal := args[0].(list)
	if len(listVal.elements) == 0 {
		return nil, fmt.Errorf("min() arguments must not be empty")
	}
	var min value = MaxNumber
	for _, v := range listVal.elements {
		cmp, err := v.compare(min)
		if err != nil {
			return nil, err
		}
		if n, ok := cmp.(Number); ok && n < 0 {
			min = v
		}
	}
	return min, nil
}

func invokeList(it *interpreter, args []value) (value, error) {
	count := args[0].(Number)
	val := args[1].(Number)
	values := make([]value, int(count))
	for i := range values {
		values[i] = val
	}
	return list{elements: values}, nil
}

func invokeListFn(it *interpreter, args []value) (value, error) {
	count := args[0].(Number)
	fn := args[1].(function)
	values := make([]value, int(count))
	fnArgs := make([]value, 1)
	for i := range values {
		fnArgs[0] = Number(i)
		retVal, err := it.invokeFunctionExpr("<list_fn>", fn, fnArgs)
		if err != nil {
			return nil, err
		}
		values[i] = retVal
	}
	return list{elements: values}, nil
}

func invokeKernel(ir *interpreter, args []value) (value, error) {
	width := args[0].(Number)
	height := args[1].(Number)
	val := args[2].(Number)

	values := make([]Number, int(width*height))
	for i := range values {
		values[i] = val
	}

	return kernel{width: int(width), height: int(height), values: values}, nil
}

func invokeKernelFn(ir *interpreter, args []value) (value, error) {
	width := int(args[0].(Number))
	height := int(args[1].(Number))
	fn := args[2].(function)
	values := make([]Number, width*height)
	fnArgs := make([]value, 2)

	i := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fnArgs[0] = Number(x)
			fnArgs[1] = Number(y)
			retVal, err := ir.invokeFunctionExpr("<kernel_fn>", fn, fnArgs)
			if err != nil {
				return nil, err
			}
			retNum, ok := retVal.(Number)
			if !ok {
				return nil, fmt.Errorf("type mismatch: function passed to kernel_fn must return number, not %s", reflect.TypeOf(retVal))
			}
			values[i] = retNum
			i++
		}
	}

	return kernel{width: int(width), height: int(height), values: values}, nil
}

func invokeGauss(ir *interpreter, args []value) (value, error) {
	radius := int(args[0].(Number))
	length := int(radius*2 + 1)

	values := make([]Number, int(length*length))
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
			values[i] = Number(intVal)
			i++
		}
	}

	return kernel{width: int(length), height: int(length), values: values}, nil
}

func invokeResize(ir *interpreter, args []value) (value, error) {
	width := args[0].(Number)
	height := args[1].(Number)

	ir.bitmap.ResizeTarget(int(width), int(height))

	return rect{
		Max: image.Point{int(width), int(height)},
	}, nil
}

func invokeLine(ir *interpreter, args []value) (value, error) {
	point1, point2 := args[0].(point), args[1].(point)
	return line{point1, point2}, nil
}

func invokePolygon(ir *interpreter, args []value) (value, error) {
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

func invokePolygonList(ir *interpreter, args []value) (value, error) {
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

func invokeCircle(ir *interpreter, args []value) (value, error) {
	center, radius := args[0].(point), args[1].(Number)
	return circle{
		center: center,
		radius: radius,
	}, nil
}

func invokeIntersect(ir *interpreter, args []value) (value, error) {
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

func intersectVertical(p1 point, p2 point, x int) value {
	if p1.X == p2.X {
		// line is parallel to y axis
		return nil
	}
	return point{
		x,
		((x-p1.X)*(p2.Y-p1.Y)/(p2.X-p1.X) + p1.Y),
	}
}

func invokeTranslateLine(ir *interpreter, args []value) (value, error) {
	ln, pt := args[0].(line), args[1].(point)
	return line{
		point1: point{ln.point1.X + pt.X, ln.point1.Y + pt.Y},
		point2: point{ln.point2.X + pt.X, ln.point2.Y + pt.Y},
	}, nil
}

func invokeTranslateRect(ir *interpreter, args []value) (value, error) {
	rc, pt := args[0].(rect), args[1].(point)
	return rect(image.Rectangle(rc).Add(image.Point(pt))), nil
}

func invokeTranslatePolygon(ir *interpreter, args []value) (value, error) {
	poly, pt := args[0].(polygon), args[1].(point)
	newVertices := make([]point, len(poly.vertices))
	for i, v := range poly.vertices {
		newVertices[i] = point{v.X + pt.X, v.Y + pt.Y}
	}
	return polygon{newVertices}, nil
}

func invokeTranslateCircle(ir *interpreter, args []value) (value, error) {
	cir, pt := args[0].(circle), args[1].(point)
	return circle{
		center: point{
			X: cir.center.X + pt.X,
			Y: cir.center.Y + pt.Y,
		},
	}, nil
}

func invokeClamp(ir *interpreter, args []value) (value, error) {
	n, min, max := args[0].(Number), args[1].(Number), args[2].(Number)
	if n < min {
		n = min
	}
	if n > max {
		n = max
	}
	return n, nil
}

func invokeClampRgb(ir *interpreter, args []value) (value, error) {
	color := args[0].(Color)
	return color.Clamp(), nil
}

func invokeCompose(ir *interpreter, args []value) (value, error) {
	lower, upper := args[0].(Color), args[1].(Color)
	lowerA, upperA := lower.ScA(), upper.ScA()
	inverseUpperA := 1.0 - upperA
	a := lowerA + (1.0-lowerA)*upperA

	return NewRgba(
		clamp((upper.R*upperA+lower.R*lowerA*inverseUpperA)/a),
		clamp((upper.G*upperA+lower.G*lowerA*inverseUpperA)/a),
		clamp((upper.B*upperA+lower.B*lowerA*inverseUpperA)/a),
		clamp(255.0*a)), nil
}

func invokeSumList(ir *interpreter, args []value) (value, error) {
	list := args[0].(list)
	var sum value
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

func invokeSumKernel(ir *interpreter, args []value) (value, error) {
	kernel := args[0].(kernel)
	var sum Number
	for _, v := range kernel.values {
		sum += v
	}
	return sum, nil
}

func invokeOutlineRect(ir *interpreter, args []value) (value, error) {
	rc := args[0].(rect)
	lines := make([]value, 4)
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

func invokeOutlinePolygon(ir *interpreter, args []value) (value, error) {
	poly := args[0].(polygon)
	lines := make([]value, len(poly.vertices))

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

func invokeOutlineCircle(ir *interpreter, args []value) (value, error) {
	cir := args[0].(circle)
	deg2rad := math.Pi / 180
	var lines []value
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

func invokeRgb2Hsv(ir *interpreter, args []value) (value, error) {
	rgb := args[0].(Color)
	return hsvFromRgb(rgb), nil
}

func invokeHsv2Rgb(ir *interpreter, args []value) (value, error) {
	hsv := args[0].(colorHsv)
	return hsv.rgb(), nil
}

func invokeHsv(ir *interpreter, args []value) (value, error) {
	return colorHsv{
		h: args[0].(Number),
		s: args[1].(Number),
		v: args[2].(Number),
	}, nil
}

func invokeCompare(ir *interpreter, args []value) (value, error) {
	return args[0].compare(args[1])
}

func invokePlot(ir *interpreter, args []value) (value, error) {
	iterable := args[0]
	color := args[1].(Color)

	err := iterable.iterate(func(v value) error {
		pt, ok := v.(point)
		if !ok {
			return fmt.Errorf("type mismatch: expected point, but found %s", reflect.TypeOf(v))
		}
		ir.bitmap.SetPixel(pt.X, pt.Y, color)
		return nil
	})

	return nil, err
}
