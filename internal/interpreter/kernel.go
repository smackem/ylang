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

func (k kernel) compare(other Value) (Value, error) {
	if r, ok := other.(kernel); ok {
		if reflect.DeepEqual(k, r) {
			return number(0), nil
		}
	}
	return nil, nil
}

func (k kernel) add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel + %s not supported", reflect.TypeOf(other))
}

func (k kernel) sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel - %s not supported", reflect.TypeOf(other))
}

func (k kernel) mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel * %s not supported", reflect.TypeOf(other))
}

func (k kernel) div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel / %s not supported", reflect.TypeOf(other))
}

func (k kernel) mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel %% %s not supported", reflect.TypeOf(other))
}

func (k kernel) in(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel in %s not supported", reflect.TypeOf(other))
}

func (k kernel) neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -kernel not supported")
}

func (k kernel) not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'not kernel' not supported")
}

func (k kernel) at(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @kernel not supported")
}

func (k kernel) property(ident string) (Value, error) {
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

func (k kernel) iterate(visit func(Value) error) error {
	for _, v := range k.values {
		if err := visit(number(v)); err != nil {
			return err
		}
	}
	return nil
}

func (k kernel) index(index Value) (Value, error) {
	switch i := index.(type) {
	case number:
		return number(k.values[indexAt(i, len(k.values))]), nil
	case point:
		return number(k.values[i.Y*k.width+i.X]), nil
	}
	return nil, fmt.Errorf("type mismatch: expected kernel[number] or kernel[point] but found kernel[%s]", reflect.TypeOf(index))
}

func (k kernel) indexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel[lower..upper] not supported")
}

func (k kernel) indexAssign(index Value, val Value) error {
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

func (k kernel) concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel :: [%s] not supported", reflect.TypeOf(val))
}
