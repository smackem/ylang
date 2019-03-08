package lang

import (
	"fmt"
	"strings"
)

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
	indexRange(lower, upper value) (value, error)
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

func indexAt(n Number, count int) int {
	if n < 0 {
		return count + int(n)
	}
	return int(n)
}

/*
{
	x: [
		1,
		2,
		3,
		4,
		[
			33,
			44,
			55,
		]
	],
	y: abcdef,
}
*/

func formatValue(val value, indent string, leadingIndent bool) string {
	buf := strings.Builder{}
	if leadingIndent {
		buf.WriteString(indent)
	}
	switch v := val.(type) {
	case list:
		buf.WriteString("[\n")
		innerIndent := indent + "  "
		for _, elem := range v.elements {
			buf.WriteString(fmt.Sprintf("%s,\n", formatValue(elem, innerIndent, true)))
		}
		buf.WriteString(fmt.Sprintf("%s]", indent))
	case hashMap:
		buf.WriteString("{\n")
		innerIndent := indent + "  "
		for key, elem := range v {
			buf.WriteString(fmt.Sprintf("%s%s: %s,\n", innerIndent, key, formatValue(elem, innerIndent, false)))
		}
		buf.WriteString(fmt.Sprintf("%s}", indent))
	case kernel:
		buf.WriteByte('|')
		innerIndent := indent + " "
		for i, elem := range v.values {
			if v.width != 0 && i > 0 && i%v.width == 0 {
				buf.WriteByte('\n')
				buf.WriteString(innerIndent)
			}
			buf.WriteString(fmt.Sprintf(" %5.2f", elem))
		}
		buf.WriteString("|")
	default:
		buf.WriteString(v.printStr())
	}
	return buf.String()
}
