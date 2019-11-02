package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"image"
	"math"
	"reflect"
)

type point image.Point

func (p point) compare(other value) (value, error) {
	if r, ok := other.(point); ok {
		if p == r {
			return number(0), nil
		}
	}
	return nil, nil
}

func (p point) add(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return point{p.X + int(r+0.5), p.Y + int(r+0.5)}, nil
	case point:
		return point{p.X + r.X, p.Y + r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point + number or point + point, found point + %s", reflect.TypeOf(other))
}

func (p point) sub(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return point{p.X - int(r+0.5), p.Y - int(r+0.5)}, nil
	case point:
		return point{p.X - r.X, p.Y - r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point - number or point - point, found point - %s", reflect.TypeOf(other))
}

func (p point) mul(other value) (value, error) {
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return point{int(lang.Number(p.X)*rn + 0.5), int(lang.Number(p.Y)*rn + 0.5)}, nil
	case point:
		return point{p.X * r.X, p.Y * r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point * number or point * point, found point * %s", reflect.TypeOf(other))
}

func (p point) div(other value) (value, error) {
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return point{int(lang.Number(p.X)/rn + 0.5), int(lang.Number(p.Y)/rn + 0.5)}, nil
	case point:
		return point{p.X / r.X, p.Y / r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point / number or point / point, found point / %s", reflect.TypeOf(other))
}

func (p point) mod(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return point{p.X % int(r+0.5), p.Y % int(r+0.5)}, nil
	case point:
		return point{p.X % r.X, p.Y % r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point %% number or point %% point, found point %% %s", reflect.TypeOf(other))
}

func (p point) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		return boolean(p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: expected point == point, found point == %s", reflect.TypeOf(other))
}

func (p point) neg() (value, error) {
	return point{-p.X, -p.Y}, nil
}

func (p point) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not point not supported")
}

func (p point) at(bitmap BitmapContext) (value, error) {
	return color(bitmap.GetPixel(p.X, p.Y)), nil
}

func (p point) property(ident string) (value, error) {
	switch ident {
	case "x":
		return number(p.X), nil
	case "y":
		return number(p.Y), nil
	case "mag":
		x, y := float64(p.X), float64(p.Y)
		return number(math.Sqrt(x*x + y*y)), nil
	}
	return baseProperty(p, ident)
}

func (p point) printStr() string {
	return fmt.Sprintf("%d;%d", p.X, p.Y)
}

func (p point) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over point")
}

func (p point) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point[index] not supported")
}

func (p point) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point[lower..upper] not supported")
}

func (p point) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: point[%s] not supported", reflect.TypeOf(index))
}

func (p point) runtimeTypeName() string {
	return "point"
}

func (p point) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point  :: [%s] not supported", reflect.TypeOf(val))
}
