package lang

import (
	"fmt"
	"reflect"
)

type value interface {
	equals(other value) (bool, error)
	greaterThan(other value) (bool, error)
	greaterThanOrEqual(other value) (bool, error)
	lessThan(other value) (bool, error)
	lessThanOrEqual(other value) (bool, error)
	add(other value) (value, error)
	sub(other value) (value, error)
	mul(other value) (value, error)
	div(other value) (value, error)
	mod(other value) (value, error)
	in(other value) (bool, error)
	neg() (value, error)
	not() (value, error)
	at(bitmap Bitmap) (value, error)
	property(ident string) (value, error)
}

//////////////////////////////////////////// Number

func (n Number) equals(other value) (bool, error) {
	if r, ok := other.(Number); ok {
		return n == r, nil
	}
	return false, fmt.Errorf("type mismatch: expected number == number, found number == %s", reflect.TypeOf(other))
}

func (n Number) greaterThan(other value) (bool, error) {
	if r, ok := other.(Number); ok {
		return n > r, nil
	}
	return false, fmt.Errorf("type mismatch: expected number > number, found number > %s", reflect.TypeOf(other))
}

func (n Number) greaterThanOrEqual(other value) (bool, error) {
	if r, ok := other.(Number); ok {
		return n >= r, nil
	}
	return false, fmt.Errorf("type mismatch: expected number >= number, found number >= %s", reflect.TypeOf(other))
}

func (n Number) lessThan(other value) (bool, error) {
	if r, ok := other.(Number); ok {
		return n < r, nil
	}
	return false, fmt.Errorf("type mismatch: expected number < number, found number < %s", reflect.TypeOf(other))
}

func (n Number) lessThanOrEqual(other value) (bool, error) {
	if r, ok := other.(Number); ok {
		return n <= r, nil
	}
	return false, fmt.Errorf("type mismatch: expected number <= number, found number <= %s", reflect.TypeOf(other))
}

func (n Number) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n + r), nil
	case Color:
		return NewRgba(n+r.R, n+r.G, n+r.B, r.A), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number + number or number + color, found number + %s", reflect.TypeOf(other))
}

func (n Number) sub(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Number(n - r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number - number, found number - %s", reflect.TypeOf(other))
}

func (n Number) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n * r), nil
	case Color:
		return NewSrgba(n*r.ScR(), n*r.ScG(), n*r.ScB(), r.ScA()), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number * number or number * color, found number + %s", reflect.TypeOf(other))
}

func (n Number) div(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Number(n / r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n Number) mod(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return Number(n / r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n Number) in(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: 'number in %s' not supported", reflect.TypeOf(other))
}

func (n Number) neg() (value, error) {
	return Number(-n), nil
}

func (n Number) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: found 'not number' instead of 'not bool'")
}

func (n Number) at(bitmap Bitmap) (value, error) {
	return nil, fmt.Errorf("type mismatch: found '@number' instead of '@position'")
}

func (n Number) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'number.%s'", ident)
}

//////////////////////////////////////////// string

func (s String) equals(other value) (bool, error) {
	if r, ok := other.(String); ok {
		return s == r, nil
	}
	return false, fmt.Errorf("type mismatch: expected string == string, found string == %s", reflect.TypeOf(other))
}

func (s String) greaterThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: string > %s not supported", reflect.TypeOf(other))
}

func (s String) greaterThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: string >= %s not supported", reflect.TypeOf(other))
}

func (s String) lessThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: string < %s not supported", reflect.TypeOf(other))
}

func (s String) lessThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: string <= %s not supported", reflect.TypeOf(other))
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

func (s String) in(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: string in %s not supported", reflect.TypeOf(other))
}

func (s String) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -string not supported")
}

func (s String) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not string not supported")
}

func (s String) at(bitmap Bitmap) (value, error) {
	return nil, fmt.Errorf("type mismatch: @string not supported")
}

func (s String) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'string.%s'", ident)
}

//////////////////////////////////////////// lang.Position

func (p Position) equals(other value) (bool, error) {
	if r, ok := other.(Position); ok {
		return p == r, nil
	}
	return false, fmt.Errorf("type mismatch: expected position == position, found position == %s", reflect.TypeOf(other))
}

func (p Position) greaterThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: position > %s not supported", reflect.TypeOf(other))
}

func (p Position) greaterThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: position >= %s not supported", reflect.TypeOf(other))
}

func (p Position) lessThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: position < %s not supported", reflect.TypeOf(other))
}

func (p Position) lessThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: position <= %s not supported", reflect.TypeOf(other))
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

func (p Position) in(other value) (bool, error) {
	if r, ok := other.(Rect); ok {
		return p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y, nil
	}
	return false, fmt.Errorf("type mismatch: expected position == position, found position == %s", reflect.TypeOf(other))
}

func (p Position) neg() (value, error) {
	return Position{-p.X, -p.Y}, nil
}

func (p Position) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not position not supported")
}

func (p Position) at(bitmap Bitmap) (value, error) {
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

//////////////////////////////////////////// lang.Color

func (c Color) equals(other value) (bool, error) {
	if r, ok := other.(Color); ok {
		return c == r, nil
	}
	return false, fmt.Errorf("type mismatch: expected color == color, found color == %s", reflect.TypeOf(other))
}

func (c Color) greaterThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: color > %s not supported", reflect.TypeOf(other))
}

func (c Color) greaterThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: color >= %s not supported", reflect.TypeOf(other))
}

func (c Color) lessThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: color <= %s not supported", reflect.TypeOf(other))
}

func (c Color) lessThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: color <= %s not supported", reflect.TypeOf(other))
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

func (c Color) in(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: 'color in' not supported")
}

func (c Color) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -color not supported")
}

func (c Color) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not color' not supported")
}

func (c Color) at(bitmap Bitmap) (value, error) {
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
	}
	return nil, fmt.Errorf("unknown property 'color.%s'", ident)
}

//////////////////////////////////////////// lang.Rect

func (rect Rect) equals(other value) (bool, error) {
	if r, ok := other.(Rect); ok {
		return rect == r, nil
	}
	return false, fmt.Errorf("type mismatch: rect == %s not supported", reflect.TypeOf(other))
}

func (rect Rect) greaterThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: rect > %s not supported", reflect.TypeOf(other))
}

func (rect Rect) greaterThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: rect >= %s not supported", reflect.TypeOf(other))
}

func (rect Rect) lessThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: rect < %s not supported", reflect.TypeOf(other))
}

func (rect Rect) lessThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: rect <= %s not supported", reflect.TypeOf(other))
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

func (rect Rect) in(other value) (bool, error) {
	if r, ok := other.(Rect); ok {
		return rect.Min.X >= r.Min.X && rect.Min.Y >= r.Min.Y && rect.Max.X < r.Max.X && rect.Max.Y < r.Max.Y, nil
	}
	return false, fmt.Errorf("type mismatch: rect in %s not supported", reflect.TypeOf(other))
}

func (rect Rect) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -rect not supported")
}

func (rect Rect) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not rect' not supported")
}

func (rect Rect) at(bitmap Bitmap) (value, error) {
	return nil, fmt.Errorf("type mismatch: @rect not supported")
}

func (rect Rect) property(ident string) (value, error) {
	switch ident {
	case "x":
		return Number(rect.Min.X), nil
	case "y":
		return Number(rect.Min.Y), nil
	case "w":
		return Number(rect.Max.X - rect.Min.X), nil
	case "h":
		return Number(rect.Max.Y - rect.Min.Y), nil
	}
	return nil, fmt.Errorf("unknown property 'rect.%s'", ident)
}

//////////////////////////////////////////// lang.Kernel(radius, length, values)

type kernel struct {
	length int
	radius int
	values []Number
}

func (k kernel) equals(other value) (bool, error) {
	if r, ok := other.(kernel); ok {
		return reflect.DeepEqual(k, r), nil
	}
	return false, fmt.Errorf("type mismatch: kernel == %s not supported", reflect.TypeOf(other))
}

func (k kernel) greaterThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: kernel > %s not supported", reflect.TypeOf(other))
}

func (k kernel) greaterThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: kernel >= %s not supported", reflect.TypeOf(other))
}

func (k kernel) lessThan(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: kernel < %s not supported", reflect.TypeOf(other))
}

func (k kernel) lessThanOrEqual(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: kernel <= %s not supported", reflect.TypeOf(other))
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

func (k kernel) in(other value) (bool, error) {
	return false, fmt.Errorf("type mismatch: kernel in %s not supported", reflect.TypeOf(other))
}

func (k kernel) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -kernel not supported")
}

func (k kernel) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not kernel' not supported")
}

func (k kernel) at(bitmap Bitmap) (value, error) {
	return nil, fmt.Errorf("type mismatch: @kernel not supported")
}

func (k kernel) property(ident string) (value, error) {
	switch ident {
	case "length":
		return Number(k.length), nil
	case "radius":
		return Number(k.radius), nil
	}
	return nil, fmt.Errorf("unknown property 'kernel.%s'", ident)
}
