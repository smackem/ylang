package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
	"sort"
)

type HashMap map[Value]Value

func (h HashMap) sortedKeys() []Value {
	keys := make([]Value, 0, len(h))
	for key := range h {
		keys = append(keys, key)
	}
	sort.Sort(valueSlice(keys))
	return keys
}

// implement value

func (h HashMap) Compare(other Value) (Value, error) {
	return Boolean(lang.FalseVal), nil
}

func (h HashMap) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap + %s Not supported", reflect.TypeOf(other))
}

func (h HashMap) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap - %s Not supported", reflect.TypeOf(other))
}

func (h HashMap) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap * %s Not supported", reflect.TypeOf(other))
}

func (h HashMap) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap / %s Not supported", reflect.TypeOf(other))
}

func (h HashMap) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap %% %s Not supported", reflect.TypeOf(other))
}

func (h HashMap) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap In %s Not supported", reflect.TypeOf(other))
}

func (h HashMap) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -hashMap Not supported")
}

func (h HashMap) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: Not hashMap Not supported")
}

func (h HashMap) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @hashMap Not supported")
}

func (h HashMap) Property(ident string) (Value, error) {
	switch ident {
	case "count":
		return Number(len(h)), nil
	default:
		if val, ok := h[Str(ident)]; ok {
			return val, nil
		}
	}
	return baseProperty(h, ident)
}

func (h HashMap) PrintStr() string {
	return fmt.Sprintf("hashMap(count: %d)", len(h))
}

func (h HashMap) Iterate(visit func(Value) error) error {
	for key, val := range h {
		entry := HashMap{
			Str("key"): key,
			Str("val"): val,
		}
		if err := visit(entry); err != nil {
			return err
		}
	}
	return nil
}

func (h HashMap) Index(index Value) (Value, error) {
	if val, ok := h[index]; ok {
		return val, nil
	}
	return nil, nil
}

func (h HashMap) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap[lower..upper] Not supported")
}

func (h HashMap) IndexAssign(index Value, val Value) error {
	h[index] = val
	return nil
}

func (h HashMap) RuntimeTypeName() string {
	return "hashmap"
}

func (h HashMap) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: hashMap :: [%s] Not supported", reflect.TypeOf(val))
}
