package lang

import (
	"fmt"
	"image"
	"reflect"
)

type point image.Point

func (p point) equals(other value) (value, error) {
	if r, ok := other.(point); ok {
		return boolean(p == r), nil
	}
	return falseVal, nil
}

func (p point) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point > %s not supported", reflect.TypeOf(other))
}

func (p point) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point >= %s not supported", reflect.TypeOf(other))
}

func (p point) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point < %s not supported", reflect.TypeOf(other))
}

func (p point) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point <= %s not supported", reflect.TypeOf(other))
}

func (p point) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{p.X + int(r+0.5), p.Y + int(r+0.5)}, nil
	case point:
		return point{p.X + r.X, p.Y + r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point + number or point + point, found point + %s", reflect.TypeOf(other))
}

func (p point) sub(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{p.X - int(r+0.5), p.Y - int(r+0.5)}, nil
	case point:
		return point{p.X - r.X, p.Y - r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point - number or point - point, found point - %s", reflect.TypeOf(other))
}

func (p point) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{int(Number(p.X)*r + 0.5), int(Number(p.Y)*r + 0.5)}, nil
	case point:
		return point{p.X * r.X, p.Y * r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point * number or point * point, found point * %s", reflect.TypeOf(other))
}

func (p point) div(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{int(Number(p.X)/r + 0.5), int(Number(p.Y)/r + 0.5)}, nil
	case point:
		return point{p.X / r.X, p.Y / r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point / number or point / point, found point / %s", reflect.TypeOf(other))
}

func (p point) mod(other value) (value, error) {
	switch r := other.(type) {
	case Number:
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
	return bitmap.GetPixel(p.X, p.Y), nil
}

func (p point) property(ident string) (value, error) {
	switch ident {
	case "x":
		return Number(p.X), nil
	case "y":
		return Number(p.Y), nil
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

func (p point) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: point[%s] not supported", reflect.TypeOf(index))
}

func (p point) runtimeTypeName() string {
	return "point"
}

func (p point) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point  :: [%s] not supported", reflect.TypeOf(val))
}
