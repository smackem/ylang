package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type color lang.Color

func (c color) compare(other Value) (Value, error) {
	if r, ok := other.(color); ok {
		if c == r {
			return number(0), nil
		}
	}
	return nil, nil
}

func (c color) add(other Value) (Value, error) {
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return color(lang.NewRgba(c.R+rn, c.G+rn, c.B+rn, c.A)), nil
	case color:
		return color(lang.NewRgba(c.R+r.R, c.G+r.G, c.B+r.B, c.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: color + %s not supported", reflect.TypeOf(other))
}

func (c color) sub(other Value) (Value, error) {
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return color(lang.NewRgba(c.R-rn, c.G-rn, c.B-rn, c.A)), nil
	case color:
		return color(lang.NewRgba(c.R-r.R, c.G-r.G, c.B-r.B, c.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: color - %s not supported", reflect.TypeOf(other))
}

func (c color) mul(other Value) (Value, error) {
	cc := lang.Color(c)
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return color(lang.NewSrgba(cc.ScR()*rn, cc.ScG()*rn, cc.ScB()*rn, cc.ScA())), nil
	case color:
		rc := lang.Color(r)
		return color(lang.NewSrgba(cc.ScR()*rc.ScR(), cc.ScG()*rc.ScG(), cc.ScB()*rc.ScB(), cc.ScA())), nil
	}
	return nil, fmt.Errorf("type mismatch: color * %s not supported", reflect.TypeOf(other))
}

func (c color) div(other Value) (Value, error) {
	cc := lang.Color(c)
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return color(lang.NewSrgba(cc.ScR()/rn, cc.ScG()/rn, cc.ScB()/rn, cc.ScA())), nil
	case color:
		rc := lang.Color(r)
		return color(lang.NewSrgba(cc.ScR()/rc.ScR(), cc.ScG()/rc.ScG(), cc.ScB()/rc.ScB(), cc.ScA())), nil
	}
	return nil, fmt.Errorf("type mismatch: color / %s not supported", reflect.TypeOf(other))
}

func (c color) mod(other Value) (Value, error) {
	switch r := other.(type) {
	case number:
		rn := lang.Number(r)
		return color(lang.NewRgba(lang.Number(int(c.R+0.5)%int(rn+0.5)), lang.Number(int(c.G+0.5)%int(rn+0.5)), lang.Number(int(c.B+0.5)%int(rn+0.5)), c.A)), nil
	case color:
		rc := lang.Color(r)
		return color(lang.NewRgba(lang.Number(int(c.R+0.5)%int(rc.R+0.5)), lang.Number(int(c.G+0.5)%int(rc.G+0.5)), lang.Number(int(c.B+0.5)%int(rc.B+0.5)), c.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: color %% %s not supported", reflect.TypeOf(other))
}

func (c color) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'color in' not supported")
}

func (c color) neg() (Value, error) {
	return color(lang.NewRgba(255-c.R, 255-c.G, 255-c.B, c.A)), nil
}

func (c color) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'not color' not supported")
}

func (c color) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @color not supported")
}

func (c color) property(ident string) (Value, error) {
	switch ident {
	case "r", "red":
		return number(c.R), nil
	case "g", "green":
		return number(c.G), nil
	case "b", "blue":
		return number(c.B), nil
	case "a", "alpha":
		return number(c.A), nil
	case "r01", "red01":
		return number(lang.Color(c).ScR()), nil
	case "g01", "green01":
		return number(lang.Color(c).ScG()), nil
	case "b01", "blue01":
		return number(lang.Color(c).ScB()), nil
	case "a01", "alpha01":
		return number(lang.Color(c).ScA()), nil
	case "i", "intensity":
		return number(lang.Color(c).Intensity()), nil
	case "i01", "intensity01":
		return number(lang.Color(c).ScIntensity()), nil
	}
	return baseProperty(c, ident)
}

func (c color) printStr() string {
	return fmt.Sprintf("rgba(%g,%g,%g:%g)", c.R, c.G, c.B, c.A)
}

func (c color) iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot iterate over color")
}

func (c color) index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: color[index] not supported")
}

func (c color) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: color[lower..upper] not supported")
}

func (c color) indexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: color[%s] not supported", reflect.TypeOf(index))
}

func (c color) runtimeTypeName() string {
	return "color"
}

func (c color) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: color :: [%s] not supported", reflect.TypeOf(val))
}
