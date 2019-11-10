package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type Nilval lang.Nil

func (n Nilval) Compare(other Value) (Value, error) {
	if _, ok := other.(Nilval); ok {
		return Number(0), nil
	}
	return nil, nil
}

func (n Nilval) greaterThan(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil > %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) greaterThanOrEqual(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil >= %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) lessThan(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil < %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) lessThanOrEqual(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil <= %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil + %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil - %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil * %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil / %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil %% %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil In %s Not supported", reflect.TypeOf(other))
}

func (n Nilval) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -nil Not supported")
}

func (n Nilval) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: Not nil Not supported")
}

func (n Nilval) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @nil Not supported")
}

func (n Nilval) Property(ident string) (Value, error) {
	return baseProperty(n, ident)
}

func (n Nilval) PrintStr() string {
	return "nil"
}

func (n Nilval) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over nil")
}

func (n Nilval) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil[Index] Not supported")
}

func (n Nilval) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil[lower..upper] Not supported")
}

func (n Nilval) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: nil[%s] Not supported", reflect.TypeOf(index))
}

func (n Nilval) RuntimeTypeName() string {
	return "nil"
}

func (n Nilval) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil :: [%s] Not supported", reflect.TypeOf(val))
}
