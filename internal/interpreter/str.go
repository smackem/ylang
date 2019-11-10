package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
	"strings"
	"unicode/utf8"
)

type Str lang.Str

func (s Str) Compare(other Value) (Value, error) {
	if r, ok := other.(Str); ok {
		return Number(strings.Compare(string(s), string(r))), nil
	}
	return nil, nil
}

func (s Str) Add(other Value) (Value, error) {
	return Str(fmt.Sprintf("%s%s", s, other.PrintStr())), nil
}

func (s Str) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string - %s Not supported", reflect.TypeOf(other))
}

func (s Str) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string * %s Not supported", reflect.TypeOf(other))
}

func (s Str) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string / %s Not supported", reflect.TypeOf(other))
}

func (s Str) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string %% %s Not supported", reflect.TypeOf(other))
}

func (s Str) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string In %s Not supported", reflect.TypeOf(other))
}

func (s Str) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -string Not supported")
}

func (s Str) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: Not string Not supported")
}

func (s Str) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @string Not supported")
}

func (s Str) Property(ident string) (Value, error) {
	switch ident {
	case "len", "length":
		return Number(len(s)), nil
	}
	return baseProperty(s, ident)
}

func (s Str) PrintStr() string {
	return string(s)
}

func (s Str) Iterate(visit func(Value) error) error {
	return fmt.Errorf("cannot Iterate over string")
}

func (s Str) Index(index Value) (Value, error) {
	if n, ok := index.(Number); ok {
		i := int(n)
		runeVal, _ := utf8.DecodeRuneInString(string(s[i:]))
		return Str(fmt.Sprint(runeVal)), nil
	}
	return nil, fmt.Errorf("type mismatch: str[%s] Not supported", reflect.TypeOf(index))
}

func (s Str) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string[lower..upper] Not supported")
}

func (s Str) IndexAssign(index Value, val Value) error {
	return fmt.Errorf("type mismatch: str[%s] Not supported", reflect.TypeOf(index))
}

func (s Str) RuntimeTypeName() string {
	return "string"
}

func (s Str) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: string :: [%s] Not supported", reflect.TypeOf(val))
}
