package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/parser"
	"reflect"
	"strings"
)

type Function struct {
	ParameterNames []string
	Body           []parser.Statement
	closure        []scope
}

func (f Function) Compare(other Value) (Value, error) {
	return nil, nil
}

func (f Function) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function + %s Not supported", reflect.TypeOf(other))
}

func (f Function) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function - %s Not supported", reflect.TypeOf(other))
}

func (f Function) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function * %s Not supported", reflect.TypeOf(other))
}

func (f Function) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function / %s Not supported", reflect.TypeOf(other))
}

func (f Function) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function %% %s Not supported", reflect.TypeOf(other))
}

func (f Function) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function In %s Not supported", reflect.TypeOf(other))
}

func (f Function) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -function Not supported")
}

func (f Function) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: Not function Not supported")
}

func (f Function) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @function Not supported")
}

func (f Function) Property(ident string) (Value, error) {
	return baseProperty(f, ident)
}

func (f Function) PrintStr() string {
	return fmt.Sprintf("fn(%s)", strings.Join(f.ParameterNames, ", "))
}

func (f Function) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over function")
}

func (f Function) Index(index Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function[Index] Not supported")
}

func (f Function) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function[lower..upper] Not supported")
}

func (f Function) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: function[%s] Not supported", reflect.TypeOf(index))
}

func (f Function) RuntimeTypeName() string {
	return "function"
}

func (f Function) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: function :: [%s] Not supported", reflect.TypeOf(val))
}
