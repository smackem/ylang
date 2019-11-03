package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type nilval lang.Nil

func (n nilval) compare(other Value) (Value, error) {
	if _, ok := other.(nilval); ok {
		return number(0), nil
	}
	return nil, nil
}

func (n nilval) greaterThan(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil > %s not supported", reflect.TypeOf(other))
}

func (n nilval) greaterThanOrEqual(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil >= %s not supported", reflect.TypeOf(other))
}

func (n nilval) lessThan(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil < %s not supported", reflect.TypeOf(other))
}

func (n nilval) lessThanOrEqual(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil <= %s not supported", reflect.TypeOf(other))
}

func (n nilval) add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil + %s not supported", reflect.TypeOf(other))
}

func (n nilval) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil - %s not supported", reflect.TypeOf(other))
}

func (n nilval) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil * %s not supported", reflect.TypeOf(other))
}

func (n nilval) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil / %s not supported", reflect.TypeOf(other))
}

func (n nilval) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil %% %s not supported", reflect.TypeOf(other))
}

func (n nilval) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil in %s not supported", reflect.TypeOf(other))
}

func (n nilval) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -nil not supported")
}

func (n nilval) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: not nil not supported")
}

func (n nilval) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @nil not supported")
}

func (n nilval) property(ident string) (Value, error) {
	return baseProperty(n, ident)
}

func (n nilval) printStr() string {
	return "nil"
}

func (n nilval) iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot iterate over nil")
}

func (n nilval) index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil[index] not supported")
}

func (n nilval) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil[lower..upper] not supported")
}

func (n nilval) indexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: nil[%s] not supported", reflect.TypeOf(index))
}

func (n nilval) runtimeTypeName() string {
	return "nil"
}

func (n nilval) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: nil :: [%s] not supported", reflect.TypeOf(val))
}
