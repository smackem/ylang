package lang

import (
	"fmt"
	"math"
	"reflect"
)

type line struct {
	point1 point
	point2 point
}

func (l line) equals(other value) (value, error) {
	if r, ok := other.(line); ok {
		return boolean(l == r), nil
	}
	return falseVal, nil
}

func (l line) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line > %s not supported", reflect.TypeOf(other))
}

func (l line) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line >= %s not supported", reflect.TypeOf(other))
}

func (l line) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line < %s not supported", reflect.TypeOf(other))
}

func (l line) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line <= %s not supported", reflect.TypeOf(other))
}

func (l line) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line + %s not supported", reflect.TypeOf(other))
}

func (l line) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line - %s not supported", reflect.TypeOf(other))
}

func (l line) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line * %s not supported", reflect.TypeOf(other))
}

func (l line) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line / %s not supported", reflect.TypeOf(other))
}

func (l line) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line %% %s not supported", reflect.TypeOf(other))
}

func (l line) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		p1, _ := l.point1.in(r)
		p2, _ := l.point2.in(r)
		return boolean(p1.(boolean) && p2.(boolean)), nil
	}
	return nil, fmt.Errorf("type mismatch: line in %s not supported", reflect.TypeOf(other))
}

func (l line) neg() (value, error) {
	return line{point1: l.point2, point2: l.point1}, nil
}

func (l line) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not line' not supported")
}

func (l line) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @line not supported")
}

func (l line) property(ident string) (value, error) {
	switch ident {
	case "p1", "point1":
		return l.point1, nil
	case "p2", "point2":
		return l.point2, nil
	case "dx":
		return Number(l.point2.X - l.point1.X), nil
	case "dy":
		return Number(l.point2.Y - l.point1.Y), nil
	case "len":
		dx, dy := l.point2.X-l.point1.X, l.point2.Y-l.point1.Y
		return Number(math.Sqrt(float64(dx*dx + dy*dy))), nil
	}
	return baseProperty(l, ident)
}

func (l line) printStr() string {
	return fmt.Sprintf("line(point1:%s, point2:%s)", l.point1.printStr(), l.point2.printStr())
}

func (l line) iterate(visit func(value) error) error {
	dx, dy := Number(l.point2.X-l.point1.X), Number(l.point2.Y-l.point1.Y)
	dxabs, dyabs := math.Abs(float64(dx)), math.Abs(float64(dy))

	var steps int
	if dxabs > dyabs {
		steps = int(dxabs)
	} else {
		steps = int(dyabs)
	}

	stepsN := Number(steps)
	dx, dy = dx/stepsN, dy/stepsN
	x, y := Number(l.point1.X), Number(l.point1.Y)

	for i := 0; i < steps; i++ {
		if err := visit(point{int(x + 0.5), int(y + 0.5)}); err != nil {
			return err
		}
		x = x + dx
		y = y + dy
	}
	return nil
}

func (l line) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line[index] not supported")
}

func (l line) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line[lower..upper] not supported")
}

func (l line) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: line[%s] not supported", reflect.TypeOf(index))
}

func (l line) runtimeTypeName() string {
	return "line"
}

func (l line) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line :: %s not supported", reflect.TypeOf(val))
}
