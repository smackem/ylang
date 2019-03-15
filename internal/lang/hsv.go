package lang

import "math"

type colorHsv struct {
	H Number
	S Number
	V Number
}

func (hsv colorHsv) clamp() colorHsv {
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

	return colorHsv{h, s, v}
}

