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

type FunctionDecl struct {
	body   func(ir *interpreter, values []Value) (Value, error)
	params []reflect.Type
}

func PrintFunctions() string {
	names := make([]string, 0, len(functions))
	for name := range functions {
		names = append(names, name)
	}
	sort.Sort(sort.StringSlice(names))
	buf := strings.Builder{}
	for _, name := range names {
		buf.WriteString(name)
		buf.WriteString("\n")
		for _, decl := range functions[name] {
			buf.WriteString("  ")
			buf.WriteString(signature(name, decl))
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

var numberType = reflect.TypeOf(Number(0))
var pointType = reflect.TypeOf(Point{})
var kernelType = reflect.TypeOf(Kernel{})
var listType = reflect.TypeOf(List{})
var colorType = reflect.TypeOf(Color{})
var functionType = reflect.TypeOf(Function{})
var lineType = reflect.TypeOf(Line{})
var rectType = reflect.TypeOf(Rect{})
var polygonType = reflect.TypeOf(Polygon{})
var numberSliceType = reflect.TypeOf([]Number{})
var pointSliceType = reflect.TypeOf([]Point{})
var circleType = reflect.TypeOf(Circle{})
var hsvType = reflect.TypeOf(ColorHsv{})
var valueType = reflect.TypeOf((*Value)(nil)).Elem()

var functions map[string][]FunctionDecl

func initFunctions() {
	if functions != nil {
		return
	}
	functions = map[string][]FunctionDecl{
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
		"trunc": {
			{
				body:   invokeTrunc,
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

func signature(name string, f FunctionDecl) string {
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
	return Color(lang.NewRgba(
		lang.Number(args[0].(Number)),
		lang.Number(args[1].(Number)),
		lang.Number(args[2].(Number)), 255)), nil
}

func invokeSrgb(ir *interpreter, args []Value) (Value, error) {
	return Color(lang.NewSrgba(
		lang.Number(args[0].(Number)),
		lang.Number(args[1].(Number)),
		lang.Number(args[2].(Number)), 1.0)), nil
}

func invokeRgba(ir *interpreter, args []Value) (Value, error) {
	return Color(lang.NewRgba(
		lang.Number(args[0].(Number)),
		lang.Number(args[1].(Number)),
		lang.Number(args[2].(Number)),
		lang.Number(args[3].(Number)))), nil
}

func invokeRgb2Rgba(ir *interpreter, args []Value) (Value, error) {
	rgb := args[0].(Color)
	a := args[1].(Number)
	return Color(lang.NewRgba(rgb.R, rgb.G, rgb.B, lang.Number(a))), nil
}

func invokeHsv2Rgba(ir *interpreter, args []Value) (Value, error) {
	hsv := args[0].(ColorHsv)
	rgb := hsv.rgb()
	a := args[1].(Number)
	return Color(lang.NewRgba(rgb.R, rgb.G, rgb.B, lang.Number(a))), nil
}

func invokeSrgba(ir *interpreter, args []Value) (Value, error) {
	return Color(lang.NewSrgba(
		lang.Number(args[0].(Number)),
		lang.Number(args[1].(Number)),
		lang.Number(args[2].(Number)),
		lang.Number(args[3].(Number)))), nil
}

func invokeGrey(ir *interpreter, args []Value) (Value, error) {
	v := lang.Number(args[0].(Number))
	return Color(lang.NewRgba(v, v, v, 255)), nil
}

func invokeSgrey(ir *interpreter, args []Value) (Value, error) {
	v := lang.Number(args[0].(Number))
	return Color(lang.NewSrgba(v, v, v, 1.0)), nil
}

func invokeRect(ir *interpreter, args []Value) (Value, error) {
	x, y := int(args[0].(Number)), int(args[1].(Number))
	return Rect{
		Min: image.Point{x, y},
		Max: image.Point{x + int(args[2].(Number)+0.5), y + int(args[3].(Number)+0.5)},
	}, nil
}

func invokeConvolute(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(Point)
	kernelVal := args[1].(Kernel)
	return Color(ir.bitmap.Convolute(posVal.X, posVal.Y, kernelVal.Width, kernelVal.Height, kernelVal.Values)), nil
}

func invokeBlt(ir *interpreter, args []Value) (Value, error) {
	rect := args[0].(Rect)
	ir.bitmap.Blt(rect.Min.X, rect.Min.Y, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
	return nil, nil
}

func invokeFlip(ir *interpreter, args []Value) (Value, error) {
	imageID := ir.bitmap.Flip()
	ir.assignBounds(false)
	return Number(imageID), nil
}

func invokeRecall(ir *interpreter, args []Value) (Value, error) {
	imageID := args[0].(Number)
	if err := ir.bitmap.Recall(int(imageID)); err != nil {
		return nil, err
	}
	ir.assignBounds(false)
	return nil, nil
}

func invokeSortKernel(ir *interpreter, args []Value) (Value, error) {
	kernelVal := args[0].(Kernel)
	numbers := convertLangNumbersToNumbers(kernelVal.Values)
	sort.Sort(numberSlice(numbers))
	result := kernelVal
	result.Values = convertNumbersToLangNumbers(numbers)
	return result, nil
}

func invokeSortList(ir *interpreter, args []Value) (Value, error) {
	listVal := args[0].(List)
	result := listVal
	result.Elements = append([]Value(nil), result.Elements...) // clone elements
	sort.Sort(valueSlice(result.Elements))
	return result, nil
}

func invokeSortListFn(ir *interpreter, args []Value) (Value, error) {
	listVal := args[0].(List)
	fn := args[1].(Function)

	fnArgs := make([]Value, 2)
	result := listVal
	result.Elements = append([]Value(nil), result.Elements...) // clone elements

	sort.Slice(result.Elements, func(i, j int) bool {
		fnArgs[0] = result.Elements[i]
		fnArgs[1] = result.Elements[j]
		retVal, err := ir.invokeFunctionExpr("<sort_fn>", fn, fnArgs)
		if err != nil {
			return false
		}
		retNum, ok := retVal.(Number)
		return ok && retNum < 0
	})
	return result, nil
}

func invokeFetchRed(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(Point)
	kernelVal := args[1].(Kernel)
	result := kernelVal
	result.Values = ir.bitmap.MapRed(posVal.X, posVal.Y, kernelVal.Width, kernelVal.Height, kernelVal.Values)
	return result, nil
}

func invokeFetchGreen(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(Point)
	kernelVal := args[1].(Kernel)
	result := kernelVal
	result.Values = ir.bitmap.MapGreen(posVal.X, posVal.Y, kernelVal.Width, kernelVal.Height, kernelVal.Values)
	return result, nil
}

func invokeFetchBlue(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(Point)
	kernelVal := args[1].(Kernel)
	result := kernelVal
	result.Values = ir.bitmap.MapBlue(posVal.X, posVal.Y, kernelVal.Width, kernelVal.Height, kernelVal.Values)
	return result, nil
}

func invokeFetchAlpha(ir *interpreter, args []Value) (Value, error) {
	posVal := args[0].(Point)
	kernelVal := args[1].(Kernel)
	result := kernelVal
	result.Values = ir.bitmap.MapAlpha(posVal.X, posVal.Y, kernelVal.Width, kernelVal.Height, kernelVal.Values)
	return result, nil
}

func invokeSin(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Sin(float64(args[0].(Number)))), nil
}

func invokeCos(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Cos(float64(args[0].(Number)))), nil
}

func invokeTan(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Tan(float64(args[0].(Number)))), nil
}

func invokeAsin(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Asin(float64(args[0].(Number)))), nil
}

func invokeAcos(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Acos(float64(args[0].(Number)))), nil
}

func invokeAtan(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Atan(float64(args[0].(Number)))), nil
}

func invokeAtan2(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Atan2(float64(args[0].(Number)), float64(args[1].(Number)))), nil
}

func invokeSqrt(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Sqrt(float64(args[0].(Number)))), nil
}

func invokePow(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Pow(float64(args[0].(Number)), float64(args[1].(Number)))), nil
}

func invokeAbs(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Abs(float64(args[0].(Number)))), nil
}

func invokeRound(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Round(float64(args[0].(Number)))), nil
}

func invokeFloor(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Floor(float64(args[0].(Number)))), nil
}

func invokeCeil(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Ceil(float64(args[0].(Number)))), nil
}

func invokeTrunc(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Trunc(float64(args[0].(Number)))), nil
}

func invokeHypot(ir *interpreter, args []Value) (Value, error) {
	return Number(math.Hypot(float64(args[0].(Number)), float64(args[1].(Number)))), nil
}

func invokeHypotRgb(ir *interpreter, args []Value) (Value, error) {
	a, b := args[0].(Color), args[1].(Color)
	return Color{
		R: lang.Number(math.Hypot(float64(a.R), float64(b.R))),
		G: lang.Number(math.Hypot(float64(a.G), float64(b.G))),
		B: lang.Number(math.Hypot(float64(a.B), float64(b.B))),
		A: a.A,
	}, nil
}

func invokeHypotPoint(ir *interpreter, args []Value) (Value, error) {
	p := args[0].(Point)
	return Number(math.Hypot(float64(p.X), float64(p.Y))), nil
}

func invokeRandom(ir *interpreter, args []Value) (Value, error) {
	min := args[0].(Number)
	max := args[1].(Number)
	if min < 0 || max-min <= 0 {
		return nil, fmt.Errorf("invalid range [%g - %g] for random", min, max)
	}
	return Number(int(min) + rand.Intn(int(max-min))), nil
}

func invokeRandom01(ir *interpreter, args []Value) (Value, error) {
	return Number(rand.Float32()), nil
}

func invokeMax(it *interpreter, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() arguments must Not be empty")
	}
	max := Number(lang.MinNumber)
	for _, v := range args {
		n := v.(Number)
		if n > max {
			max = n
		}
	}
	return max, nil
}

func invokeMaxKernel(it *interpreter, args []Value) (Value, error) {
	kernelVal := args[0].(Kernel)
	max := Number(lang.MinNumber)
	for _, n := range kernelVal.Values {
		nn := Number(n)
		if nn > max {
			max = nn
		}
	}
	return max, nil
}

func invokeMaxList(it *interpreter, args []Value) (Value, error) {
	listVal := args[0].(List)
	if len(listVal.Elements) == 0 {
		return nil, fmt.Errorf("max() arguments must Not be empty")
	}
	var max Value = Number(lang.MinNumber)
	for _, v := range listVal.Elements {
		cmp, err := v.Compare(max)
		if err != nil {
			return nil, err
		}
		if n, ok := cmp.(Number); ok && n > 0 {
			max = v
		}
	}
	return max, nil
}

func invokeMin(it *interpreter, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("min() arguments must Not be empty")
	}
	min := Number(lang.MaxNumber)
	for _, v := range args {
		n := v.(Number)
		if n < min {
			min = n
		}
	}
	return min, nil
}

func invokeMinKernel(it *interpreter, args []Value) (Value, error) {
	kernelVal := args[0].(Kernel)
	min := Number(lang.MaxNumber)
	for _, n := range kernelVal.Values {
		nn := Number(n)
		if nn < min {
			min = nn
		}
	}
	return Number(min), nil
}

func invokeMinList(it *interpreter, args []Value) (Value, error) {
	listVal := args[0].(List)
	if len(listVal.Elements) == 0 {
		return nil, fmt.Errorf("min() arguments must Not be empty")
	}
	var min Value = Number(lang.MaxNumber)
	for _, v := range listVal.Elements {
		cmp, err := v.Compare(min)
		if err != nil {
			return nil, err
		}
		if n, ok := cmp.(Number); ok && n < 0 {
			min = v
		}
	}
	return min, nil
}

func invokeList(it *interpreter, args []Value) (Value, error) {
	count := args[0].(Number)
	val := args[1].(Number)
	values := make([]Value, int(count))
	for i := range values {
		values[i] = val
	}
	return List{Elements: values}, nil
}

func invokeListFn(it *interpreter, args []Value) (Value, error) {
	count := args[0].(Number)
	fn := args[1].(Function)
	values := make([]Value, int(count))
	fnArgs := make([]Value, 1)
	for i := range values {
		fnArgs[0] = Number(i)
		retVal, err := it.invokeFunctionExpr("<list_fn>", fn, fnArgs)
		if err != nil {
			return nil, err
		}
		values[i] = retVal
	}
	return List{Elements: values}, nil
}

func invokeKernel(ir *interpreter, args []Value) (Value, error) {
	width := args[0].(Number)
	height := args[1].(Number)
	val := args[2].(Number)

	values := make([]Number, int(width*height))
	for i := range values {
		values[i] = val
	}

	return Kernel{Width: int(width), Height: int(height), Values: convertNumbersToLangNumbers(values)}, nil
}

func invokeKernelFn(ir *interpreter, args []Value) (Value, error) {
	width := int(args[0].(Number))
	height := int(args[1].(Number))
	fn := args[2].(Function)
	values := make([]Number, width*height)
	fnArgs := make([]Value, 2)

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
				return nil, fmt.Errorf("type mismatch: function passed to kernel_fn must return number, Not %s", reflect.TypeOf(retVal))
			}
			values[i] = retNum
			i++
		}
	}

	return Kernel{Width: int(width), Height: int(height), Values: convertNumbersToLangNumbers(values)}, nil
}

func invokeGauss(ir *interpreter, args []Value) (Value, error) {
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

	return Kernel{Width: int(length), Height: int(length), Values: convertNumbersToLangNumbers(values)}, nil
}

func invokeResize(ir *interpreter, args []Value) (Value, error) {
	width := args[0].(Number)
	height := args[1].(Number)

	ir.bitmap.ResizeTarget(int(width), int(height))

	return Rect{
		Max: image.Point{int(width), int(height)},
	}, nil
}

func invokeLine(ir *interpreter, args []Value) (Value, error) {
	point1, point2 := args[0].(Point), args[1].(Point)
	return Line{point1, point2}, nil
}

func invokePolygon(ir *interpreter, args []Value) (Value, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("polygon must Not be empty")
	}

	points := make([]Point, len(args))
	for i, param := range args {
		points[i] = param.(Point)
	}

	if points[0] == points[len(points)-1] {
		// first and last are equal -> remove last
		points = points[:len(points)-1]
	}

	return Polygon{
		Vertices: points,
	}, nil
}

func invokePolygonList(ir *interpreter, args []Value) (Value, error) {
	list := args[0].(List)

	points := make([]Point, len(list.Elements))
	for i, param := range list.Elements {
		pt, ok := param.(Point)
		if !ok {
			return nil, fmt.Errorf("type mismatch: polygon_list expects a list of points but found a %s", reflect.TypeOf(param))
		}
		points[i] = pt
	}

	if points[0] == points[len(points)-1] {
		// first and last are equal -> remove last
		points = points[:len(points)-1]
	}

	return Polygon{
		Vertices: points,
	}, nil
}

func invokeCircle(ir *interpreter, args []Value) (Value, error) {
	center, radius := args[0].(Point), args[1].(Number)
	return Circle{
		Center: center,
		Radius: radius,
	}, nil
}

func invokeIntersect(ir *interpreter, args []Value) (Value, error) {
	line1, line2 := args[0].(Line), args[1].(Line)
	p1, p2 := line1.Point1, line1.Point2
	p3, p4 := line2.Point1, line2.Point2
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
	return Point{int(x + 0.5), int(y + 0.5)}, nil
}

func intersectVertical(p1 Point, p2 Point, x int) Value {
	if p1.X == p2.X {
		// line is parallel to y axis
		return nil
	}
	return Point{
		x,
		((x-p1.X)*(p2.Y-p1.Y)/(p2.X-p1.X) + p1.Y),
	}
}

func invokeTranslateLine(ir *interpreter, args []Value) (Value, error) {
	ln, pt := args[0].(Line), args[1].(Point)
	return Line{
		Point1: Point{ln.Point1.X + pt.X, ln.Point1.Y + pt.Y},
		Point2: Point{ln.Point2.X + pt.X, ln.Point2.Y + pt.Y},
	}, nil
}

func invokeTranslateRect(ir *interpreter, args []Value) (Value, error) {
	rc, pt := args[0].(Rect), args[1].(Point)
	return Rect(image.Rectangle(rc).Add(image.Point(pt))), nil
}

func invokeTranslatePolygon(ir *interpreter, args []Value) (Value, error) {
	poly, pt := args[0].(Polygon), args[1].(Point)
	newVertices := make([]Point, len(poly.Vertices))
	for i, v := range poly.Vertices {
		newVertices[i] = Point{v.X + pt.X, v.Y + pt.Y}
	}
	return Polygon{newVertices}, nil
}

func invokeTranslateCircle(ir *interpreter, args []Value) (Value, error) {
	cir, pt := args[0].(Circle), args[1].(Point)
	return Circle{
		Center: Point{
			X: cir.Center.X + pt.X,
			Y: cir.Center.Y + pt.Y,
		},
	}, nil
}

func invokeClamp(ir *interpreter, args []Value) (Value, error) {
	n, min, max := args[0].(Number), args[1].(Number), args[2].(Number)
	if n < min {
		n = min
	}
	if n > max {
		n = max
	}
	return n, nil
}

func invokeClampRgb(ir *interpreter, args []Value) (Value, error) {
	col := args[0].(Color)
	return Color(lang.Color(col).Clamp()), nil
}

func invokeCompose(ir *interpreter, args []Value) (Value, error) {
	lower, upper := args[0].(Color), args[1].(Color)
	lowerA, upperA := lang.Color(lower).ScA(), lang.Color(upper).ScA()
	inverseUpperA := 1.0 - upperA
	a := lowerA + (1.0-lowerA)*upperA

	return Color(lang.NewRgba(
		((upper.R*upperA + lower.R*lowerA*inverseUpperA) / a).Clamp(),
		((upper.G*upperA + lower.G*lowerA*inverseUpperA) / a).Clamp(),
		((upper.B*upperA + lower.B*lowerA*inverseUpperA) / a).Clamp(),
		255.0*a)), nil
}

func invokeSumList(ir *interpreter, args []Value) (Value, error) {
	list := args[0].(List)
	var sum Value
	var err error
	for _, v := range list.Elements {
		if sum == nil {
			sum = v
		} else {
			sum, err = sum.Add(v)
		}
		if err != nil {
			return nil, err
		}
	}
	return sum, nil
}

func invokeSumKernel(ir *interpreter, args []Value) (Value, error) {
	kernel := args[0].(Kernel)
	var sum Number
	for _, v := range kernel.Values {
		sum += Number(v)
	}
	return sum, nil
}

func invokeOutlineRect(ir *interpreter, args []Value) (Value, error) {
	rc := args[0].(Rect)
	lines := make([]Value, 4)
	lines[0] = Line{
		Point1: Point(rc.Min),
		Point2: Point{rc.Max.X, rc.Min.Y},
	}
	lines[1] = Line{
		Point1: Point{rc.Max.X, rc.Min.Y},
		Point2: Point(rc.Max),
	}
	lines[2] = Line{
		Point1: Point(rc.Max),
		Point2: Point{rc.Min.X, rc.Max.Y},
	}
	lines[3] = Line{
		Point1: Point{rc.Min.X, rc.Max.Y},
		Point2: Point(rc.Min),
	}
	return List{lines}, nil
}

func invokeOutlinePolygon(ir *interpreter, args []Value) (Value, error) {
	poly := args[0].(Polygon)
	lines := make([]Value, len(poly.Vertices))

	for i, vertex := range poly.Vertices {
		inext := i + 1
		if inext >= len(poly.Vertices) {
			inext = 0
		}
		nextVertex := poly.Vertices[inext]
		lines[i] = Line{
			Point1: vertex,
			Point2: nextVertex,
		}
	}

	return List{lines}, nil
}

func invokeOutlineCircle(ir *interpreter, args []Value) (Value, error) {
	cir := args[0].(Circle)
	deg2rad := math.Pi / 180
	var lines []Value
	var prevPt Point

	for i := 0; i <= 360; i++ {
		angle := i % 360
		pt := Point{
			X: int(float64(cir.Center.X) + math.Sin(float64(angle)*deg2rad)*float64(cir.Radius)),
			Y: int(float64(cir.Center.Y) + math.Cos(float64(angle)*deg2rad)*float64(cir.Radius)),
		}
		if i > 0 && pt != prevPt {
			l := Line{
				Point1: prevPt,
				Point2: pt,
			}
			lines = append(lines, l)
		}
		prevPt = pt
	}
	return List{lines}, nil
}

func invokeRgb2Hsv(ir *interpreter, args []Value) (Value, error) {
	rgb := lang.Color(args[0].(Color))
	return hsvFromRgb(rgb), nil
}

func invokeHsv2Rgb(ir *interpreter, args []Value) (Value, error) {
	hsv := args[0].(ColorHsv)
	return Color(hsv.rgb()), nil
}

func invokeHsv(ir *interpreter, args []Value) (Value, error) {
	return ColorHsv{
		H: lang.Number(args[0].(Number)),
		S: lang.Number(args[1].(Number)),
		V: lang.Number(args[2].(Number)),
	}, nil
}

func invokeCompare(ir *interpreter, args []Value) (Value, error) {
	return args[0].Compare(args[1])
}

func invokePlot(ir *interpreter, args []Value) (Value, error) {
	iterable := args[0]
	color := lang.Color(args[1].(Color))

	err := iterable.Iterate(func(v Value) error {
		pt, ok := v.(Point)
		if !ok {
			return fmt.Errorf("type mismatch: expected point, but found %s", reflect.TypeOf(v))
		}
		ir.bitmap.SetPixel(pt.X, pt.Y, color)
		return nil
	})

	return nil, err
}

func convertNumbersToLangNumbers(numbers []Number) []lang.Number {
	result := make([]lang.Number, len(numbers))
	for i, n := range numbers {
		result[i] = lang.Number(n)
	}
	return result
}

func convertLangNumbersToNumbers(numbers []lang.Number) []Number {
	result := make([]Number, len(numbers))
	for i, n := range numbers {
		result[i] = Number(n)
	}
	return result
}
