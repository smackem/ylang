package lang

import (
	"fmt"
	"reflect"
	"sort"
)

type hashMap map[value]value

func (h hashMap) sortedKeys() []value {
	keys := make([]value, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	sort.Sort(valueSlice(keys))
	return keys
}

// implement value

func (h hashMap) compare(other value) (value, error) {
	return falseVal, nil
}

func (h hashMap) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap + %s not supported", reflect.TypeOf(other))
}

func (h hashMap) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap - %s not supported", reflect.TypeOf(other))
}

func (h hashMap) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap * %s not supported", reflect.TypeOf(other))
}

func (h hashMap) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap / %s not supported", reflect.TypeOf(other))
}

func (h hashMap) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap %% %s not supported", reflect.TypeOf(other))
}

func (h hashMap) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap in %s not supported", reflect.TypeOf(other))
}

func (h hashMap) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -hashMap not supported")
}

func (h hashMap) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: not hashMap not supported")
}

func (h hashMap) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @hashMap not supported")
}

func (h hashMap) property(ident string) (value, error) {
	switch ident {
	case "count":
		return Number(len(h)), nil
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

func (h hashMap) iterate(visit func(value) error) error {
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

func (h hashMap) index(index value) (value, error) {
	if val, ok := h[index]; ok {
		return val, nil
	}
	return nil, nil
}

func (h hashMap) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap[lower..upper] not supported")
}

func (h hashMap) indexAssign(index value, val value) error {
	h[index] = val
	return nil
}

func (h hashMap) runtimeTypeName() string {
	return "hashmap"
}

func (h hashMap) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap :: [%s] not supported", reflect.TypeOf(val))
}
