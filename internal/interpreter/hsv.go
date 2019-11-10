package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"math"
	"reflect"
)

type ColorHsv struct {
	H lang.Number
	S lang.Number
	V lang.Number
}

func hsvFromRgb(rgb lang.Color) ColorHsv {
	r := rgb.ScR()
	g := rgb.ScG()
	b := rgb.ScB()

	max := lang.Number(math.Max(float64(r), math.Max(float64(g), float64(b))))
	min := lang.Number(math.Min(float64(r), math.Min(float64(g), float64(b))))

	var hue lang.Number

	if max == min {
		hue = 0
	} else if max == r && g >= b {
		hue = 60.0*(g-b)/(max-min) + 0.0
	} else if max == r && g < b {
		hue = 60.0*(g-b)/(max-min) + 360.0
	} else if max == g {
		hue = 60.0*(b-r)/(max-min) + 120.0
	} else if max == b {
		hue = 60.0*(r-g)/(max-min) + 240.0
	} else {
		hue = 0.0
	}

	var s lang.Number
	if max == 0.0 {
		s = 0.0
	} else {
		s = 1.0 - min/max
	}

	return ColorHsv{hue, s, max}
}

func (hsv ColorHsv) clamp() ColorHsv {
	h := hsv.H
	s := hsv.S
	v := hsv.V

	if h >= 360 {
		h = 360 - math.SmallestNonzeroFloat32
	} else if h < 0 {
		h = 0
	}

	if s > 1 {
		s = 1
	} else if s < 0 {
		s = 0
	}

	if v > 1 {
		v = 1
	} else if v < 0 {
		v = 0
	}

	return ColorHsv{h, s, v}
}

func (hsv ColorHsv) rgb() lang.Color {
	hsv = hsv.clamp()
	h := hsv.H
	s := hsv.S
	v := hsv.V

	hi := int(h) / 60 % 6
	f := h/60.0 - lang.Number(hi)

	r := lang.Number(0.0)
	g := lang.Number(0.0)
	b := lang.Number(0.0)

	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)

	switch hi {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	case 5:
		r = v
		g = p
		b = q
	}

	return lang.NewSrgba(r, g, b, 1.0)
}

func (hsv ColorHsv) Compare(other Value) (Value, error) {
	if r, ok := other.(ColorHsv); ok {
		if hsv == r {
			return Number(0), nil
		}
	}
	return Boolean(lang.FalseVal), nil
}

func (hsv ColorHsv) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv + %s Not supported", reflect.TypeOf(other))
}

func (hsv ColorHsv) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv - %s Not supported", reflect.TypeOf(other))
}

func (hsv ColorHsv) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv * %s Not supported", reflect.TypeOf(other))
}

func (hsv ColorHsv) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv / %s Not supported", reflect.TypeOf(other))
}

func (hsv ColorHsv) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv %% %s Not supported", reflect.TypeOf(other))
}

func (hsv ColorHsv) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv In %s Not supported", reflect.TypeOf(other))
}

func (hsv ColorHsv) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: '-hsv' Not supported")
}

func (hsv ColorHsv) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not hsv' Not supported")
}

func (hsv ColorHsv) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @hsv Not supported")
}

func (hsv ColorHsv) Property(ident string) (Value, error) {
	switch ident {
	case "h", "hue":
		return Number(hsv.H), nil
	case "s", "saturation":
		return Number(hsv.S), nil
	case "v", "value":
		return Number(hsv.V), nil
	}
	return baseProperty(hsv, ident)
}

func (hsv ColorHsv) PrintStr() string {
	return fmt.Sprintf("hsv(h:%s, s:%s, v:%s)", Number(hsv.H).PrintStr(), Number(hsv.S).PrintStr(), Number(hsv.V).PrintStr())
}

func (hsv ColorHsv) Iterate(visit func(Value) error) error {
	return fmt.Errorf("type mismatch: iteration over hsv Not supported")
}

func (hsv ColorHsv) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv[Index] Not supported")
}

func (hsv ColorHsv) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv[lower..upper] Not supported")
}

func (hsv ColorHsv) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: hsv[%s] Not supported", reflect.TypeOf(index))
}

func (hsv ColorHsv) RuntimeTypeName() string {
	return "hsv"
}

func (hsv ColorHsv) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hsv :: %s Not supported", reflect.TypeOf(val))
}
