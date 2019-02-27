package lang

import (
	"fmt"
	"image"
	"math"
	"reflect"
	"sort"
)

type polygon struct {
	vertices []point
}

func (p polygon) equals(other value) (value, error) {
	if r, ok := other.(polygon); ok {
		if len(r.vertices) != len(p.vertices) {
			return falseVal, nil
		}
		for i, v := range p.vertices {
			if v != r.vertices[i] {
				return falseVal, nil
			}
		}
		return boolean(true), nil
	}
	return falseVal, nil
}

func (p polygon) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon > %s not supported", reflect.TypeOf(other))
}

func (p polygon) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon >= %s not supported", reflect.TypeOf(other))
}

func (p polygon) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon < %s not supported", reflect.TypeOf(other))
}

func (p polygon) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon <= %s not supported", reflect.TypeOf(other))
}

func (p polygon) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon + %s not supported", reflect.TypeOf(other))
}

func (p polygon) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon - %s not supported", reflect.TypeOf(other))
}

func (p polygon) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon * %s not supported", reflect.TypeOf(other))
}

func (p polygon) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon / %s not supported", reflect.TypeOf(other))
}

func (p polygon) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon %% %s not supported", reflect.TypeOf(other))
}

func (p polygon) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		for _, pt := range p.vertices {
			v, err := pt.in(r)
			if err != nil {
				return nil, err
			}
			if !v.(boolean) {
				return falseVal, nil
			}
		}
		return boolean(true), nil
	}
	return nil, fmt.Errorf("type mismatch: polygon in %s not supported", reflect.TypeOf(other))
}

func (p polygon) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: '-polygon' not supported")
}

func (p polygon) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not polygon' not supported")
}

func (p polygon) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @polygon not supported")
}

func (p polygon) property(ident string) (value, error) {
	switch ident {
	case "bounds":
		return p.bounds(), nil
	case "vertices":
		values := make([]value, len(p.vertices))
		for i, vertex := range p.vertices {
			values[i] = vertex
		}
		return list{
			elements: values,
		}, nil
	}
	return baseProperty(p, ident)
}

func (p polygon) printStr() string {
	return fmt.Sprintf("polygon(%v)", p.vertices)
}

func (p polygon) iterateHorizLine(x1, y, x2 int, visit func(value) error) error {
	if x1 > x2 {
		temp := x1
		x1 = x2
		x2 = temp
	}
	for x := x1; x <= x2; x++ {
		if err := visit(point{x, y}); err != nil {
			return err
		}
	}
	return nil
}

func (p polygon) bounds() rect {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for _, pt := range p.vertices {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}
	return rect{
		Min: image.Point{minX, minY},
		Max: image.Point{maxX, maxY},
	}
}

func (p polygon) iterate(visit func(value) error) error {
	// buffer to store the x-coordinates of intersections of the polygon with some horizontal line
	xs := make([]int, 0, len(p.vertices))

	// determine maxima
	bounds := p.bounds()

	if bounds.Min.Y == bounds.Max.Y {
		//Special case: polygon only 1 pixel high.
		if err := p.iterateHorizLine(bounds.Min.X, bounds.Min.Y, bounds.Max.Y, visit); err != nil {
			return err
		}
	}

	// draw, scanning y
	// ----------------
	// the algorithm uses a horizontal line (y) that moves from top to the
	// bottom of the polygon:
	// 1. search intersections with the border lines
	// 2. sort intersections (x_intersect)
	// 3. each two x-coordinates in x_intersect are then inside the polygon
	//    (drawhorzlineclip for a pair of two such points)
	//
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		intersectionXs := xs[:0] // re-use xs buffer to avoid allocations

		for i, pt := range p.vertices {
			previ := i - 1
			if previ < 0 {
				previ = len(p.vertices) - 1
			}
			prevPt := p.vertices[previ]
			y1, y2 := prevPt.Y, pt.Y
			var x1, x2 int
			if y1 < y2 {
				x1 = prevPt.X
				x2 = pt.X
			} else if y1 > y2 {
				y2 = prevPt.Y
				y1 = pt.Y
				x2 = prevPt.X
				x1 = pt.X
			} else { // y1 == y2 : has to be handled as special case (below)
				continue
			}
			if y >= y1 && y < y2 || y == bounds.Max.Y && y2 == bounds.Max.Y {
				// add intersection if y crosses the edge (excluding the lower end), or when we are on the lowest line (maxy)
				intersectionXs = append(intersectionXs, (y-y1)*(x2-x1)/(y2-y1)+x1)
			}
		}
		sort.Ints(intersectionXs)

		for i := 0; i < len(intersectionXs); i += 2 {
			if err := p.iterateHorizLine(intersectionXs[i], y, intersectionXs[i+1], visit); err != nil {
				return err
			}
		}
	}

	// finally, a special case is not handled by above algorithm:
	// for two border points with same height miny < y < maxy,
	// sometimes the line between them is not colored:
	// this happens when the line will be a lower border line of the polygon
	// (eg we are inside the polygon with a smaller y, and outside with a bigger y),
	// So we loop for border lines that are horizontal.
	//
	previ := len(p.vertices) - 1
	for i, pt := range p.vertices {
		if (bounds.Min.Y < pt.Y) && p.vertices[previ].Y == pt.Y && pt.Y <= bounds.Max.Y {
			if err := p.iterateHorizLine(pt.X, pt.Y, p.vertices[previ].X, visit); err != nil {
				return err
			}
		}
		previ = i
	}
	return nil
}

func (p polygon) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon[index] not supported")
}

func (p polygon) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: polygon[%s] not supported", reflect.TypeOf(index))
}

func (p polygon) runtimeTypeName() string {
	return "polygon"
}

func (p polygon) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon  :: [%s] not supported", reflect.TypeOf(val))
}
