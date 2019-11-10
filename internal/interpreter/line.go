package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"math"
	"reflect"
)

type Line struct {
	Point1 Point
	Point2 Point
}

func (l Line) Compare(other Value) (Value, error) {
	if r, ok := other.(Line); ok {
		if l == r {
			return Number(0), nil
		}
	}
	return Boolean(lang.FalseVal), nil
}

func (l Line) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line + %s Not supported", reflect.TypeOf(other))
}

func (l Line) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line - %s Not supported", reflect.TypeOf(other))
}

func (l Line) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line * %s Not supported", reflect.TypeOf(other))
}

func (l Line) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line / %s Not supported", reflect.TypeOf(other))
}

func (l Line) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line %% %s Not supported", reflect.TypeOf(other))
}

func (l Line) In(other Value) (Value, error) {
	if r, ok := other.(Rect); ok {
		p1, _ := l.Point1.In(r)
		p2, _ := l.Point2.In(r)
		return Boolean(p1.(Boolean) && p2.(Boolean)), nil
	}
	return nil, fmt.Errorf("type mismatch: line In %s Not supported", reflect.TypeOf(other))
}

func (l Line) Neg() (Value, error) {
	return Line{Point1: l.Point2, Point2: l.Point1}, nil
}

func (l Line) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not line' Not supported")
}

func (l Line) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @line Not supported")
}

func (l Line) Property(ident string) (Value, error) {
	switch ident {
	case "p1", "point1":
		return l.Point1, nil
	case "p2", "point2":
		return l.Point2, nil
	case "dx":
		return Number(l.Point2.X - l.Point1.X), nil
	case "dy":
		return Number(l.Point2.Y - l.Point1.Y), nil
	case "len":
		dx, dy := l.Point2.X-l.Point1.X, l.Point2.Y-l.Point1.Y
		return Number(math.Sqrt(float64(dx*dx + dy*dy))), nil
	}
	return baseProperty(l, ident)
}

func (l Line) PrintStr() string {
	return fmt.Sprintf("line(point1:%s, point2:%s)", l.Point1.PrintStr(), l.Point2.PrintStr())
}

func (l Line) Iterate(visit func(Value) error) error {
	dx, dy := lang.Number(l.Point2.X-l.Point1.X), lang.Number(l.Point2.Y-l.Point1.Y)
	dxabs, dyabs := math.Abs(float64(dx)), math.Abs(float64(dy))

	var steps int
	if dxabs > dyabs {
		steps = int(dxabs)
	} else {
		steps = int(dyabs)
	}

	stepsN := lang.Number(steps)
	dx, dy = dx/stepsN, dy/stepsN
	x, y := lang.Number(l.Point1.X), lang.Number(l.Point1.Y)

	for i := 0; i < steps; i++ {
		if err := visit(Point{int(x + 0.5), int(y + 0.5)}); err != nil {
			return err
		}
		x = x + dx
		y = y + dy
	}
	return nil
}

func (l Line) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line[Index] Not supported")
}

func (l Line) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line[lower..upper] Not supported")
}

func (l Line) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: line[%s] Not supported", reflect.TypeOf(index))
}

func (l Line) RuntimeTypeName() string {
	return "line"
}

func (l Line) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: line :: %s Not supported", reflect.TypeOf(val))
}
