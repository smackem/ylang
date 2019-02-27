package lang

import (
	"fmt"
	"reflect"
)

// Number is the number type used by the interpreter
type Number float32

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
	return baseProperty(n, ident)
}

func (n Number) printStr() string {
	return fmt.Sprintf("%f", n)
}

func (n Number) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over number")
}

func (n Number) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: number[index] not supported")
}

// implement sort.Interface for Number slice

type numberSlice []Number

func (p numberSlice) Len() int           { return len(p) }
func (p numberSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p numberSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (n Number) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: number[%s] not supported", reflect.TypeOf(index))
}

func (n Number) runtimeTypeName() string {
	return "number"
}

func (n Number) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: number :: [%s] not supported", reflect.TypeOf(val))
}
