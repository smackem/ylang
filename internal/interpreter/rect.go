package interpreter

import (
	"fmt"
	"image"
	"reflect"
)

type rect image.Rectangle

func (rc rect) compare(other Value) (Value, error) {
	if r, ok := other.(rect); ok {
		if rc == r {
			return number(0), nil
		}
	}
	return nil, nil
}

func (rc rect) add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect + %s not supported", reflect.TypeOf(other))
}

func (rc rect) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect - %s not supported", reflect.TypeOf(other))
}

func (rc rect) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect * %s not supported", reflect.TypeOf(other))
}

func (rc rect) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect / %s not supported", reflect.TypeOf(other))
}

func (rc rect) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect %% %s not supported", reflect.TypeOf(other))
}

func (rc rect) in(other Value) (Value, error) {
	if r, ok := other.(rect); ok {
		return boolean(rc.Min.X >= r.Min.X && rc.Min.Y >= r.Min.Y && rc.Max.X < r.Max.X && rc.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: rect in %s not supported", reflect.TypeOf(other))
}

func (rc rect) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -rect not supported")
}

func (rc rect) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'not rect' not supported")
}

func (rc rect) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @rect not supported")
}

func (rc rect) property(ident string) (Value, error) {
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

func (rc rect) iterate(visit func(Value) error) error {
	for y := rc.Min.Y; y < rc.Max.Y; y++ {
		for x := rc.Min.X; x < rc.Max.X; x++ {
			if err := visit(point{x, y}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rc rect) index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect[index] not supported")
}

func (rc rect) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect[lower..upper] not supported")
}

func (rc rect) indexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: rect[%s] not supported", reflect.TypeOf(index))
}

func (rc rect) runtimeTypeName() string {
	return "rect"
}

func (rc rect) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect  :: [%s] not supported", reflect.TypeOf(val))
}
