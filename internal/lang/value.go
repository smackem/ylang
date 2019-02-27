package lang

import "fmt"

type value interface {
	equals(other value) (value, error)
	greaterThan(other value) (value, error)
	greaterThanOrEqual(other value) (value, error)
	lessThan(other value) (value, error)
	lessThanOrEqual(other value) (value, error)
	add(other value) (value, error)
	sub(other value) (value, error)
	mul(other value) (value, error)
	div(other value) (value, error)
	mod(other value) (value, error)
	in(other value) (value, error)
	neg() (value, error)
	not() (value, error)
	at(bitmap BitmapContext) (value, error)
	property(ident string) (value, error)
	printStr() string
	iterate(visit func(value) error) error
	index(index value) (value, error)
	indexAssign(index value, val value) error
	runtimeTypeName() string
	concat(val value) (value, error)
}

var falseVal = boolean(false)

func baseProperty(val value, ident string) (value, error) {
	switch ident {
	case "__type":
		return str(val.runtimeTypeName()), nil
	}
	return nil, fmt.Errorf("unknown property '%s.%s'", val.runtimeTypeName(), ident)
}
