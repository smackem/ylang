package lang

import (
	"fmt"
	"image"
	"reflect"
)

type rect image.Rectangle

func (rc rect) equals(other value) (value, error) {
	if r, ok := other.(rect); ok {
		return boolean(rc == r), nil
	}
	return falseVal, nil
}

func (rc rect) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect > %s not supported", reflect.TypeOf(other))
}

func (rc rect) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect >= %s not supported", reflect.TypeOf(other))
}

func (rc rect) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect < %s not supported", reflect.TypeOf(other))
}

func (rc rect) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect <= %s not supported", reflect.TypeOf(other))
}

func (rc rect) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect + %s not supported", reflect.TypeOf(other))
}

func (rc rect) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect - %s not supported", reflect.TypeOf(other))
}

func (rc rect) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect * %s not supported", reflect.TypeOf(other))
}

func (rc rect) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect / %s not supported", reflect.TypeOf(other))
}

func (rc rect) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect %% %s not supported", reflect.TypeOf(other))
}

func (rc rect) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		return boolean(rc.Min.X >= r.Min.X && rc.Min.Y >= r.Min.Y && rc.Max.X < r.Max.X && rc.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: rect in %s not supported", reflect.TypeOf(other))
}

func (rc rect) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -rect not supported")
}

func (rc rect) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not rect' not supported")
}

func (rc rect) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @rect not supported")
}

func (rc rect) property(ident string) (value, error) {
	switch ident {
	case "x", "left":
		return Number(rc.Min.X), nil
	case "y", "top":
		return Number(rc.Min.Y), nil
	case "w", "width":
		return Number(rc.Max.X - rc.Min.X), nil
	case "h", "height":
		return Number(rc.Max.Y - rc.Min.Y), nil
	case "right":
		return Number(rc.Max.X), nil
	case "bottom":
		return Number(rc.Max.Y), nil
	}
	return nil, fmt.Errorf("unknown property 'rect.%s'", ident)
}

func (rc rect) printStr() string {
	return fmt.Sprintf("rect(x:%d, y:%d, w:%d, h:%d)", rc.Min.X, rc.Min.Y, rc.Max.X-rc.Min.X, rc.Max.Y-rc.Min.Y)
}

func (rc rect) iterate(visit func(value) error) error {
	for y := rc.Min.Y; y < rc.Max.Y; y++ {
		for x := rc.Min.X; x < rc.Max.X; x++ {
			if err := visit(point{x, y}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rc rect) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect[index] not supported")
}
