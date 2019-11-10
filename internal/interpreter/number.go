package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type Number lang.Number

func (n Number) Compare(other Value) (Value, error) {
	if r, ok := other.(Number); ok {
		return n - r, nil
	}
	return nil, nil
}

func (n Number) Add(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n + r), nil
	case Point:
		return Point{int(n + Number(r.X) + 0.5), int(n + Number(r.Y) + 0.5)}, nil
	case Color:
		nn := lang.Number(n)
		return Color(lang.NewRgba(nn+r.R, nn+r.G, nn+r.B, r.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number + number or number + color, found number + %s", reflect.TypeOf(other))
}

func (n Number) Sub(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n - r), nil
	case Point:
		return Point{int(n - Number(r.X) + 0.5), int(n - Number(r.Y) + 0.5)}, nil
	case Color:
		nn := lang.Number(n)
		return Color(lang.NewRgba(nn-r.R, nn-r.G, nn-r.B, r.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number - number or number - color, found number - %s", reflect.TypeOf(other))
}

func (n Number) Mul(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n * r), nil
	case Point:
		return Point{int(n*Number(r.X) + 0.5), int(n*Number(r.Y) + 0.5)}, nil
	case Color:
		rc := lang.Color(r)
		nn := lang.Number(n)
		return Color(lang.NewSrgba(nn*rc.ScR(), nn*rc.ScG(), nn*rc.ScB(), rc.ScA())), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number * number or number * color, found number + %s", reflect.TypeOf(other))
}

func (n Number) Div(other Value) (Value, error) {
	switch r := other.(type) {
	case Number:
		return Number(n / r), nil
	case Point:
		return Point{int(n/Number(r.X) + 0.5), int(n/Number(r.Y) + 0.5)}, nil
	case Color:
		nn := lang.Number(n)
		return Color(lang.NewRgba(nn/r.R, nn/r.G, nn/r.B, r.A)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n Number) Mod(other Value) (Value, error) {
	if r, ok := other.(Number); ok {
		return Number(int(n+0.5) % int(r+0.5)), nil
	}
	return nil, fmt.Errorf("type mismatch: expected number / number, found number / %s", reflect.TypeOf(other))
}

func (n Number) In(other Value) (Value, error) {
	if k, ok := other.(Kernel); ok {
		for _, kn := range k.Values {
			if kn == lang.Number(n) {
				return Boolean(true), nil
			}
		}
		return falseVal, nil
	}
	return nil, fmt.Errorf("type mismatch: 'number In %s' Not supported", reflect.TypeOf(other))
}

func (n Number) Neg() (Value, error) {
	return Number(-n), nil
}

func (n Number) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: found 'Not number' instead of 'Not bool'")
}

func (n Number) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: found '@number' instead of '@point'")
}

func (n Number) Property(ident string) (Value, error) {
	return baseProperty(n, ident)
}

func (n Number) PrintStr() string {
	return fmt.Sprintf("%g", n)
}

func (n Number) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over number")
}

func (n Number) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: number[Index] Not supported")
}

func (n Number) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: number[lower..upper] Not supported")
}

// implement sort.Interface for number slice

type numberSlice []Number

func (p numberSlice) Len() int           { return len(p) }
func (p numberSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p numberSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (n Number) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: number[%s] Not supported", reflect.TypeOf(index))
}

func (n Number) RuntimeTypeName() string {
	return "number"
}

func (n Number) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: number :: [%s] Not supported", reflect.TypeOf(val))
}
