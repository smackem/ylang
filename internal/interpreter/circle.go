package interpreter

import (
	"fmt"
	"image"
	"reflect"
)

type Circle struct {
	Center Point
	Radius Number
}

func (c Circle) Compare(other Value) (Value, error) {
	if r, ok := other.(Circle); ok {
		if c == r {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (c Circle) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle + %s Not supported", reflect.TypeOf(other))
}

func (c Circle) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle - %s Not supported", reflect.TypeOf(other))
}

func (c Circle) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle * %s Not supported", reflect.TypeOf(other))
}

func (c Circle) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle / %s Not supported", reflect.TypeOf(other))
}

func (c Circle) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle %% %s Not supported", reflect.TypeOf(other))
}

func (c Circle) bounds() Rect {
	radiusInt := int(c.Radius + 0.5)
	return Rect{
		Min: image.Point{
			X: c.Center.X - radiusInt,
			Y: c.Center.Y - radiusInt,
		},
		Max: image.Point{
			X: c.Center.X + radiusInt,
			Y: c.Center.Y + radiusInt,
		},
	}
}

func (c Circle) In(other Value) (Value, error) {
	if r, ok := other.(Rect); ok {
		rc := c.bounds()
		return Boolean(rc.Min.X >= r.Min.X && rc.Min.Y >= r.Min.Y && rc.Max.X < r.Max.X && rc.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: line In %s Not supported", reflect.TypeOf(other))
}

func (c Circle) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: '-circle' Not supported")
}

func (c Circle) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not circle' Not supported")
}

func (c Circle) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @circle Not supported")
}

func (c Circle) Property(ident string) (Value, error) {
	switch ident {
	case "center":
		return c.Center, nil
	case "radius":
		return c.Radius, nil
	case "bounds":
		return c.bounds(), nil
	}
	return baseProperty(c, ident)
}

func (c Circle) PrintStr() string {
	return fmt.Sprintf("circle(center:%s, radius:%s)", c.Center.PrintStr(), c.Radius.PrintStr())
}

func (c Circle) Iterate(visit func(Value) error) error {
	x0, y0 := c.Center.X, c.Center.Y
	radius := int(c.Radius + 0.5)
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0+x; i++ {
			if err := visit(Point{i, y0 + y}); err != nil {
				return err
			}
			if err := visit(Point{i, y0 - y}); err != nil {
				return err
			}
		}
		for i := x0 - y; i <= x0+y; i++ {
			if err := visit(Point{i, y0 + x}); err != nil {
				return err
			}
			if err := visit(Point{i, y0 - x}); err != nil {
				return err
			}
		}

		y++
		radiusError += yChange
		yChange += 2
		if (radiusError<<1)+xChange > 0 {
			x--
			radiusError += xChange
			xChange += 2
		}
	}
	return nil
}

func (c Circle) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle[Index] Not supported")
}

func (c Circle) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle[lower..upper] Not supported")
}

func (c Circle) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: circle[%s] Not supported", reflect.TypeOf(index))
}

func (c Circle) RuntimeTypeName() string {
	return "circle"
}

func (c Circle) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: circle :: %s Not supported", reflect.TypeOf(val))
}
