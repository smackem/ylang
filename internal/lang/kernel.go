package lang

import (
	"fmt"
	"reflect"
)

type kernel struct {
	width  int
	height int
	values []Number
}

func (k kernel) equals(other value) (value, error) {
	if r, ok := other.(kernel); ok {
		return boolean(reflect.DeepEqual(k, r)), nil
	}
	return falseVal, nil
}

func (k kernel) greaterThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel > %s not supported", reflect.TypeOf(other))
}

func (k kernel) greaterThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel >= %s not supported", reflect.TypeOf(other))
}

func (k kernel) lessThan(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel < %s not supported", reflect.TypeOf(other))
}

func (k kernel) lessThanOrEqual(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel <= %s not supported", reflect.TypeOf(other))
}

func (k kernel) add(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel + %s not supported", reflect.TypeOf(other))
}

func (k kernel) sub(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel - %s not supported", reflect.TypeOf(other))
}

func (k kernel) mul(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel * %s not supported", reflect.TypeOf(other))
}

func (k kernel) div(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel / %s not supported", reflect.TypeOf(other))
}

func (k kernel) mod(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel %% %s not supported", reflect.TypeOf(other))
}

func (k kernel) in(other value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel in %s not supported", reflect.TypeOf(other))
}

func (k kernel) neg() (value, error) {
	return nil, fmt.Errorf("type mismatch: -kernel not supported")
}

func (k kernel) not() (value, error) {
	return nil, fmt.Errorf("type mismatch: 'not kernel' not supported")
}

func (k kernel) at(bitmap BitmapContext) (value, error) {
	return nil, fmt.Errorf("type mismatch: @kernel not supported")
}

func (k kernel) property(ident string) (value, error) {
	switch ident {
	case "width":
		return Number(k.width), nil
	case "height":
		return Number(k.height), nil
	case "count":
		return Number(len(k.values)), nil
	}
	return baseProperty(k, ident)
}

func (k kernel) printStr() string {
	return fmt.Sprintf("kernel(width: %d, height: %d)", k.width, k.height)
}

func (k kernel) iterate(visit func(value) error) error {
	for _, v := range k.values {
		if err := visit(v); err != nil {
			return err
		}
	}
	return nil
}

func (k kernel) index(index value) (value, error) {
	switch i := index.(type) {
	case Number:
		return k.values[int(i)], nil
	case point:
		return k.values[i.Y*k.width+i.X], nil
	}
	return nil, fmt.Errorf("type mismatch: expected kernel[number] or kernel[point] but found kernel[%s]", reflect.TypeOf(index))
}

func (k kernel) indexAssign(index value, val value) error {
	nval, ok := val.(Number)
	if !ok {
		return fmt.Errorf("type mismatch: expected kernel[index] = number but found kernel[index] = %s", reflect.TypeOf(val))
	}
	switch i := index.(type) {
	case Number:
		k.values[int(i)] = nval
		return nil
	case point:
		k.values[i.Y*k.width+i.X] = nval
		return nil
	}
	return fmt.Errorf("type mismatch: expected kernel[number] or kernel[point] but found kernel[%s]", reflect.TypeOf(index))
}

func (k kernel) runtimeTypeName() string {
	return "kernel"
}

func (k kernel) concat(val value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel :: [%s] not supported", reflect.TypeOf(val))
}
