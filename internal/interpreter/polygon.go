package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"image"
	"math"
	"reflect"
	"sort"
	"strings"
)

type Polygon struct {
	Vertices []Point
}

func (p Polygon) Compare(other Value) (Value, error) {
	if r, ok := other.(Polygon); ok {
		if len(r.Vertices) != len(p.Vertices) {
			return Number(1), nil
		}
		for i, v := range p.Vertices {
			if v != r.Vertices[i] {
				return Number(1), nil
			}
		}
		return Number(0), nil
	}
	return nil, nil
}

func (p Polygon) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon + %s Not supported", reflect.TypeOf(other))
}

func (p Polygon) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon - %s Not supported", reflect.TypeOf(other))
}

func (p Polygon) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon * %s Not supported", reflect.TypeOf(other))
}

func (p Polygon) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon / %s Not supported", reflect.TypeOf(other))
}

func (p Polygon) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon %% %s Not supported", reflect.TypeOf(other))
}

func (p Polygon) In(other Value) (Value, error) {
	if r, ok := other.(Rect); ok {
		for _, pt := range p.Vertices {
			v, err := pt.In(r)
			if err != nil {
				return nil, err
			}
			if !v.(Boolean) {
				return Boolean(lang.FalseVal), nil
			}
		}
		return Boolean(true), nil
	}
	return nil, fmt.Errorf("type mismatch: polygon In %s Not supported", reflect.TypeOf(other))
}

func (p Polygon) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: '-polygon' Not supported")
}

func (p Polygon) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not polygon' Not supported")
}

func (p Polygon) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @polygon Not supported")
}

func (p Polygon) Property(ident string) (Value, error) {
	switch ident {
	case "bounds":
		return p.bounds(), nil
	case "vertices":
		values := make([]Value, len(p.Vertices))
		for i, vertex := range p.Vertices {
			values[i] = vertex
		}
		return List{
			Elements: values,
		}, nil
	}
	return baseProperty(p, ident)
}

func (p Polygon) PrintStr() string {
	verts := make([]string, len(p.Vertices))
	for i, vertex := range p.Vertices {
		verts[i] = vertex.PrintStr()
	}
	return fmt.Sprintf("polygon(%s)", strings.Join(verts, ", "))
}

func (p Polygon) iterateHorizLine(x1, y, x2 int, visit func(Value) error) error {
	if x1 > x2 {
		temp := x1
		x1 = x2
		x2 = temp
	}
	for x := x1; x <= x2; x++ {
		if err := visit(Point{x, y}); err != nil {
			return err
		}
	}
	return nil
}

func (p Polygon) bounds() Rect {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for _, pt := range p.Vertices {
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
	return Rect{
		Min: image.Point{minX, minY},
		Max: image.Point{maxX, maxY},
	}
}

func (p Polygon) Iterate(visit func(Value) error) error {
	// buffer to store the x-coordinates of intersections of the polygon with some horizontal line
	xs := make([]int, 0, len(p.Vertices))

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
	// 3. each two x-coordinates In x_intersect are then inside the polygon
	//    (drawhorzlineclip for a pair of two such points)
	//
	for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
		intersectionXs := xs[:0] // re-use xs buffer to avoid allocations

		for i, pt := range p.Vertices {
			previ := i - 1
			if previ < 0 {
				previ = len(p.Vertices) - 1
			}
			prevPt := p.Vertices[previ]
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
				// Add intersection if y crosses the edge (excluding the lower end), or when we are on the lowest line (maxy)
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

	// finally, a special case is Not handled by above algorithm:
	// for two border points with same height miny < y < maxy,
	// sometimes the line between them is Not colored:
	// this happens when the line will be a lower border line of the polygon
	// (eg we are inside the polygon with a smaller y, and outside with a bigger y),
	// So we loop for border lines that are horizontal.
	//
	previ := len(p.Vertices) - 1
	for i, pt := range p.Vertices {
		if (bounds.Min.Y < pt.Y) && p.Vertices[previ].Y == pt.Y && pt.Y <= bounds.Max.Y {
			if err := p.iterateHorizLine(pt.X, pt.Y, p.Vertices[previ].X, visit); err != nil {
				return err
			}
		}
		previ = i
	}
	return nil
}

func (p Polygon) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon[Index] Not supported")
}

func (p Polygon) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon[lower..upper] Not supported")
}

func (p Polygon) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: polygon[%s] Not supported", reflect.TypeOf(index))
}

func (p Polygon) RuntimeTypeName() string {
	return "polygon"
}

func (p Polygon) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: polygon  :: [%s] Not supported", reflect.TypeOf(val))
}
