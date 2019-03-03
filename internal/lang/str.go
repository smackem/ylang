package lang

import (
	"fmt"
	"reflect"
	"unicode/utf8"
)

type str string

func (s str) equals(other value) (value, error) {
	if r, ok := other.(str); ok {
		return boolean(s == r), nil
	}
	return falseVal, nil
}

func (s str) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string > %s not supported", reflect.TypeOf(other))
}

func (s str) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string >= %s not supported", reflect.TypeOf(other))
}

func (s str) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string < %s not supported", reflect.TypeOf(other))
}

func (s str) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string <= %s not supported", reflect.TypeOf(other))
}

func (s str) add(other value) (value, error) {
	return str(fmt.Sprintf("%s%s", s, other)), nil
}

func (s str) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string - %s not supported", reflect.TypeOf(other))
}

func (s str) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string * %s not supported", reflect.TypeOf(other))
}

func (s str) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string / %s not supported", reflect.TypeOf(other))
}

func (s str) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string %% %s not supported", reflect.TypeOf(other))
}

func (s str) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string in %s not supported", reflect.TypeOf(other))
}

func (s str) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -string not supported")
}

func (s str) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not string not supported")
}

func (s str) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @string not supported")
}

func (s str) property(ident string) (value, error) {
	switch ident {
	case "len", "length":
		return Number(len(s)), nil
	}
	return baseProperty(s, ident)
}

func (s str) printStr() string {
	return string(s)
}

func (s str) iterate(visit func(value) error) error {
	return fmt.Errorf("cannot iterate over string")
}

func (s str) index(index value) (value, error) {
	if n, ok := index.(Number); ok {
		i := int(n)
		runeVal, _ := utf8.DecodeRuneInString(string(s[i:]))
		return str(fmt.Sprint(runeVal)), nil
	}
	return nil, fmt.Errorf("type mismatch: str[%s] not supported", reflect.TypeOf(index))
}

func (s str) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string[lower..upper] not supported")
}

func (s str) indexAssign(index value, val value) error {
	return fmt.Errorf("type mismatch: str[%s] not supported", reflect.TypeOf(index))
}

func (s str) runtimeTypeName() string {
	return "string"
}

func (s str) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: string :: [%s] not supported", reflect.TypeOf(val))
}
