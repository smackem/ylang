package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type Color lang.Color

func (c Color) Compare(other Value) (Value, error) {
	if r, ok := other.(Color); ok {
		if c == r {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (c Color) Add(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Color(lang.NewRgba(c.R+rn, c.G+rn, c.B+rn, c.A)), nil
	case Color:
		return Color(lang.NewRgba(c.R+r.R, c.G+r.G, c.B+r.B, c.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: color + %s Not supported", reflect.TypeOf(other))
}

func (c Color) Sub(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Color(lang.NewRgba(c.R-rn, c.G-rn, c.B-rn, c.A)), nil
	case Color:
		return Color(lang.NewRgba(c.R-r.R, c.G-r.G, c.B-r.B, c.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: color - %s Not supported", reflect.TypeOf(other))
}

func (c Color) Mul(other Value) (Value, error) {
	cc := lang.Color(c)
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Color(lang.NewSrgba(cc.ScR()*rn, cc.ScG()*rn, cc.ScB()*rn, cc.ScA())), nil
	case Color:
		rc := lang.Color(r)
		return Color(lang.NewSrgba(cc.ScR()*rc.ScR(), cc.ScG()*rc.ScG(), cc.ScB()*rc.ScB(), cc.ScA())), nil
	}
	return nil, fmt.Errorf("type mismatch: color * %s Not supported", reflect.TypeOf(other))
}

func (c Color) Div(other Value) (Value, error) {
	cc := lang.Color(c)
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Color(lang.NewSrgba(cc.ScR()/rn, cc.ScG()/rn, cc.ScB()/rn, cc.ScA())), nil
	case Color:
		rc := lang.Color(r)
		return Color(lang.NewSrgba(cc.ScR()/rc.ScR(), cc.ScG()/rc.ScG(), cc.ScB()/rc.ScB(), cc.ScA())), nil
	}
	return nil, fmt.Errorf("type mismatch: color / %s Not supported", reflect.TypeOf(other))
}

func (c Color) Mod(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		rn := lang.Number(r)
		return Color(lang.NewRgba(lang.Number(int(c.R+0.5)%int(rn+0.5)), lang.Number(int(c.G+0.5)%int(rn+0.5)), lang.Number(int(c.B+0.5)%int(rn+0.5)), c.A)), nil
	case Color:
		rc := lang.Color(r)
		return Color(lang.NewRgba(lang.Number(int(c.R+0.5)%int(rc.R+0.5)), lang.Number(int(c.G+0.5)%int(rc.G+0.5)), lang.Number(int(c.B+0.5)%int(rc.B+0.5)), c.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: color %% %s Not supported", reflect.TypeOf(other))
}

func (c Color) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'color In' Not supported")
}

func (c Color) Neg() (Value, error) {
	return Color(lang.NewRgba(255-c.R, 255-c.G, 255-c.B, c.A)), nil
}

func (c Color) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not color' Not supported")
}

func (c Color) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @color Not supported")
}

func (c Color) Property(ident string) (Value, error) {
	switch ident {
	case "r", "red":
		return Number(c.R), nil
	case "g", "green":
		return Number(c.G), nil
	case "b", "blue":
		return Number(c.B), nil
	case "a", "alpha":
		return Number(c.A), nil
	case "r01", "red01":
		return Number(lang.Color(c).ScR()), nil
	case "g01", "green01":
		return Number(lang.Color(c).ScG()), nil
	case "b01", "blue01":
		return Number(lang.Color(c).ScB()), nil
	case "a01", "alpha01":
		return Number(lang.Color(c).ScA()), nil
	case "i", "intensity":
		return Number(lang.Color(c).Intensity()), nil
	case "i01", "intensity01":
		return Number(lang.Color(c).ScIntensity()), nil
	}
	return baseProperty(c, ident)
}

func (c Color) PrintStr() string {
	return fmt.Sprintf("rgba(%g,%g,%g:%g)", c.R, c.G, c.B, c.A)
}

func (c Color) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over color")
}

func (c Color) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: color[Index] Not supported")
}

func (c Color) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: color[lower..upper] Not supported")
}

func (c Color) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: color[%s] Not supported", reflect.TypeOf(index))
}

func (c Color) RuntimeTypeName() string {
	return "color"
}

func (c Color) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: color :: [%s] Not supported", reflect.TypeOf(val))
}
