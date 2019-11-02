package interpreter

import (
	"fmt"
	"reflect"
)

type list struct {
	elements []value
}

func (l list) compare(other value) (value, error) {
	if r, ok := other.(list); ok {
		if reflect.DeepEqual(l, r) {
			return number(0), nil
		}
	}
	return nil, nil
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
		return number(len(l.elements)), nil
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
	i, ok := index.(number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	return l.elements[indexAt(i, len(l.elements))], nil
}

func (l list) indexRange(lower, upper value) (value, error) {
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

func (l list) indexAssign(index value, val value) error {
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
