package lang

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"reflect"
	"sort"
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
	"grey": {
		body:   invokeGrey,
		params: []reflect.Type{numberType},
	},
	"sgrey": {
		body:   invokeSgrey,
		params: []reflect.Type{numberType},
	},
	"rect": {
		body:   invokeRect,
		params: []reflect.Type{numberType, numberType, numberType, numberType},
	},
	"convolute": {
		body:   invokeConvolute,
		params: []reflect.Type{pointType, kernelType},
	},
	"blt": {
		body:   invokeBlt,
		params: []reflect.Type{reflect.TypeOf(rect{})},
	},
	"flip": {
		body:   invokeFlip,
		params: []reflect.Type{},
	},
	"recall": {
		body:   invokeRecall,
		params: []reflect.Type{numberType},
	},
	"sort_kernel": {
		body:   invokeSortKernel,
		params: []reflect.Type{kernelType},
	},
	"sort_list": {
		body:   invokeSortList,
		params: []reflect.Type{listType},
	},
	"fetch_red": {
		body:   invokeFetchRed,
		params: []reflect.Type{pointType, kernelType},
	},
	"fetch_green": {
		body:   invokeFetchGreen,
		params: []reflect.Type{pointType, kernelType},
	},
	"fetch_blue": {
		body:   invokeFetchBlue,
		params: []reflect.Type{pointType, kernelType},
	},
	"fetch_alpha": {
		body:   invokeFetchAlpha,
		params: []reflect.Type{pointType, kernelType},
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
	"pow": {
		body:   invokePow,
		params: []reflect.Type{numberType, numberType},
	},
	"abs": {
		body:   invokeAbs,
		params: []reflect.Type{numberType},
	},
	"round": {
		body:   invokeRound,
		params: []reflect.Type{numberType},
	},
	"hypot": {
		body:   invokeHypot,
		params: []reflect.Type{numberType, numberType},
	},
	"hypot_rgb": {
		body:   invokeHypotRgb,
		params: []reflect.Type{colorType, colorType},
	},
	"random": {
		body:   invokeRandom,
		params: []reflect.Type{numberType},
	},
	"min": {
		body:   invokeMin,
		params: []reflect.Type{reflect.TypeOf([]Number{})},
	},
	"min_kernel": {
		body:   invokeMinKernel,
		params: []reflect.Type{kernelType},
	},
	"min_list": {
		body:   invokeMinList,
		params: []reflect.Type{listType},
	},
	"max": {
		body:   invokeMax,
		params: []reflect.Type{reflect.TypeOf([]Number{})},
	},
	"max_kernel": {
		body:   invokeMaxKernel,
		params: []reflect.Type{kernelType},
	},
	"max_list": {
		body:   invokeMaxList,
		params: []reflect.Type{listType},
	},
	"list": {
		body:   invokeList,
		params: []reflect.Type{numberType, numberType},
	},
	"kernel": {
		body:   invokeKernel,
		params: []reflect.Type{numberType, numberType, numberType},
	},
	"gauss": {
		body:   invokeGauss,
		params: []reflect.Type{numberType},
	},
	"resize": {
		body:   invokeResize,
		params: []reflect.Type{numberType, numberType},
	},
	"line": {
		body:   invokeLine,
		params: []reflect.Type{pointType, pointType},
	},
	"polygon": {
		body:   invokePolygon,
		params: []reflect.Type{reflect.TypeOf([]point{})},
	},
	"intersect": {
		body:   invokeIntersect,
		params: []reflect.Type{reflect.TypeOf(line{}), reflect.TypeOf(line{})},
	},
	"translate_line": {
		body:   invokeTranslateLine,
		params: []reflect.Type{reflect.TypeOf(line{}), pointType},
	},
	"translate_rect": {
		body:   invokeTranslateRect,
		params: []reflect.Type{reflect.TypeOf(rect{}), pointType},
	},
	"translate_polygon": {
		body:   invokeTranslatePolygon,
		params: []reflect.Type{reflect.TypeOf(polygon{}), pointType},
	},
	"clamp": {
		body:   invokeClamp,
		params: []reflect.Type{numberType, numberType, numberType},
	},
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

func invokeRandom(ir *interpreter, args []value) (value, error) {
	return Number(rand.Intn(int(args[0].(Number)))), nil
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
		isGt, err := v.greaterThan(max)
		if err != nil {
			return nil, err
		}
		if isGt.(boolean) {
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
		isLess, err := v.lessThan(min)
		if err != nil {
			return nil, err
		}
		if isLess.(boolean) {
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

	points := make([]point, 0, len(args))
	for _, param := range args {
		points = append(points, param.(point))
	}

	if points[0] == points[len(points)-1] {
		// first and last are equal -> remove last
		points = points[:len(points)-1]
	}

	return polygon{
		vertices: points,
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
