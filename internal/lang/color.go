package lang

import (
	"fmt"
	"reflect"
)

// Color represents a color with RGBA channels, each channel value held as a floating-point number
// with range 0..255. This range can be exceeded as a result of a computation.
type Color struct {
	R Number
	G Number
	B Number
	A Number
}

// NewRgba creates a Color from r,g,b,a channel values.
func NewRgba(r Number, g Number, b Number, a Number) Color {
	return Color{
		r,
		g,
		b,
		a,
	}
}

// NewSrgba creates a Color from r,g,b,a channel values which are normalized to the range 0..1.
func NewSrgba(scr Number, scg Number, scb Number, sca Number) Color {
	return Color{
		scr * 255.0,
		scg * 255.0,
		scb * 255.0,
		sca * 255.0,
	}
}

// Clamp returns a Color with all channel values clamped to 0..255
func (c Color) Clamp() Color {
	return Color{
		clamp(c.R),
		clamp(c.G),
		clamp(c.B),
		clamp(c.A),
	}
}

// Intensity returns the brightness of a color normalized to 0..255.
func (c Color) Intensity() Number {
	return 0.299*c.R + 0.587*c.G + 0.114*c.B
}

// ScIntensity returns the brightness of a color normalized to 0..1
func (c Color) ScIntensity() Number {
	return c.Intensity() / 255.0
}

// ScR returns the red channel value normalized to 0..1
func (c Color) ScR() Number {
	return c.R / 255.0
}

// ScG returns the green channel value normalized to 0..1
func (c Color) ScG() Number {
	return c.G / 255.0
}

// ScB returns the blue channel value normalized to 0..1
func (c Color) ScB() Number {
	return c.B / 255.0
}

// ScA returns the alpha value normalized to 0..1
func (c Color) ScA() Number {
	return c.A / 255.0
}

// Implement value

func (c Color) equals(other value) (value, error) {
	if r, ok := other.(Color); ok {
		return boolean(c == r), nil
	}
	return falseVal, nil
}

func (c Color) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color > %s not supported", reflect.TypeOf(other))
}

func (c Color) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color >= %s not supported", reflect.TypeOf(other))
}

func (c Color) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color <= %s not supported", reflect.TypeOf(other))
}

func (c Color) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color <= %s not supported", reflect.TypeOf(other))
}

func (c Color) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return NewRgba(c.R+r, c.G+r, c.B+r, c.A), nil
	case Color:
		return NewRgba(c.R+r.R, c.G+r.G, c.B+r.B, c.A), nil
	}
	return nil, fmt.Errorf("type mismatch: color + %s not supported", reflect.TypeOf(other))
}

func (c Color) sub(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return NewRgba(c.R-r, c.G-r, c.B-r, c.A), nil
	case Color:
		return NewRgba(c.R-r.R, c.G-r.G, c.B-r.B, c.A), nil
	}
	return nil, fmt.Errorf("type mismatch: color - %s not supported", reflect.TypeOf(other))
}

func (c Color) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return NewSrgba(c.ScR()*r, c.ScG()*r, c.ScB()*r, c.ScA()), nil
	case Color:
		return NewSrgba(c.ScR()*r.ScR(), c.ScG()*r.ScG(), c.ScB()*r.ScB(), c.ScA()), nil
	}
	return nil, fmt.Errorf("type mismatch: color * %s not supported", reflect.TypeOf(other))
}

func (c Color) div(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return NewSrgba(c.ScR()/r, c.ScG()/r, c.ScB()/r, c.ScA()), nil
	case Color:
		return NewSrgba(c.ScR()/r.ScR(), c.ScG()/r.ScG(), c.ScB()/r.ScB(), c.ScA()), nil
	}
	return nil, fmt.Errorf("type mismatch: color / %s not supported", reflect.TypeOf(other))
}

func (c Color) mod(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return NewRgba(Number(int(c.R+0.5)%int(r+0.5)), Number(int(c.G+0.5)%int(r+0.5)), Number(int(c.B+0.5)%int(r+0.5)), c.A), nil
	case Color:
		return NewRgba(Number(int(c.R+0.5)%int(r.R+0.5)), Number(int(c.G+0.5)%int(r.G+0.5)), Number(int(c.B+0.5)%int(r.B+0.5)), c.A), nil
	}
	return nil, fmt.Errorf("type mismatch: color %% %s not supported", reflect.TypeOf(other))
}

func (c Color) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: 'color in' not supported")
}

func (c Color) neg() (value, error) {
	return NewRgba(255-c.R, 255-c.G, 255-c.B, c.A), nil
}

func (c Color) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not color' not supported")
}

func (c Color) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @color not supported")
}

func (c Color) property(ident string) (value, error) {
	switch ident {
	case "r":
		return c.R, nil
	case "g":
		return c.G, nil
	case "b":
		return c.B, nil
	case "a":
		return c.A, nil
	case "scr":
		return Number(c.ScR()), nil
	case "scg":
		return Number(c.ScG()), nil
	case "scb":
		return Number(c.ScB()), nil
	case "sca":
		return Number(c.ScA()), nil
	case "i":
		return c.Intensity(), nil
	case "sci":
		return c.ScIntensity(), nil
	}
	return baseProperty(c, ident)
}

func (c Color) printStr() string {
	return fmt.Sprintf("rgba(%g,%g:%g,%g)", c.R, c.G, c.B, c.A)
}

func (c Color) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over color")
}

func (c Color) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color[index] not supported")
}

func (c Color) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color[lower..upper] not supported")
}

func clamp(n Number) Number {
	if n > 255 {
		return 255
	}
	if n < 0 {
		return 0
	}
	return n
}

func (c Color) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: color[%s] not supported", reflect.TypeOf(index))
}

func (c Color) runtimeTypeName() string {
	return "color"
}

func (c Color) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: color :: [%s] not supported", reflect.TypeOf(val))
}
