package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type boolean lang.Boolean

func (b boolean) compare(other Value) (Value, error) {
	if r, ok := other.(boolean); ok {
		if b && !r {
			return number(1), nil
		}
		if b == r {
			return number(0), nil
		}
		if !b && r {
			return number(-1), nil
		}
	}
	return nil, nil
}

func (b boolean) add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool + %s not supported", reflect.TypeOf(other))
}

func (b boolean) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool - %s not supported", reflect.TypeOf(other))
}

func (b boolean) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool * %s not supported", reflect.TypeOf(other))
}

func (b boolean) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool / %s not supported", reflect.TypeOf(other))
}

func (b boolean) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool %% %s not supported", reflect.TypeOf(other))
}

func (b boolean) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool in %s not supported", reflect.TypeOf(other))
}

func (b boolean) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -bool not supported")
}

func (b boolean) not() (Value, error) {
	return !b, nil
}

func (b boolean) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @bool not supported")
}

func (b boolean) property(ident string) (Value, error) {
	return baseProperty(b, ident)
}

func (b boolean) printStr() string {
	if b {
		return "true"
	}
	return "false"
}

func (b boolean) iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot iterate over bool")
}

func (b boolean) index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool[index] not supported")
}

func (b boolean) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool[lower..upper] not supported")
}

func (b boolean) indexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: bool[%s] not supported", reflect.TypeOf(index))
}

func (b boolean) runtimeTypeName() string {
	return "boolean"
}

func (b boolean) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: bool :: [%s] not supported", reflect.TypeOf(val))
}
