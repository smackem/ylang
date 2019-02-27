package lang

import (
	"fmt"
	"reflect"
)

type list struct {
	elements []value
}

func (l list) equals(other value) (value, error) {
	if r, ok := other.(list); ok {
		return boolean(reflect.DeepEqual(l, r)), nil
	}
	return falseVal, nil
}

func (l list) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list > %s not supported", reflect.TypeOf(other))
}

func (l list) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list >= %s not supported", reflect.TypeOf(other))
}

func (l list) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list < %s not supported", reflect.TypeOf(other))
}

func (l list) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list <= %s not supported", reflect.TypeOf(other))
}

func (l list) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list + %s not supported", reflect.TypeOf(other))
}

func (l list) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list - %s not supported", reflect.TypeOf(other))
}

func (l list) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list * %s not supported", reflect.TypeOf(other))
}

func (l list) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list / %s not supported", reflect.TypeOf(other))
}

func (l list) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list %% %s not supported", reflect.TypeOf(other))
}

func (l list) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: list in %s not supported", reflect.TypeOf(other))
}

func (l list) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -list not supported")
}

func (l list) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not list' not supported")
}

func (l list) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @list not supported")
}

func (l list) property(ident string) (value, error) {
	switch ident {
	case "count":
		return Number(len(l.elements)), nil
	}
	return baseProperty(l, ident)
}

func (l list) printStr() string {
	return fmt.Sprintf("list(count: %d)", len(l.elements))
}

func (l list) iterate(visit func(value) error) error {
	for _, v := range l.elements {
		if err := visit(v); err != nil {
			return err
		}
	}
	return nil
}

func (l list) index(index value) (value, error) {
	i, ok := index.(Number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	return l.elements[int(i)], nil
}

func (l list) indexAssign(index value, val value) error {
	i, ok := index.(Number)
	if !ok {
		return fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	l.elements[int(i)] = val
	return nil
}

func (l list) runtimeTypeName() string {
	return "list"
}

func (l list) concat(val value) (value, error) {
	if r, ok := val.(list); ok {
		return list{
			elements: append(l.elements, r.elements...),
		}, nil
	}
	return list{
		elements: append(l.elements, val),
	}, nil
}
