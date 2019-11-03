package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
	"sort"
)

type hashMap map[Value]Value

func (h hashMap) sortedKeys() []Value {
	keys := make([]Value, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	sort.Sort(valueSlice(keys))
	return keys
}

// implement value

func (h hashMap) compare(other Value) (Value, error) {
	return boolean(lang.FalseVal), nil
}

func (h hashMap) add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap + %s not supported", reflect.TypeOf(other))
}

func (h hashMap) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap - %s not supported", reflect.TypeOf(other))
}

func (h hashMap) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap * %s not supported", reflect.TypeOf(other))
}

func (h hashMap) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap / %s not supported", reflect.TypeOf(other))
}

func (h hashMap) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap %% %s not supported", reflect.TypeOf(other))
}

func (h hashMap) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap in %s not supported", reflect.TypeOf(other))
}

func (h hashMap) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -hashMap not supported")
}

func (h hashMap) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: not hashMap not supported")
}

func (h hashMap) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @hashMap not supported")
}

func (h hashMap) property(ident string) (Value, error) {
	switch ident {
	case "count":
		return number(len(h)), nil
	default:
		if val, ok := h[str(ident)]; ok {
			return val, nil
		}
	}
	return baseProperty(h, ident)
}

func (h hashMap) printStr() string {
	return fmt.Sprintf("hashMap(count: %d)", len(h))
}

func (h hashMap) iterate(visit func(Value) error) error {
	for key, val := range h {
		entry := hashMap{
			str("key"): key,
			str("val"): val,
		}
		if err := visit(entry); err != nil {
			return err
		}
	}
	return nil
}

func (h hashMap) index(index Value) (Value, error) {
	if val, ok := h[index]; ok {
		return val, nil
	}
	return nil, nil
}

func (h hashMap) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap[lower..upper] not supported")
}

func (h hashMap) indexAssign(index Value, val Value) error {
	h[index] = val
	return nil
}

func (h hashMap) runtimeTypeName() string {
	return "hashmap"
}

func (h hashMap) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap :: [%s] not supported", reflect.TypeOf(val))
}
