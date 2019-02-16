package lang

import (
	"fmt"
	"reflect"
)

type value interface {
	equals(other value) (value, error)
	greaterThan(other value) (value, error)
	greaterThanOrEqual(other value) (value, error)
	lessThan(other value) (value, error)
	lessThanOrEqual(other value) (value, error)
	add(other value) (value, error)
	sub(other value) (value, error)
	mul(other value) (value, error)
	div(other value) (value, error)
	mod(other value) (value, error)
	in(other value) (value, error)
	neg() (value, error)
	not() (value, error)
	at(bitmap BitmapContext) (value, error)
	property(ident string) (value, error)
	printStr() string
}

var falseVal = Bool(false)

//////////////////////////////////////////// Number

func (n Number) equals(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Bool(n == r), nil
	}
	return falseVal, nil
}

func (n Number) greaterThan(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Bool(n > r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number > number, found number > %s", reflect.TypeOf(other))
}

func (n Number) greaterThanOrEqual(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Bool(n >= r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number >= number, found number >= %s", reflect.TypeOf(other))
}

func (n Number) lessThan(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Bool(n < r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number < number, found number < %s", reflect.TypeOf(other))
}

func (n Number) lessThanOrEqual(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Bool(n <= r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number <= number, found number <= %s", reflect.TypeOf(other))
}

func (n Number) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n + r), nil
	case Position:
		return Position{int(n + Number(r.X) + 0.5), int(n + Number(r.Y) + 0.5)}, nil
	case Color:
		return NewRgba(n+r.R, n+r.G, n+r.B, r.A), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number + number or number + color, found number + %s", reflect.TypeOf(other))
}

func (n Number) sub(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n - r), nil
	case Position:
		return Position{int(n - Number(r.X) + 0.5), int(n - Number(r.Y) + 0.5)}, nil
	case Color:
		return NewRgba(n-r.R, n-r.G, n-r.B, r.A), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number - number or number - color, found number - %s", reflect.TypeOf(other))
}

func (n Number) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n * r), nil
	case Position:
		return Position{int(n*Number(r.X) + 0.5), int(n*Number(r.Y) + 0.5)}, nil
	case Color:
		return NewSrgba(n*r.ScR(), n*r.ScG(), n*r.ScB(), r.ScA()), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number * number or number * color, found number + %s", reflect.TypeOf(other))
}

func (n Number) div(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n / r), nil
	case Position:
		return Position{int(n/Number(r.X) + 0.5), int(n/Number(r.Y) + 0.5)}, nil
	case Color:
		return NewSrgba(n/r.ScR(), n/r.ScG(), n/r.ScB(), r.ScA()), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n Number) mod(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Number(int(n+0.5) % int(r+0.5)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n Number) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: 'number in %s' not supported", reflect.TypeOf(other))
}

func (n Number) neg() (value, error) {
	return Number(-n), nil
}

func (n Number) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: found 'not number' instead of 'not bool'")
}

func (n Number) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: found '@number' instead of '@position'")
}

func (n Number) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'number.%s'", ident)
}

func (n Number) printStr() string {
	return fmt.Sprintf("%f", n)
}

//////////////////////////////////////////// string

func (s String) equals(other value) (value, error) {
	if r, ok := other.(String); ok {
		return Bool(s == r), nil
	}
	return falseVal, nil
}

func (s String) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string > %s not supported", reflect.TypeOf(other))
}

func (s String) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string >= %s not supported", reflect.TypeOf(other))
}

func (s String) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string < %s not supported", reflect.TypeOf(other))
}

func (s String) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string <= %s not supported", reflect.TypeOf(other))
}

func (s String) add(other value) (value, error) {
	return String(fmt.Sprintf("%s%s", s, other)), nil
}

func (s String) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string - %s not supported", reflect.TypeOf(other))
}

func (s String) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string * %s not supported", reflect.TypeOf(other))
}

func (s String) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string / %s not supported", reflect.TypeOf(other))
}

func (s String) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string %% %s not supported", reflect.TypeOf(other))
}

func (s String) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string in %s not supported", reflect.TypeOf(other))
}

func (s String) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -string not supported")
}

func (s String) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not string not supported")
}

func (s String) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @string not supported")
}

func (s String) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'string.%s'", ident)
}

func (s String) printStr() string {
	return string(s)
}

//////////////////////////////////////////// lang.Position

func (p Position) equals(other value) (value, error) {
	if r, ok := other.(Position); ok {
		return Bool(p == r), nil
	}
	return falseVal, nil
}

func (p Position) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: position > %s not supported", reflect.TypeOf(other))
}

func (p Position) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: position >= %s not supported", reflect.TypeOf(other))
}

func (p Position) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: position < %s not supported", reflect.TypeOf(other))
}

func (p Position) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: position <= %s not supported", reflect.TypeOf(other))
}

func (p Position) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Position{p.X + int(r+0.5), p.Y + int(r+0.5)}, nil
	case Position:
		return Position{p.X + r.X, p.Y + r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected position + number or position + position, found position + %s", reflect.TypeOf(other))
}

func (p Position) sub(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Position{p.X - int(r+0.5), p.Y - int(r+0.5)}, nil
	case Position:
		return Position{p.X - r.X, p.Y - r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected position - number or position - position, found position - %s", reflect.TypeOf(other))
}

func (p Position) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Position{int(Number(p.X)*r + 0.5), int(Number(p.Y)*r + 0.5)}, nil
	case Position:
		return Position{p.X * r.X, p.Y * r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected position * number or position * position, found position * %s", reflect.TypeOf(other))
}

func (p Position) div(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Position{int(Number(p.X)/r + 0.5), int(Number(p.Y)/r + 0.5)}, nil
	case Position:
		return Position{p.X / r.X, p.Y / r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected position / number or position / position, found position / %s", reflect.TypeOf(other))
}

func (p Position) mod(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Position{p.X % int(r+0.5), p.Y % int(r+0.5)}, nil
	case Position:
		return Position{p.X % r.X, p.Y % r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected position %% number or position %% position, found position %% %s", reflect.TypeOf(other))
}

func (p Position) in(other value) (value, error) {
	if r, ok := other.(Rect); ok {
		return Bool(p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: expected position == position, found position == %s", reflect.TypeOf(other))
}

func (p Position) neg() (value, error) {
	return Position{-p.X, -p.Y}, nil
}

func (p Position) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not position not supported")
}

func (p Position) at(bitmap BitmapContext) (value, error) {
	return bitmap.GetPixel(p.X, p.Y), nil
}

func (p Position) property(ident string) (value, error) {
	switch ident {
	case "x":
		return Number(p.X), nil
	case "y":
		return Number(p.Y), nil
	}
	return nil, fmt.Errorf("unknown property 'position.%s'", ident)
}

func (p Position) printStr() string {
	return fmt.Sprintf("%d;%d", p.X, p.Y)
}

//////////////////////////////////////////// lang.Color

func (c Color) equals(other value) (value, error) {
	if r, ok := other.(Color); ok {
		return Bool(c == r), nil
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
	return nil, fmt.Errorf("unknown property 'color.%s'", ident)
}

func (c Color) printStr() string {
	return fmt.Sprintf("%f:%f:%f@%f", c.R, c.G, c.B, c.A)
}

//////////////////////////////////////////// lang.Rect

func (rect Rect) equals(other value) (value, error) {
	if r, ok := other.(Rect); ok {
		return Bool(rect == r), nil
	}
	return falseVal, nil
}

func (rect Rect) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect > %s not supported", reflect.TypeOf(other))
}

func (rect Rect) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect >= %s not supported", reflect.TypeOf(other))
}

func (rect Rect) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect < %s not supported", reflect.TypeOf(other))
}

func (rect Rect) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect <= %s not supported", reflect.TypeOf(other))
}

func (rect Rect) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect + %s not supported", reflect.TypeOf(other))
}

func (rect Rect) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect - %s not supported", reflect.TypeOf(other))
}

func (rect Rect) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect * %s not supported", reflect.TypeOf(other))
}

func (rect Rect) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect / %s not supported", reflect.TypeOf(other))
}

func (rect Rect) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect %% %s not supported", reflect.TypeOf(other))
}

func (rect Rect) in(other value) (value, error) {
	if r, ok := other.(Rect); ok {
		return Bool(rect.Min.X >= r.Min.X && rect.Min.Y >= r.Min.Y && rect.Max.X < r.Max.X && rect.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: rect in %s not supported", reflect.TypeOf(other))
}

func (rect Rect) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -rect not supported")
}

func (rect Rect) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not rect' not supported")
}

func (rect Rect) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @rect not supported")
}

func (rect Rect) property(ident string) (value, error) {
	switch ident {
	case "x", "left":
		return Number(rect.Min.X), nil
	case "y", "top":
		return Number(rect.Min.Y), nil
	case "w":
		return Number(rect.Max.X - rect.Min.X), nil
	case "h":
		return Number(rect.Max.Y - rect.Min.Y), nil
	case "right":
		return Number(rect.Max.X), nil
	case "bottom":
		return Number(rect.Max.Y), nil
	}
	return nil, fmt.Errorf("unknown property 'rect.%s'", ident)
}

func (rect Rect) printStr() string {
	return fmt.Sprintf("rect(x:%d, y:%d, w:%d, h:%d)", rect.Min.X, rect.Min.Y, rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
}

//////////////////////////////////////////// lang.Kernel(radius, length, values)

type kernel struct {
	width  int
	height int
	values []Number
}

func (k kernel) equals(other value) (value, error) {
	if r, ok := other.(kernel); ok {
		return Bool(reflect.DeepEqual(k, r)), nil
	}
	return falseVal, nil
}

func (k kernel) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel > %s not supported", reflect.TypeOf(other))
}

func (k kernel) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel >= %s not supported", reflect.TypeOf(other))
}

func (k kernel) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel < %s not supported", reflect.TypeOf(other))
}

func (k kernel) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel <= %s not supported", reflect.TypeOf(other))
}

func (k kernel) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel + %s not supported", reflect.TypeOf(other))
}

func (k kernel) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel - %s not supported", reflect.TypeOf(other))
}

func (k kernel) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel * %s not supported", reflect.TypeOf(other))
}

func (k kernel) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel / %s not supported", reflect.TypeOf(other))
}

func (k kernel) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel %% %s not supported", reflect.TypeOf(other))
}

func (k kernel) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel in %s not supported", reflect.TypeOf(other))
}

func (k kernel) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -kernel not supported")
}

func (k kernel) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not kernel' not supported")
}

func (k kernel) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @kernel not supported")
}

func (k kernel) property(ident string) (value, error) {
	switch ident {
	case "width":
		return Number(k.width), nil
	case "height":
		return Number(k.height), nil
	case "count":
		return Number(len(k.values)), nil
	}
	return nil, fmt.Errorf("unknown property 'kernel.%s'", ident)
}

func (k kernel) printStr() string {
	return fmt.Sprintf("kernel(width: %d, height: %d)", k.width, k.height)
}

//////////////////////////////////////////// Bool

func (b Bool) equals(other value) (value, error) {
	if r, ok := other.(Bool); ok {
		return Bool(b == r), nil
	}
	return falseVal, nil
}

func (b Bool) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool > %s not supported", reflect.TypeOf(other))
}

func (b Bool) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool >= %s not supported", reflect.TypeOf(other))
}

func (b Bool) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool < %s not supported", reflect.TypeOf(other))
}

func (b Bool) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool <= %s not supported", reflect.TypeOf(other))
}

func (b Bool) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool + %s not supported", reflect.TypeOf(other))
}

func (b Bool) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool - %s not supported", reflect.TypeOf(other))
}

func (b Bool) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool * %s not supported", reflect.TypeOf(other))
}

func (b Bool) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool / %s not supported", reflect.TypeOf(other))
}

func (b Bool) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool %% %s not supported", reflect.TypeOf(other))
}

func (b Bool) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool in %s not supported", reflect.TypeOf(other))
}

func (b Bool) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -bool not supported")
}

func (b Bool) not() (value, error) {
	return !b, nil
}

func (b Bool) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @bool not supported")
}

func (b Bool) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'bool.%s'", ident)
}

func (b Bool) printStr() string {
	if b {
		return "true"
	}
	return "false"
}

//////////////////////////////////////////// functionExpr

type function struct {
	parameterNames []string
	body           []statement
	closure        []scope
}

func (f function) equals(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function == %s not supported", reflect.TypeOf(other))
}

func (f function) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function > %s not supported", reflect.TypeOf(other))
}

func (f function) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function >= %s not supported", reflect.TypeOf(other))
}

func (f function) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function < %s not supported", reflect.TypeOf(other))
}

func (f function) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function <= %s not supported", reflect.TypeOf(other))
}

func (f function) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function + %s not supported", reflect.TypeOf(other))
}

func (f function) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function - %s not supported", reflect.TypeOf(other))
}

func (f function) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function * %s not supported", reflect.TypeOf(other))
}

func (f function) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function / %s not supported", reflect.TypeOf(other))
}

func (f function) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function %% %s not supported", reflect.TypeOf(other))
}

func (f function) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function in %s not supported", reflect.TypeOf(other))
}

func (f function) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -function not supported")
}

func (f function) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not function not supported")
}

func (f function) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @function not supported")
}

func (f function) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'function.%s'", ident)
}

func (f function) printStr() string {
	return fmt.Sprintf("fn(%v) {...}", f.parameterNames)
}
