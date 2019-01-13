package lang

type value interface {
	isTrue() bool
	equals(other value) bool
	greaterThan(other value) bool
	greaterThanOrEqual(other value) bool
	lessThan(other value) bool
	lessThanOrEqual(other value) bool
	add(other value) value
	sub(other value) value
	mul(other value) value
	div(other value) value
	mod(other value) value
	in(other value) value
	neg() value
	not() value
	at(bitmap Bitmap) value
	property(ident string) value
}
