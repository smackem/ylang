package interpreter

import (
	"fmt"
	"reflect"
)

type list struct {
	elements []Value
}

func (l list) compare(other Value) (Value, error) {
	if r, ok := other.(list); ok {
		if reflect.DeepEqual(l, r) {
			return number(0), nil
		}
	}
	return nil, nil
}

func (l list) add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list + %s not supported", reflect.TypeOf(other))
}

func (l list) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list - %s not supported", reflect.TypeOf(other))
}

func (l list) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list * %s not supported", reflect.TypeOf(other))
}

func (l list) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list / %s not supported", reflect.TypeOf(other))
}

func (l list) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list %% %s not supported", reflect.TypeOf(other))
}

func (l list) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list in %s not supported", reflect.TypeOf(other))
}

func (l list) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -list not supported")
}

func (l list) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'not list' not supported")
}

func (l list) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @list not supported")
}

func (l list) property(ident string) (Value, error) {
	switch ident {
	case "count":
		return number(len(l.elements)), nil
	}
	return baseProperty(l, ident)
}

func (l list) printStr() string {
	return fmt.Sprintf("list(count: %d)", len(l.elements))
}

func (l list) iterate(visit func(Value) error) error {
	for _, v := range l.elements {
		if err := visit(v); err != nil {
			return err
		}
	}
	return nil
}

func (l list) index(index Value) (Value, error) {
	i, ok := index.(number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	return l.elements[indexAt(i, len(l.elements))], nil
}

func (l list) indexRange(lower, upper Value) (Value, error) {
	lowern, ok := lower.(number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: kernel[%s..upper] not supported", reflect.TypeOf(lower))
	}
	loweri := indexAt(lowern, len(l.elements))
	uppern, ok := upper.(number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: kernel[lower..%s] not supported", reflect.TypeOf(upper))
	}
	upperi := indexAt(uppern, len(l.elements))
	return list{
		elements: l.elements[int(loweri) : upperi+1],
	}, nil
}

func (l list) indexAssign(index Value, val Value) error {
	i, ok := index.(number)
	if !ok {
		return fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	l.elements[indexAt(i, len(l.elements))] = val
	return nil
}

func (l list) runtimeTypeName() string {
	return "list"
}

func (l list) concat(val Value) (Value, error) {
	if r, ok := val.(list); ok {
		return list{
			elements: append(l.elements, r.elements...),
		}, nil
	}
	return list{
		elements: append(l.elements, val),
	}, nil
}
