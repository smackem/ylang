package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type number lang.Number

func (n number) compare(other value) (value, error) {
	if r, ok := other.(number); ok {
		return n - r, nil
	}
	return nil, nil
}

func (n number) add(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return number(n + r), nil
	case point:
		return point{int(n + number(r.X) + 0.5), int(n + number(r.Y) + 0.5)}, nil
	case color:
		nn := lang.Number(n)
		return color(lang.NewRgba(nn+r.R, nn+r.G, nn+r.B, r.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number + number or number + color, found number + %s", reflect.TypeOf(other))
}

func (n number) sub(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return number(n - r), nil
	case point:
		return point{int(n - number(r.X) + 0.5), int(n - number(r.Y) + 0.5)}, nil
	case color:
		nn := lang.Number(n)
		return color(lang.NewRgba(nn-r.R, nn-r.G, nn-r.B, r.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number - number or number - color, found number - %s", reflect.TypeOf(other))
}

func (n number) mul(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return number(n * r), nil
	case point:
		return point{int(n*number(r.X) + 0.5), int(n*number(r.Y) + 0.5)}, nil
	case color:
		rc := lang.Color(r)
		nn := lang.Number(n)
		return color(lang.NewSrgba(nn*rc.ScR(), nn*rc.ScG(), nn*rc.ScB(), rc.ScA())), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number * number or number * color, found number + %s", reflect.TypeOf(other))
}

func (n number) div(other value) (value, error) {
	switch r := other.(type) {
	case number:
		return number(n / r), nil
	case point:
		return point{int(n/number(r.X) + 0.5), int(n/number(r.Y) + 0.5)}, nil
	case color:
		nn := lang.Number(n)
		return color(lang.NewRgba(nn/r.R, nn/r.G, nn/r.B, r.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n number) mod(other value) (value, error) {
	if r, ok := other.(number); ok {
		return number(int(n+0.5) % int(r+0.5)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n number) in(other value) (value, error) {
	if k, ok := other.(kernel); ok {
		for _, kn := range k.values {
			if kn == lang.Number(n) {
				return boolean(true), nil
			}
		}
		return falseVal, nil
	}
	return nil, fmt.Errorf("type mismatch: 'number in %s' not supported", reflect.TypeOf(other))
}

func (n number) neg() (value, error) {
	return number(-n), nil
}

func (n number) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: found 'not number' instead of 'not bool'")
}

func (n number) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: found '@number' instead of '@point'")
}

func (n number) property(ident string) (value, error) {
	return baseProperty(n, ident)
}

func (n number) printStr() string {
	return fmt.Sprintf("%g", n)
}

func (n number) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over number")
}

func (n number) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: number[index] not supported")
}

func (n number) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: number[lower..upper] not supported")
}

// implement sort.Interface for number slice

type numberSlice []number

func (p numberSlice) Len() int           { return len(p) }
func (p numberSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p numberSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (n number) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: number[%s] not supported", reflect.TypeOf(index))
}

func (n number) runtimeTypeName() string {
	return "number"
}

func (n number) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: number :: [%s] not supported", reflect.TypeOf(val))
}
