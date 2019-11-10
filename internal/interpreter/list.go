package interpreter

import (
	"fmt"
	"reflect"
)

type List struct {
	Elements []Value
}

func (l List) Compare(other Value) (Value, error) {
	if r, ok := other.(List); ok {
		if reflect.DeepEqual(l, r) {
			return Number(0), nil
		}
	}
	return nil, nil
}

func (l List) Add(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list + %s Not supported", reflect.TypeOf(other))
}

func (l List) Sub(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list - %s Not supported", reflect.TypeOf(other))
}

func (l List) Mul(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list * %s Not supported", reflect.TypeOf(other))
}

func (l List) Div(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list / %s Not supported", reflect.TypeOf(other))
}

func (l List) Mod(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list %% %s Not supported", reflect.TypeOf(other))
}

func (l List) In(other Value) (Value, error) {
	return nil, fmt.Errorf("type mismatch: list In %s Not supported", reflect.TypeOf(other))
}

func (l List) Neg() (Value, error) {
	return nil, fmt.Errorf("type mismatch: -list Not supported")
}

func (l List) Not() (Value, error) {
	return nil, fmt.Errorf("type mismatch: 'Not list' Not supported")
}

func (l List) At(bitmap BitmapContext) (Value, error) {
	return nil, fmt.Errorf("type mismatch: @list Not supported")
}

func (l List) Property(ident string) (Value, error) {
	switch ident {
	case "count":
		return Number(len(l.Elements)), nil
	}
	return baseProperty(l, ident)
}

func (l List) PrintStr() string {
	return fmt.Sprintf("list(count: %d)", len(l.Elements))
}

func (l List) Iterate(visit func(Value) error) error {
	for _, v := range l.Elements {
		if err := visit(v); err != nil {
			return err
		}
	}
	return nil
}

func (l List) Index(index Value) (Value, error) {
	i, ok := index.(Number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	return l.Elements[indexAt(i, len(l.Elements))], nil
}

func (l List) IndexRange(lower, upper Value) (Value, error) {
	lowern, ok := lower.(Number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: kernel[%s..upper] Not supported", reflect.TypeOf(lower))
	}
	loweri := indexAt(lowern, len(l.Elements))
	uppern, ok := upper.(Number)
	if !ok {
		return nil, fmt.Errorf("type mismatch: kernel[lower..%s] Not supported", reflect.TypeOf(upper))
	}
	upperi := indexAt(uppern, len(l.Elements))
	return List{
		Elements: l.Elements[int(loweri) : upperi+1],
	}, nil
}

func (l List) IndexAssign(index Value, val Value) error {
	i, ok := index.(Number)
	if !ok {
		return fmt.Errorf("type mismatch: expected list[number] but found list[%s]", reflect.TypeOf(index))
	}
	l.Elements[indexAt(i, len(l.Elements))] = val
	return nil
}

func (l List) RuntimeTypeName() string {
	return "list"
}

func (l List) Concat(val Value) (Value, error) {
	if r, ok := val.(List); ok {
		return List{
			Elements: append(l.Elements, r.Elements...),
		}, nil
	}
	return List{
		Elements: append(l.Elements, val),
	}, nil
}
