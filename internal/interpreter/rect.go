package interpreter

import (
	"fmt"
	"image"
	"reflect"
)

type Rect image.Rectangle

func (rc Rect) Compare(other Value) (Value, error) {
	if r, ok := other.(Rect); ok {
		if rc == r {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (rc Rect) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect + %s Not supported", reflect.TypeOf(other))
}

func (rc Rect) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect - %s Not supported", reflect.TypeOf(other))
}

func (rc Rect) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect * %s Not supported", reflect.TypeOf(other))
}

func (rc Rect) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect / %s Not supported", reflect.TypeOf(other))
}

func (rc Rect) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect %% %s Not supported", reflect.TypeOf(other))
}

func (rc Rect) In(other Value) (Value, error) {
	if r, ok := other.(Rect); ok {
		return Boolean(rc.Min.X >= r.Min.X && rc.Min.Y >= r.Min.Y && rc.Max.X < r.Max.X && rc.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: rect In %s Not supported", reflect.TypeOf(other))
}

func (rc Rect) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -rect Not supported")
}

func (rc Rect) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not rect' Not supported")
}

func (rc Rect) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @rect Not supported")
}

func (rc Rect) Property(ident string) (Value, error) {
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
	return baseProperty(rc, ident)
}

func (rc Rect) PrintStr() string {
	return fmt.Sprintf("rect(x:%d, y:%d, w:%d, h:%d)", rc.Min.X, rc.Min.Y, rc.Max.X-rc.Min.X, rc.Max.Y-rc.Min.Y)
}

func (rc Rect) Iterate(visit func(Value) error) error {
	for y := rc.Min.Y; y < rc.Max.Y; y++ {
		for x := rc.Min.X; x < rc.Max.X; x++ {
			if err := visit(Point{x, y}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rc Rect) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect[Index] Not supported")
}

func (rc Rect) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect[lower..upper] Not supported")
}

func (rc Rect) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: rect[%s] Not supported", reflect.TypeOf(index))
}

func (rc Rect) RuntimeTypeName() string {
	return "rect"
}

func (rc Rect) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: rect  :: [%s] Not supported", reflect.TypeOf(val))
}
