package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
	"strings"
	"unicode/utf8"
)

type str lang.Str

func (s str) compare(other Value) (Value, error) {
	if r, ok := other.(str); ok {
		return number(strings.Compare(string(s), string(r))), nil
	}
	return nil, nil
}

func (s str) add(other Value) (Value, error) {
	return str(fmt.Sprintf("%s%s", s, other.printStr())), nil
}

func (s str) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string - %s not supported", reflect.TypeOf(other))
}

func (s str) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string * %s not supported", reflect.TypeOf(other))
}

func (s str) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string / %s not supported", reflect.TypeOf(other))
}

func (s str) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string %% %s not supported", reflect.TypeOf(other))
}

func (s str) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string in %s not supported", reflect.TypeOf(other))
}

func (s str) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -string not supported")
}

func (s str) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: not string not supported")
}

func (s str) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @string not supported")
}

func (s str) property(ident string) (Value, error) {
	switch ident {
	case "len", "length":
		return number(len(s)), nil
	}
	return baseProperty(s, ident)
}

func (s str) printStr() string {
	return string(s)
}

func (s str) iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot iterate over string")
}

func (s str) index(index Value) (Value, error) {
	if n, ok := index.(number); ok {
		i := int(n)
		runeVal, _ := utf8.DecodeRuneInString(string(s[i:]))
		return str(fmt.Sprint(runeVal)), nil
	}
	return nil, fmt.Errorf("type mismatch: str[%s] not supported", reflect.TypeOf(index))
}

func (s str) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string[lower..upper] not supported")
}

func (s str) indexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: str[%s] not supported", reflect.TypeOf(index))
}

func (s str) runtimeTypeName() string {
	return "string"
}

func (s str) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string :: [%s] not supported", reflect.TypeOf(val))
}
