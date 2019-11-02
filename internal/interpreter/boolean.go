package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type boolean lang.Boolean

func (b boolean) compare(other value) (value, error) {
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
	return baseProperty(b, ident)
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

func (b boolean) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool[index] not supported")
}

func (b boolean) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool[lower..upper] not supported")
}

func (b boolean) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: bool[%s] not supported", reflect.TypeOf(index))
}

func (b boolean) runtimeTypeName() string {
	return "boolean"
}

func (b boolean) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool :: [%s] not supported", reflect.TypeOf(val))
}
