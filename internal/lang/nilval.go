package lang

import (
	"fmt"
	"reflect"
)

type nilval struct{}

func (n nilval) compare(other value) (value, error) {
	if _, ok := other.(nilval); ok {
		return Number(0), nil
	}
	return nil, nil
}

func (n nilval) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil > %s not supported", reflect.TypeOf(other))
}

func (n nilval) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil >= %s not supported", reflect.TypeOf(other))
}

func (n nilval) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil < %s not supported", reflect.TypeOf(other))
}

func (n nilval) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil <= %s not supported", reflect.TypeOf(other))
}

func (n nilval) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil + %s not supported", reflect.TypeOf(other))
}

func (n nilval) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil - %s not supported", reflect.TypeOf(other))
}

func (n nilval) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil * %s not supported", reflect.TypeOf(other))
}

func (n nilval) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil / %s not supported", reflect.TypeOf(other))
}

func (n nilval) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil %% %s not supported", reflect.TypeOf(other))
}

func (n nilval) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil in %s not supported", reflect.TypeOf(other))
}

func (n nilval) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -nil not supported")
}

func (n nilval) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not nil not supported")
}

func (n nilval) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @nil not supported")
}

func (n nilval) property(ident string) (value, error) {
	return baseProperty(n, ident)
}

func (n nilval) printStr() string {
	return "nil"
}

func (n nilval) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over nil")
}

func (n nilval) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil[index] not supported")
}

func (n nilval) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil[lower..upper] not supported")
}

func (n nilval) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: nil[%s] not supported", reflect.TypeOf(index))
}

func (n nilval) runtimeTypeName() string {
	return "nil"
}

func (n nilval) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: nil :: [%s] not supported", reflect.TypeOf(val))
}
