package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"reflect"
)

type Kernel struct {
	Width  int
	Height int
	Values []lang.Number
}

func (k Kernel) Compare(other Value) (Value, error) {
	if r, ok := other.(Kernel); ok {
		if reflect.DeepEqual(k, r) {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (k Kernel) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel + %s Not supported", reflect.TypeOf(other))
}

func (k Kernel) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel - %s Not supported", reflect.TypeOf(other))
}

func (k Kernel) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel * %s Not supported", reflect.TypeOf(other))
}

func (k Kernel) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel / %s Not supported", reflect.TypeOf(other))
}

func (k Kernel) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel %% %s Not supported", reflect.TypeOf(other))
}

func (k Kernel) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel In %s Not supported", reflect.TypeOf(other))
}

func (k Kernel) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -kernel Not supported")
}

func (k Kernel) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not kernel' Not supported")
}

func (k Kernel) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @kernel Not supported")
}

func (k Kernel) Property(ident string) (Value, error) {
	switch ident {
	case "w", "width":
		return Number(k.Width), nil
	case "h", "height":
		return Number(k.Height), nil
	case "count":
		return Number(len(k.Values)), nil
	}
	return baseProperty(k, ident)
}

func (k Kernel) PrintStr() string {
	return fmt.Sprintf("kernel(width: %d, height: %d)", k.Width, k.Height)
}

func (k Kernel) Iterate(visit func(Value) error) error {
	for _, v := range k.Values {
		if err := visit(Number(v)); err != nil {
			return err
		}
	}
	return nil
}

func (k Kernel) Index(index Value) (Value, error) {
	switch i := index.(type) {
	case Number:
		return Number(k.Values[indexAt(i, len(k.Values))]), nil
	case Point:
		return Number(k.Values[i.Y*k.Width+i.X]), nil
	}
	return nil, fmt.Errorf("type mismatch: expected kernel[number] or kernel[point] but found kernel[%s]", reflect.TypeOf(index))
}

func (k Kernel) IndexRange(lower, upper Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel[lower..upper] Not supported")
}

func (k Kernel) IndexAssign(index Value, val Value) error {
	nval, ok := val.(Number)
	if !ok {
		return fmt.Errorf("type mismatch: expected kernel[Index] = number but found kernel[Index] = %s", reflect.TypeOf(val))
	}
	switch i := index.(type) {
	case Number:
		k.Values[indexAt(i, len(k.Values))] = lang.Number(nval)
		return nil
	case Point:
		k.Values[i.Y*k.Width+i.X] = lang.Number(nval)
		return nil
	}
	return fmt.Errorf("type mismatch: expected kernel[number] or kernel[point] but found kernel[%s]", reflect.TypeOf(index))
}

func (k Kernel) RuntimeTypeName() string {
	return "kernel"
}

func (k Kernel) Concat(val Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: kernel :: [%s] Not supported", reflect.TypeOf(val))
}
