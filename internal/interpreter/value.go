package interpreter

// Value represents a ylang value at runtime
type Value interface {
	Compare(other Value) (Value, error)
	Add(other Value) (Value, error)
	Sub(other Value) (Value, error)
	Mul(other Value) (Value, error)
	Div(other Value) (Value, error)
	Mod(other Value) (Value, error)
	In(other Value) (Value, error)
	Neg() (Value, error)
	Not() (Value, error)
	At(bitmap BitmapContext) (Value, error)
	Property(ident string) (Value, error)
	PrintStr() string
	Iterate(visit func(Value) error) error
	Index(index Value) (Value, error)
	IndexRange(lower, upper Value) (Value, error)
	IndexAssign(index Value, val Value) error
	RuntimeTypeName() string
	Concat(val Value) (Value, error)
}
