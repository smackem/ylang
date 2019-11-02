package lang

import (
	"math"
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
		c.R.Clamp(),
		c.G.Clamp(),
		c.B.Clamp(),
		c.A.Clamp(),
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

// Number is the number type used by the interpreter
type Number float32

// Number constants
const (
	MinNumber Number = Number(math.MinInt64)
	MaxNumber Number = Number(math.MaxInt64)
)

// Clamp clamps a number to the range of 0..255
func (n Number) Clamp() Number {
	if n > 255 {
		return 255
	}
	if n < 0 {
		return 0
	}
	return n
}

// Boolean is the type used to express booleans
type Boolean bool

// FalseVal is the Boolean false
var FalseVal Boolean = false

// TrueVal is the Boolean false
var TrueVal Boolean = true

// Str is used to express strings
type Str string

// Nil is the type used to express the absence of a value
type Nil struct{}

// NilVal is the only value a Nil type can have
var NilVal Nil = Nil{}
