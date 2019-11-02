package interpreter

type value interface {
	compare(other value) (value, error)
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
	indexRange(lower, upper value) (value, error)
	indexAssign(index value, val value) error
	runtimeTypeName() string
	concat(val value) (value, error)
}
