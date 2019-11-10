package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"image"
	"math"
	"reflect"
)

type Point image.Point

func (p Point) Compare(other Value) (Value, error) {
	if r, ok := other.(Point); ok {
		if p == r {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (p Point) Add(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Point{p.X + int(r+0.5), p.Y + int(r+0.5)}, nil
	case Point:
		return Point{p.X + r.X, p.Y + r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point + number or point + point, found point + %s", reflect.TypeOf(other))
}

func (p Point) Sub(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Point{p.X - int(r+0.5), p.Y - int(r+0.5)}, nil
	case Point:
		return Point{p.X - r.X, p.Y - r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point - number or point - point, found point - %s", reflect.TypeOf(other))
}

func (p Point) Mul(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Point{int(lang.Number(p.X)*rn + 0.5), int(lang.Number(p.Y)*rn + 0.5)}, nil
	case Point:
		return Point{p.X * r.X, p.Y * r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point * number or point * point, found point * %s", reflect.TypeOf(other))
}

func (p Point) Div(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Point{int(lang.Number(p.X)/rn + 0.5), int(lang.Number(p.Y)/rn + 0.5)}, nil
	case Point:
		return Point{p.X / r.X, p.Y / r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point / number or point / point, found point / %s", reflect.TypeOf(other))
}

func (p Point) Mod(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Point{p.X % int(r+0.5), p.Y % int(r+0.5)}, nil
	case Point:
		return Point{p.X % r.X, p.Y % r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point %% number or point %% point, found point %% %s", reflect.TypeOf(other))
}

func (p Point) In(other Value) (Value, error) {
	if r, ok := other.(Rect); ok {
		return Boolean(p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: expected point == point, found point == %s", reflect.TypeOf(other))
}

func (p Point) Neg() (Value, error) {
	return Point{-p.X, -p.Y}, nil
}

func (p Point) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: Not point Not supported")
}

func (p Point) At(bitmap BitmapContext) (Value, error) {
	return Color(bitmap.GetPixel(p.X, p.Y)), nil
}

func (p Point) Property(ident string) (Value, error) {
	switch ident {
	case "x":
		return Number(p.X), nil
	case "y":
		return Number(p.Y), nil
	case "mag":
		x, y := float64(p.X), float64(p.Y)
		return Number(math.Sqrt(x*x + y*y)), nil
	}
	return baseProperty(p, ident)
}

func (p Point) PrintStr() string {
	return fmt.Sprintf("%d;%d", p.X, p.Y)
}

func (p Point) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over point")
}

func (p Point) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: point[Index] Not supported")
}

func (p Point) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: point[lower..upper] Not supported")
}

func (p Point) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: point[%s] Not supported", reflect.TypeOf(index))
}

func (p Point) RuntimeTypeName() string {
	return "point"
}

func (p Point) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: point  :: [%s] Not supported", reflect.TypeOf(val))
}
