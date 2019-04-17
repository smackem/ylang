package lang

import (
	"fmt"
	"math"
	"reflect"
)

type colorHsv struct {
	h Number
	s Number
	v Number
}

func hsvFromRgb(rgb Color) colorHsv {
	r := rgb.ScR()
	g := rgb.ScG()
	b := rgb.ScB()

	max := Number(math.Max(float64(r), math.Max(float64(g), float64(b))))
	min := Number(math.Min(float64(r), math.Min(float64(g), float64(b))))

	var hue Number

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

	var s Number
	if max == 0.0 {
		s = 0.0
	} else {
		s = 1.0 - min/max
	}

	return colorHsv{hue, s, max}
}

func (hsv colorHsv) clamp() colorHsv {
	h := hsv.h
	s := hsv.s
	v := hsv.v

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

	return colorHsv{h, s, v}
}

func (hsv colorHsv) rgb() Color {
	hsv = hsv.clamp()
	h := hsv.h
	s := hsv.s
	v := hsv.v

	hi := int(h) / 60 % 6
	f := h/60.0 - Number(hi)

	r := Number(0.0)
	g := Number(0.0)
	b := Number(0.0)

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

	return NewSrgba(r, g, b, 1.0)
}

func (hsv colorHsv) compare(other value) (value, error) {
	if r, ok := other.(colorHsv); ok {
		if hsv == r {
			return Number(0), nil
		}
	}
	return falseVal, nil
}

func (hsv colorHsv) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv + %s not supported", reflect.TypeOf(other))
}

func (hsv colorHsv) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv - %s not supported", reflect.TypeOf(other))
}

func (hsv colorHsv) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv * %s not supported", reflect.TypeOf(other))
}

func (hsv colorHsv) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv / %s not supported", reflect.TypeOf(other))
}

func (hsv colorHsv) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv %% %s not supported", reflect.TypeOf(other))
}

func (hsv colorHsv) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv in %s not supported", reflect.TypeOf(other))
}

func (hsv colorHsv) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: '-hsv' not supported")
}

func (hsv colorHsv) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not hsv' not supported")
}

func (hsv colorHsv) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @hsv not supported")
}

func (hsv colorHsv) property(ident string) (value, error) {
	switch ident {
	case "h", "hue":
		return hsv.h, nil
	case "s", "saturation":
		return hsv.s, nil
	case "v", "value":
		return hsv.v, nil
	}
	return baseProperty(hsv, ident)
}

func (hsv colorHsv) printStr() string {
	return fmt.Sprintf("hsv(h:%s, s:%s, v:%s)", hsv.h.printStr(), hsv.s.printStr(), hsv.v.printStr())
}

func (hsv colorHsv) iterate(visit func(value) error) error {
	return fmt.Errorf("type mismatch: iteration over hsv not supported")
}

func (hsv colorHsv) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv[index] not supported")
}

func (hsv colorHsv) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv[lower..upper] not supported")
}

func (hsv colorHsv) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: hsv[%s] not supported", reflect.TypeOf(index))
}

func (hsv colorHsv) runtimeTypeName() string {
	return "hsv"
}

func (hsv colorHsv) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hsv :: %s not supported", reflect.TypeOf(val))
}
