package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"strings"
)

var falseVal = Boolean(false)

func baseProperty(val Value, ident string) (Value, error) {
	switch ident {
	case "__type":
		return Str(val.RuntimeTypeName()), nil
	}
	return nil, fmt.Errorf("unknown Property '%s.%s'", val.RuntimeTypeName(), ident)
}

func indexAt(n Number, count int) int {
	if n < 0 {
		return count + int(n)
	}
	return int(n)
}

func formatValue(val Value, indent string, leadingIndent bool) string {
	buf := strings.Builder{}
	if leadingIndent {
		buf.WriteString(indent)
	}
	switch v := val.(type) {
	case List:
		buf.WriteString("[\n")
		innerIndent := indent + "  "
		for _, elem := range v.Elements {
			buf.WriteString(fmt.Sprintf("%s,\n", formatValue(elem, innerIndent, true)))
		}
		buf.WriteString(fmt.Sprintf("%s]", indent))
	case HashMap:
		buf.WriteString("{\n")
		innerIndent := indent + "  "
		keys := v.sortedKeys()
		for _, key := range keys {
			elem := v[key]
			buf.WriteString(fmt.Sprintf("%s%s: %s,\n", innerIndent, key, formatValue(elem, innerIndent, false)))
		}
		buf.WriteString(fmt.Sprintf("%s}", indent))
	case Kernel:
		buf.WriteByte('|')
		innerIndent := indent + " "
		width := findMaxWidth(v.Values) + 3
		for i, elem := range v.Values {
			if v.Width != 0 && i > 0 && i%v.Width == 0 {
				buf.WriteByte('\n')
				buf.WriteString(innerIndent)
			}
			buf.WriteString(fmt.Sprintf(" %*.2f", width, elem))
		}
		buf.WriteString("|")
	default:
		buf.WriteString(v.PrintStr())
	}
	return buf.String()
}

func findMaxWidth(values []lang.Number) int {
	max := 1
	for _, v := range values {
		count := 0
		for v > 1 {
			v /= 10
			count++
		}
		if count > max {
			max = count
		}
	}
	return max
}

// implement sort.Interface for slice of values

type valueSlice []Value

func (s valueSlice) Len() int { return len(s) }
func (s valueSlice) Less(i int, j int) bool {
	cmp, _ := s[i].Compare(s[j])
	if cmp == nil {
		return false
	}
	n, ok := cmp.(Number)
	return ok && n < 0
}
func (s valueSlice) Swap(i int, j int) { s[i], s[j] = s[j], s[i] }
