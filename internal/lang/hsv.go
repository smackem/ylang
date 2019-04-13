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
