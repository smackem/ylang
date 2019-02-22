package lang

import (
	"fmt"
	"reflect"
)

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

func (b boolean) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: bool[index] not supported")
}
