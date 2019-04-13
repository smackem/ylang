package lang

import (
	"fmt"
	"image"
	"reflect"
)

type circle struct {
	center point
	radius Number
}

func (c circle) compare(other value) (value, error) {
	if r, ok := other.(circle); ok {
		if c == r {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (c circle) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle + %s not supported", reflect.TypeOf(other))
}

func (c circle) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle - %s not supported", reflect.TypeOf(other))
}

func (c circle) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle * %s not supported", reflect.TypeOf(other))
}

func (c circle) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle / %s not supported", reflect.TypeOf(other))
}

func (c circle) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle %% %s not supported", reflect.TypeOf(other))
}

func (c circle) bounds() rect {
	radiusInt := int(c.radius + 0.5)
	return rect{
		Min: image.Point{
			X: c.center.X - radiusInt,
			Y: c.center.Y - radiusInt,
		},
		Max: image.Point{
			X: c.center.X + radiusInt,
			Y: c.center.Y + radiusInt,
		},
	}
}

func (c circle) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		rc := c.bounds()
		return boolean(rc.Min.X >= r.Min.X && rc.Min.Y >= r.Min.Y && rc.Max.X < r.Max.X && rc.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: line in %s not supported", reflect.TypeOf(other))
}

func (c circle) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: '-circle' not supported")
}

func (c circle) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not circle' not supported")
}

func (c circle) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @circle not supported")
}

func (c circle) property(ident string) (value, error) {
	switch ident {
	case "center":
		return c.center, nil
	case "radius":
		return c.radius, nil
	case "bounds":
		return c.bounds(), nil
	}
	return baseProperty(c, ident)
}

func (c circle) printStr() string {
	return fmt.Sprintf("circle(center:%s, radius:%s)", c.center.printStr(), c.radius.printStr())
}

func (c circle) iterate(visit func(value) error) error {
	x0, y0 := c.center.X, c.center.Y
	radius := int(c.radius + 0.5)
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0+x; i++ {
			if err := visit(point{i, y0 + y}); err != nil {
				return err
			}
			if err := visit(point{i, y0 - y}); err != nil {
				return err
			}
		}
		for i := x0 - y; i <= x0+y; i++ {
			if err := visit(point{i, y0 + x}); err != nil {
				return err
			}
			if err := visit(point{i, y0 - x}); err != nil {
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

func (c circle) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle[index] not supported")
}

func (c circle) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle[lower..upper] not supported")
}

func (c circle) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: circle[%s] not supported", reflect.TypeOf(index))
}

func (c circle) runtimeTypeName() string {
	return "circle"
}

func (c circle) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: circle :: %s not supported", reflect.TypeOf(val))
}
