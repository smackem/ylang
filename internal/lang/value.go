package lang

import (
	"fmt"
	"image"
	"math"
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
	iterate(visit func(value) error) error
}

var falseVal = boolean(false)

//////////////////////////////////////////// Number

func (n Number) equals(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return boolean(n == r), nil
	}
	return falseVal, nil
}

func (n Number) greaterThan(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return boolean(n > r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number > number, found number > %s", reflect.TypeOf(other))
}

func (n Number) greaterThanOrEqual(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return boolean(n >= r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number >= number, found number >= %s", reflect.TypeOf(other))
}

func (n Number) lessThan(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return boolean(n < r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number < number, found number < %s", reflect.TypeOf(other))
}

func (n Number) lessThanOrEqual(other value) (value, error) {
	if r, ok := other.(Number); ok {
		return boolean(n <= r), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number <= number, found number <= %s", reflect.TypeOf(other))
}

func (n Number) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n + r), nil
	case point:
		return point{int(n + Number(r.X) + 0.5), int(n + Number(r.Y) + 0.5)}, nil
	case Color:
		return NewRgba(n+r.R, n+r.G, n+r.B, r.A), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number + number or number + color, found number + %s", reflect.TypeOf(other))
}

func (n Number) sub(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n - r), nil
	case point:
		return point{int(n - Number(r.X) + 0.5), int(n - Number(r.Y) + 0.5)}, nil
	case Color:
		return NewRgba(n-r.R, n-r.G, n-r.B, r.A), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number - number or number - color, found number - %s", reflect.TypeOf(other))
}

func (n Number) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n * r), nil
	case point:
		return point{int(n*Number(r.X) + 0.5), int(n*Number(r.Y) + 0.5)}, nil
	case Color:
		return NewSrgba(n*r.ScR(), n*r.ScG(), n*r.ScB(), r.ScA()), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number * number or number * color, found number + %s", reflect.TypeOf(other))
}

func (n Number) div(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n / r), nil
	case point:
		return point{int(n/Number(r.X) + 0.5), int(n/Number(r.Y) + 0.5)}, nil
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
	if k, ok := other.(kernel); ok {
		for _, kn := range k.values {
			if kn == n {
				return boolean(true), nil
			}
		}
		return falseVal, nil
	}
	return nil, fmt.Errorf("type mismatch: 'number in %s' not supported", reflect.TypeOf(other))
}

func (n Number) neg() (value, error) {
	return Number(-n), nil
}

func (n Number) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: found 'not number' instead of 'not bool'")
}

func (n Number) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: found '@number' instead of '@point'")
}

func (n Number) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'number.%s'", ident)
}

func (n Number) printStr() string {
	return fmt.Sprintf("%f", n)
}

func (n Number) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over number")
}

//////////////////////////////////////////// str

type str string

func (s str) equals(other value) (value, error) {
	if r, ok := other.(str); ok {
		return boolean(s == r), nil
	}
	return falseVal, nil
}

func (s str) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string > %s not supported", reflect.TypeOf(other))
}

func (s str) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string >= %s not supported", reflect.TypeOf(other))
}

func (s str) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string < %s not supported", reflect.TypeOf(other))
}

func (s str) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string <= %s not supported", reflect.TypeOf(other))
}

func (s str) add(other value) (value, error) {
	return str(fmt.Sprintf("%s%s", s, other)), nil
}

func (s str) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string - %s not supported", reflect.TypeOf(other))
}

func (s str) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string * %s not supported", reflect.TypeOf(other))
}

func (s str) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string / %s not supported", reflect.TypeOf(other))
}

func (s str) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string %% %s not supported", reflect.TypeOf(other))
}

func (s str) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string in %s not supported", reflect.TypeOf(other))
}

func (s str) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -string not supported")
}

func (s str) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not string not supported")
}

func (s str) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @string not supported")
}

func (s str) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'string.%s'", ident)
}

func (s str) printStr() string {
	return string(s)
}

func (s str) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over string")
}

//////////////////////////////////////////// lang.point

type point image.Point

func (p point) equals(other value) (value, error) {
	if r, ok := other.(point); ok {
		return boolean(p == r), nil
	}
	return falseVal, nil
}

func (p point) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point > %s not supported", reflect.TypeOf(other))
}

func (p point) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point >= %s not supported", reflect.TypeOf(other))
}

func (p point) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point < %s not supported", reflect.TypeOf(other))
}

func (p point) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: point <= %s not supported", reflect.TypeOf(other))
}

func (p point) add(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{p.X + int(r+0.5), p.Y + int(r+0.5)}, nil
	case point:
		return point{p.X + r.X, p.Y + r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point + number or point + point, found point + %s", reflect.TypeOf(other))
}

func (p point) sub(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{p.X - int(r+0.5), p.Y - int(r+0.5)}, nil
	case point:
		return point{p.X - r.X, p.Y - r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point - number or point - point, found point - %s", reflect.TypeOf(other))
}

func (p point) mul(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{int(Number(p.X)*r + 0.5), int(Number(p.Y)*r + 0.5)}, nil
	case point:
		return point{p.X * r.X, p.Y * r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point * number or point * point, found point * %s", reflect.TypeOf(other))
}

func (p point) div(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{int(Number(p.X)/r + 0.5), int(Number(p.Y)/r + 0.5)}, nil
	case point:
		return point{p.X / r.X, p.Y / r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point / number or point / point, found point / %s", reflect.TypeOf(other))
}

func (p point) mod(other value) (value, error) {
	switch r := other.(type) {
	case Number:
		return point{p.X % int(r+0.5), p.Y % int(r+0.5)}, nil
	case point:
		return point{p.X % r.X, p.Y % r.Y}, nil
	}
	return nil, fmt.Errorf("type mismatch: expected point %% number or point %% point, found point %% %s", reflect.TypeOf(other))
}

func (p point) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		return boolean(p.X >= r.Min.X && p.X < r.Max.X && p.Y >= r.Min.Y && p.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: expected point == point, found point == %s", reflect.TypeOf(other))
}

func (p point) neg() (value, error) {
	return point{-p.X, -p.Y}, nil
}

func (p point) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not point not supported")
}

func (p point) at(bitmap BitmapContext) (value, error) {
	return bitmap.GetPixel(p.X, p.Y), nil
}

func (p point) property(ident string) (value, error) {
	switch ident {
	case "x":
		return Number(p.X), nil
	case "y":
		return Number(p.Y), nil
	}
	return nil, fmt.Errorf("unknown property 'point.%s'", ident)
}

func (p point) printStr() string {
	return fmt.Sprintf("%d;%d", p.X, p.Y)
}

func (p point) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over point")
}

//////////////////////////////////////////// lang.Color

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
	return nil, fmt.Errorf("unknown property 'color.%s'", ident)
}

func (c Color) printStr() string {
	return fmt.Sprintf("%f:%f:%f@%f", c.R, c.G, c.B, c.A)
}

func (c Color) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over color")
}

//////////////////////////////////////////// lang.rect

type rect image.Rectangle

func (rc rect) equals(other value) (value, error) {
	if r, ok := other.(rect); ok {
		return boolean(rc == r), nil
	}
	return falseVal, nil
}

func (rc rect) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect > %s not supported", reflect.TypeOf(other))
}

func (rc rect) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect >= %s not supported", reflect.TypeOf(other))
}

func (rc rect) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect < %s not supported", reflect.TypeOf(other))
}

func (rc rect) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect <= %s not supported", reflect.TypeOf(other))
}

func (rc rect) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect + %s not supported", reflect.TypeOf(other))
}

func (rc rect) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect - %s not supported", reflect.TypeOf(other))
}

func (rc rect) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect * %s not supported", reflect.TypeOf(other))
}

func (rc rect) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect / %s not supported", reflect.TypeOf(other))
}

func (rc rect) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: rect %% %s not supported", reflect.TypeOf(other))
}

func (rc rect) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		return boolean(rc.Min.X >= r.Min.X && rc.Min.Y >= r.Min.Y && rc.Max.X < r.Max.X && rc.Max.Y < r.Max.Y), nil
	}
	return nil, fmt.Errorf("type mismatch: rect in %s not supported", reflect.TypeOf(other))
}

func (rc rect) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -rect not supported")
}

func (rc rect) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not rect' not supported")
}

func (rc rect) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @rect not supported")
}

func (rc rect) property(ident string) (value, error) {
	switch ident {
	case "x", "left":
		return Number(rc.Min.X), nil
	case "y", "top":
		return Number(rc.Min.Y), nil
	case "w", "width":
		return Number(rc.Max.X - rc.Min.X), nil
	case "h", "height":
		return Number(rc.Max.Y - rc.Min.Y), nil
	case "right":
		return Number(rc.Max.X), nil
	case "bottom":
		return Number(rc.Max.Y), nil
	}
	return nil, fmt.Errorf("unknown property 'rect.%s'", ident)
}

func (rc rect) printStr() string {
	return fmt.Sprintf("rect(x:%d, y:%d, w:%d, h:%d)", rc.Min.X, rc.Min.Y, rc.Max.X-rc.Min.X, rc.Max.Y-rc.Min.Y)
}

func (rc rect) iterate(visit func(value) error) error {
	for y := rc.Min.Y; y < rc.Max.Y; y++ {
		for x := rc.Min.X; x < rc.Max.X; x++ {
			if err := visit(point{x, y}); err != nil {
				return err
			}
		}
	}
	return nil
}

//////////////////////////////////////////// lang.Kernel(radius, length, values)

type kernel struct {
	width  int
	height int
	values []Number
}

func (k kernel) equals(other value) (value, error) {
	if r, ok := other.(kernel); ok {
		return boolean(reflect.DeepEqual(k, r)), nil
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

func (k kernel) iterate(visit func(value) error) error {
	for _, v := range k.values {
		if err := visit(v); err != nil {
			return err
		}
	}
	return nil
}

//////////////////////////////////////////// Bool

type boolean bool

func (b boolean) equals(other value) (value, error) {
	if r, ok := other.(boolean); ok {
		return boolean(b == r), nil
	}
	return falseVal, nil
}

func (b boolean) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool > %s not supported", reflect.TypeOf(other))
}

func (b boolean) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool >= %s not supported", reflect.TypeOf(other))
}

func (b boolean) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool < %s not supported", reflect.TypeOf(other))
}

func (b boolean) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool <= %s not supported", reflect.TypeOf(other))
}

func (b boolean) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool + %s not supported", reflect.TypeOf(other))
}

func (b boolean) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool - %s not supported", reflect.TypeOf(other))
}

func (b boolean) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool * %s not supported", reflect.TypeOf(other))
}

func (b boolean) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool / %s not supported", reflect.TypeOf(other))
}

func (b boolean) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool %% %s not supported", reflect.TypeOf(other))
}

func (b boolean) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool in %s not supported", reflect.TypeOf(other))
}

func (b boolean) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -bool not supported")
}

func (b boolean) not() (value, error) {
	return !b, nil
}

func (b boolean) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @bool not supported")
}

func (b boolean) property(ident string) (value, error) {
	return nil, fmt.Errorf("unknown property 'bool.%s'", ident)
}

func (b boolean) printStr() string {
	if b {
		return "true"
	}
	return "false"
}

func (b boolean) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over bool")
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

func (f function) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over function")
}

//////////////////////////////////////////// line

type line struct {
	point1 point
	point2 point
}

func (l line) equals(other value) (value, error) {
	if r, ok := other.(line); ok {
		return boolean(l == r), nil
	}
	return falseVal, nil
}

func (l line) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line > %s not supported", reflect.TypeOf(other))
}

func (l line) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line >= %s not supported", reflect.TypeOf(other))
}

func (l line) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line < %s not supported", reflect.TypeOf(other))
}

func (l line) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line <= %s not supported", reflect.TypeOf(other))
}

func (l line) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line + %s not supported", reflect.TypeOf(other))
}

func (l line) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line - %s not supported", reflect.TypeOf(other))
}

func (l line) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line * %s not supported", reflect.TypeOf(other))
}

func (l line) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line / %s not supported", reflect.TypeOf(other))
}

func (l line) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: line %% %s not supported", reflect.TypeOf(other))
}

func (l line) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		p1, _ := l.point1.in(r)
		p2, _ := l.point2.in(r)
		return boolean(p1.(boolean) && p2.(boolean)), nil
	}
	return nil, fmt.Errorf("type mismatch: line in %s not supported", reflect.TypeOf(other))
}

func (l line) neg() (value, error) {
	return line{point1: l.point2, point2: l.point1}, nil
}

func (l line) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not line' not supported")
}

func (l line) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @line not supported")
}

func (l line) property(ident string) (value, error) {
	switch ident {
	case "p1", "point1":
		return l.point1, nil
	case "p2", "point2":
		return l.point1, nil
	case "len":
		dx, dy := l.point2.X-l.point1.X, l.point2.Y-l.point1.Y
		return Number(math.Sqrt(float64(dx*dx + dy*dy))), nil
	}
	return nil, fmt.Errorf("unknown property 'rect.%s'", ident)
}

func (l line) printStr() string {
	return fmt.Sprintf("line(point1:%v, point2:%v)", l.point1, l.point2)
}

func (l line) iterate(visit func(value) error) error {
	dx, dy := Number(l.point2.X-l.point1.X), Number(l.point2.Y-l.point1.Y)
	dxabs, dyabs := math.Abs(float64(dx)), math.Abs(float64(dy))

	var steps int
	if dxabs > dyabs {
		steps = int(dxabs)
	} else {
		steps = int(dyabs)
	}

	stepsN := Number(steps)
	dx, dy = dx/stepsN, dy/stepsN
	x, y := Number(l.point1.X), Number(l.point1.Y)

	for i := 0; i < steps; i++ {
		if err := visit(point{int(x + 0.5), int(y + 0.5)}); err != nil {
			return err
		}
		x = x + dx
		y = y + dy
	}
	return nil
}

//////////////////////////////////////////// polygon

type polygon struct {
	vertices []point
}

func (p polygon) equals(other value) (value, error) {
	if r, ok := other.(polygon); ok {
		if len(r.vertices) != len(p.vertices) {
			return falseVal, nil
		}
		for i, v := range p.vertices {
			if v != r.vertices[i] {
				return falseVal, nil
			}
		}
		return boolean(true), nil
	}
	return falseVal, nil
}

func (p polygon) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon > %s not supported", reflect.TypeOf(other))
}

func (p polygon) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon >= %s not supported", reflect.TypeOf(other))
}

func (p polygon) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon < %s not supported", reflect.TypeOf(other))
}

func (p polygon) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon <= %s not supported", reflect.TypeOf(other))
}

func (p polygon) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon + %s not supported", reflect.TypeOf(other))
}

func (p polygon) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon - %s not supported", reflect.TypeOf(other))
}

func (p polygon) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon * %s not supported", reflect.TypeOf(other))
}

func (p polygon) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon / %s not supported", reflect.TypeOf(other))
}

func (p polygon) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: polygon %% %s not supported", reflect.TypeOf(other))
}

func (p polygon) in(other value) (value, error) {
	if r, ok := other.(rect); ok {
		for _, pt := range p.vertices {
			v, err := pt.in(r)
			if err != nil {
				return nil, err
			}
			if !v.(boolean) {
				return falseVal, nil
			}
		}
		return boolean(true), nil
	}
	return nil, fmt.Errorf("type mismatch: polygon in %s not supported", reflect.TypeOf(other))
}

func (p polygon) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: '-polygon' not supported")
}

func (p polygon) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not polygon' not supported")
}

func (p polygon) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @polygon not supported")
}

func (p polygon) property(ident string) (value, error) {
	switch ident {
	case "bounds":
		minX, minY := math.MaxInt32, math.MaxInt32
		maxX, maxY := math.MinInt32, math.MinInt32
		for _, pt := range p.vertices {
			if pt.X < minX {
				minX = pt.X
			}
			if pt.Y < minY {
				minY = pt.Y
			}
			if pt.X > maxX {
				maxX = pt.X
			}
			if pt.Y > maxY {
				maxY = pt.Y
			}
		}
		return rect{
			Min: image.Point{minX, minY},
			Max: image.Point{maxX, maxY},
		}, nil
	}
	return nil, fmt.Errorf("unknown property 'rect.%s'", ident)
}

func (p polygon) printStr() string {
	return fmt.Sprintf("polygon(%v)", p.vertices)
}

/*
static void
draw_fillpoly(SDL_Surface *dst, int *point_x, int *point_y, int num_points,
              Uint32 color)
{
    //point_x : x coordinates of the points
    //point-y : the y coordinates of the points
    //num_points : the number of points
    //
	 int i, i_previous, y;  // i_previous is the index of the point before i
	 int miny, maxy;
	 int x1, y1;
	 int x2, y2;
	 //x_intersect are the x-coordinates of intersections of the polygon
	 //with some horizontal line
	 int *x_intersect = PyMem_New(int, num_points);
	 if (x_intersect == NULL) {
		 PyErr_NoMemory();
		 return;
	 }

	//Determine Y maxima
	 miny = point_y[0];
	 maxy = point_y[0];
	 for (i = 1; (i < num_points); i++) {
		 miny = MIN(miny, point_y[i]);
		 maxy = MAX(maxy, point_y[i]);
	 }

	 if (miny == maxy) {
		 //Special case: polygon only 1 pixel high.

		 //Determine X bounds
		 int minx = point_x[0];
		 int maxx = point_x[0];
		 for (i = 1; (i < num_points); i++) {
			 minx = MIN(minx, point_x[i]);
			 maxx = MAX(maxx, point_x[i]);
		 }
		 drawhorzlineclip(dst, color, minx, miny, maxx);
		 PyMem_Free(x_intersect);
		 return;
	 }

	 //Draw, scanning y
	 //----------------
	 //The algorithm uses a horizontal line (y) that moves from top to the
	 //bottom of the polygon:
	 //1. search intersections with the border lines
	 //2. sort intersections (x_intersect)
	 //3. each two x-coordinates in x_intersect are then inside the polygon
	 //   (drawhorzlineclip for a pair of two such points)
	 //
	 for (y = miny; (y <= maxy); y++) {
		 // n_intersections is the number of intersections with the polygon
		 int n_intersections = 0;
		 for (i = 0; (i < num_points); i++) {
			 i_previous = ((i) ? (i - 1) : (num_points - 1));

			 y1 = point_y[i_previous];
			 y2 = point_y[i];
			 if (y1 < y2) {
				 x1 = point_x[i_previous];
				 x2 = point_x[i];
			 }
			 else if (y1 > y2) {
				 y2 = point_y[i_previous];
				 y1 = point_y[i];
				 x2 = point_x[i_previous];
				 x1 = point_x[i];
			 }
			 else {  // y1 == y2 : has to be handled as special case (below)
				 continue;
			 }
			 if (((y >= y1) && (y < y2)) || ((y == maxy) && (y2 == maxy))) {
				 // add intersection if y crosses the edge (excluding the lower
				 // end), or when we are on the lowest line (maxy)
				 x_intersect[n_intersections++] =
					 (y - y1) * (x2 - x1) / (y2 - y1) + x1;
			 }
		 }
		 qsort(x_intersect, n_intersections, sizeof(int), compare_int);

		 for (i = 0; (i < n_intersections); i += 2) {
			 drawhorzlineclip(dst, color, x_intersect[i], y,
							  x_intersect[i + 1]);
		 }
	 }

	 //Finally, a special case is not handled by above algorithm:
	 //For two border points with same height miny < y < maxy,
	 //sometimes the line between them is not colored:
	 //this happens when the line will be a lower border line of the polygon
	 //(eg we are inside the polygon with a smaller y, and outside with a
	 //bigger y),
	 //So we loop for border lines that are horizontal.
	 //
	 for (i = 0; (i < num_points); i++) {
		 i_previous = ((i) ? (i - 1) : (num_points - 1));
		 y = point_y[i];

		 if ((miny < y) && (point_y[i_previous] == y) && (y < maxy)) {
			 drawhorzlineclip(dst, color, point_x[i], y, point_x[i_previous]);
		 }
	 }
*/

func (p polygon) iterate(visit func(value) error) error {
	return nil
}
