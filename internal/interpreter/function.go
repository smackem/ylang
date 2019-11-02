package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/parser"
	"reflect"
	"strings"
)

type function struct {
	parameterNames []string
	body           []parser.Statement
	closure        []scope
}

func (f function) compare(other value) (value, error) {
	return nil, nil
}

func (f function) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function + %s not supported", reflect.TypeOf(other))
}

func (f function) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function - %s not supported", reflect.TypeOf(other))
}

func (f function) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function * %s not supported", reflect.TypeOf(other))
}

func (f function) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function / %s not supported", reflect.TypeOf(other))
}

func (f function) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function %% %s not supported", reflect.TypeOf(other))
}

func (f function) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function in %s not supported", reflect.TypeOf(other))
}

func (f function) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -function not supported")
}

func (f function) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not function not supported")
}

func (f function) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @function not supported")
}

func (f function) property(ident string) (value, error) {
	return baseProperty(f, ident)
}

func (f function) printStr() string {
	return fmt.Sprintf("fn(%s)", strings.Join(f.parameterNames, ", "))
}

func (f function) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over function")
}

func (f function) index(index value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function[index] not supported")
}

func (f function) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function[lower..upper] not supported")
}

func (f function) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: function[%s] not supported", reflect.TypeOf(index))
}

func (f function) runtimeTypeName() string {
	return "function"
}

func (f function) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: function :: [%s] not supported", reflect.TypeOf(val))
}
