package interpreter

type Value interface {
	compare(other Value) (Value, error)
	add(other Value) (Value, error)
	sub(other Value) (Value, error)
	mul(other Value) (Value, error)
	div(other Value) (Value, error)
	mod(other Value) (Value, error)
	in(other Value) (Value, error)
	neg() (Value, error)
	not() (Value, error)
	at(bitmap BitmapContext) (Value, error)
	property(ident string) (Value, error)
	printStr() string
	iterate(visit func(Value) error) error
	index(index Value) (Value, error)
	indexRange(lower, upper Value) (Value, error)
	indexAssign(index Value, val Value) error
	runtimeTypeName() string
	concat(val Value) (Value, error)
}
