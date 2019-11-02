package interpreter

import (
	"fmt"
	"image"
	"reflect"
)

type rect image.Rectangle

func (rc rect) compare(other value) (value, error) {
	if r, ok := other.(rect); ok {
		if rc == r {
			return number(0), nil
		}
	}
	return nil, nil
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
		return number(rc.Min.X), nil
	case "y", "top":
		return number(rc.Min.Y), nil
	case "w", "width":
		return number(rc.Max.X - rc.Min.X), nil
	case "h", "height":
		return number(rc.Max.Y - rc.Min.Y), nil
	case "right":
		return number(rc.Max.X), nil
	case "bottom":
		return number(rc.Max.Y), nil
	}
	return baseProperty(rc, ident)
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

func (rc rect) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect[lower..upper] not supported")
}

func (rc rect) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: rect[%s] not supported", reflect.TypeOf(index))
}

func (rc rect) runtimeTypeName() string {
	return "rect"
}

func (rc rect) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect  :: [%s] not supported", reflect.TypeOf(val))
}
