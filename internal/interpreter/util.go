package interpreter

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"strings"
)

var falseVal = boolean(false)

func baseProperty(val value, ident string) (value, error) {
	switch ident {
	case "__type":
		return str(val.runtimeTypeName()), nil
	}
	return nil, fmt.Errorf("unknown property '%s.%s'", val.runtimeTypeName(), ident)
}

func indexAt(n number, count int) int {
	if n < 0 {
		return count + int(n)
	}
	return int(n)
}

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
		keys := v.sortedKeys()
		for _, key := range keys {
			elem := v[key]
			buf.WriteString(fmt.Sprintf("%s%s: %s,\n", innerIndent, key, formatValue(elem, innerIndent, false)))
		}
		buf.WriteString(fmt.Sprintf("%s}", indent))
	case kernel:
		buf.WriteByte('|')
		innerIndent := indent + " "
		width := findMaxWidth(v.values) + 3
		for i, elem := range v.values {
			if v.width != 0 && i > 0 && i%v.width == 0 {
				buf.WriteByte('\n')
				buf.WriteString(innerIndent)
			}
			buf.WriteString(fmt.Sprintf(" %*.2f", width, elem))
		}
		buf.WriteString("|")
	default:
		buf.WriteString(v.printStr())
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

type valueSlice []value

func (s valueSlice) Len() int { return len(s) }
func (s valueSlice) Less(i int, j int) bool {
	cmp, _ := s[i].compare(s[j])
	if cmp == nil {
		return false
	}
	n, ok := cmp.(number)
	return ok && n < 0
}
func (s valueSlice) Swap(i int, j int) { s[i], s[j] = s[j], s[i] }
