package lang

import (
	"image"
)

// Compile compiles the given source code into a executable Program.
func Compile(src string) (Program, error) {
	tokens, err := lex(src)
	if err != nil {
		return Program{}, err
	}
	prog, err := parse(tokens)
	if err != nil {
		return Program{}, err
	}
	return prog, nil
}

// Program is the compiled, executable ylang program.
type Program struct {
	stmts []statement
}

// Execute executes the Program against the specified Bitmap.
func (prog Program) Execute(bitmap Bitmap) error {
	return interpret(prog, bitmap)
}

// Number is the number type used by the interpreter
type Number float64

// String is the string type used by the interpreter
type String string

// Position is a x,y position in an image
type Position image.Point

// Rect holds a rectangle
type Rect image.Rectangle

// Bool is the boolean type used by the interpreter
type Bool bool

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
	cc := c.Clamp()
	return (0.299*cc.R + 0.587*cc.G + 0.114*cc.B) / 255.0
}

// ScIntensity returns the brightness of a color normalized to 0..1
func (c Color) ScIntensity() Number {
	cc := c.Clamp()
	return (0.299*cc.R + 0.587*cc.G + 0.114*cc.B) / (255.0 * 255.0)
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

// Bitmap is the surface a Program works on.
type Bitmap interface {
	GetPixel(x int, y int) Color
	SetPixel(x int, y int, color Color)
	Width() int
	Height() int
	Convolute(x int, y int, radius int, width int, kernel []Number) Color
	Blt(rect Rect)
}

///////////////////////////////////////////////////////////////////////////////

func clamp(n Number) Number {
	if n > 255 {
		return 255
	}
	if n < 0 {
		return 0
	}
	return n
}
