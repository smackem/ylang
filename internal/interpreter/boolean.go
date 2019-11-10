package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type Boolean lang.Boolean

func (b Boolean) Compare(other Value) (Value, error) {
	if r, ok := other.(Boolean); ok {
		if b && !r {
			return Number(1), nil
		}
		if b == r {
			return Number(0), nil
		}
		if !b && r {
			return Number(-1), nil
		}
	}
	return nil, nil
}

func (b Boolean) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool + %s Not supported", reflect.TypeOf(other))
}

func (b Boolean) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool - %s Not supported", reflect.TypeOf(other))
}

func (b Boolean) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool * %s Not supported", reflect.TypeOf(other))
}

func (b Boolean) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool / %s Not supported", reflect.TypeOf(other))
}

func (b Boolean) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool %% %s Not supported", reflect.TypeOf(other))
}

func (b Boolean) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool In %s Not supported", reflect.TypeOf(other))
}

func (b Boolean) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -bool Not supported")
}

func (b Boolean) Not() (Value, error) {
	return !b, nil
}

func (b Boolean) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @bool Not supported")
}

func (b Boolean) Property(ident string) (Value, error) {
	return baseProperty(b, ident)
}

func (b Boolean) PrintStr() string {
	if b {
		return "true"
	}
	return "false"
}

func (b Boolean) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over bool")
}

func (b Boolean) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool[Index] Not supported")
}

func (b Boolean) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool[lower..upper] Not supported")
}

func (b Boolean) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: bool[%s] Not supported", reflect.TypeOf(index))
}

func (b Boolean) RuntimeTypeName() string {
	return "boolean"
}

func (b Boolean) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool :: [%s] Not supported", reflect.TypeOf(val))
}
