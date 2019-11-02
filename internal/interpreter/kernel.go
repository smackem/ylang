package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type kernel struct {
	width  int
	height int
	values []lang.Number
}

func (k kernel) compare(other value) (value, error) {
	if r, ok := other.(kernel); ok {
		if reflect.DeepEqual(k, r) {
			return number(0), nil
		}
	}
	return nil, nil
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
	case "w", "width":
		return number(k.width), nil
	case "h", "height":
		return number(k.height), nil
	case "count":
		return number(len(k.values)), nil
	}
	return baseProperty(k, ident)
}

func (k kernel) printStr() string {
	return fmt.Sprintf("kernel(width: %d, height: %d)", k.width, k.height)
}

func (k kernel) iterate(visit func(value) error) error {
	for _, v := range k.values {
		if err := visit(number(v)); err != nil {
			return err
		}
	}
	return nil
}

func (k kernel) index(index value) (value, error) {
	switch i := index.(type) {
	case number:
		return number(k.values[indexAt(i, len(k.values))]), nil
	case point:
		return number(k.values[i.Y*k.width+i.X]), nil
	}
	return nil, fmt.Errorf("type mismatch: expected kernel[number] or kernel[point] but found kernel[%s]", reflect.TypeOf(index))
}

func (k kernel) indexRange(lower, upper value) (value, error) {
	return nil, fmt.Errorf("type mismatch: kernel[lower..upper] not supported")
}

func (k kernel) indexAssign(index value, val value) error {
	nval, ok := val.(number)
	if !ok {
		return fmt.Errorf("type mismatch: expected kernel[index] = number but found kernel[index] = %s", reflect.TypeOf(val))
	}
	switch i := index.(type) {
	case number:
		k.values[indexAt(i, len(k.values))] = lang.Number(nval)
		return nil
	case point:
		k.values[i.Y*k.width+i.X] = lang.Number(nval)
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
