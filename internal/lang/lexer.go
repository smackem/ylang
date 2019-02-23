package lang

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// tokenType represents the type of token
type tokenType int

// The types of tokens that can be lexed
const (
	ttNumber tokenType = iota
	ttIdent
	ttString
	ttColor
	ttAnd
	ttOr
	ttColonEq
	ttEq
	ttEqEq
	ttNeq
	ttGt
	ttGe
	ttLt
	ttLe
	ttPlus
	ttMinus
	ttStar
	ttSlash
	ttPercent
	ttLParen
	ttRParen
	ttNot
	ttAt
	ttComma
	ttDot
	ttSemicolon
	ttFor
	ttIn
	ttYield
	ttDotDot
	ttIf
	ttElse
	ttLBrace
	ttRBrace
	ttLBracket
	ttRBracket
	ttTrue
	ttFalse
	ttQMark
	ttColon
	ttLog
	ttBlt
	ttCommit
	ttFn
	ttReturn
	ttArrow
	ttNil
	ttPipe
	ttEOF
)

var emptyToken = token{Type: ttEOF, Lexeme: ""}
var tokenTypeNames = []string{
	"number",
	"ident",
	"string",
	"color",
	"and",
	"or",
	":=",
	"=",
	"==",
	"!=",
	">",
	">=",
	"<",
	"<=",
	"+",
	"-",
	"*",
	"/",
	"%",
	"(",
	")",
	"not",
	"@",
	",",
	".",
	";",
	"for",
	"in",
	"yield",
	"..",
	"if",
	"else",
	"{",
	"}",
	"[",
	"]",
	"true",
	"false",
	"?",
	":",
	"log",
	"blt",
	"commit",
	"fn",
	"return",
	"->",
	"nil",
	"|",
	"eof",
}

// token represents a lexed token
type token struct {
	Type       tokenType
	Lexeme     string
	LineNumber int
}

func (t token) String() string {
	return t.Lexeme
}

func (t token) parseNumber() Number {
	if t.Type == ttNumber {
		if n, err := strconv.ParseFloat(t.Lexeme, 64); err == nil {
			return Number(n)
		}
	}

	panic(fmt.Sprintf("error converting %s to number", t.Lexeme))
}

func (t token) parseColor() Color {
	if t.Type == ttColor {
		str := t.Lexeme[1:]
		if strings.Contains(str, ":") {
			str = strings.Replace(str, ":", "", -1)
		} else {
			str += "ff"
		}
		u, err := strconv.ParseUint(str, 16, 32)
		if err == nil {
			r := (u & 0xff000000) >> 24
			g := (u & 0x00ff0000) >> 16
			b := (u & 0x0000ff00) >> 8
			a := (u & 0x000000ff) >> 0
			return NewRgba(Number(r), Number(g), Number(b), Number(a))
		}
	}

	panic(fmt.Sprintf("error converting %s to color", t.Lexeme))
}

// lex walks the specified string and returns an array of lexed Tokens
// or an non-nil error if the input could not be lexed.
func lex(src string) ([]token, error) {
	tokens := []token{}
	lineNumber := 1
	isNotSpace := func(r rune) bool {
		if r == '\n' {
			lineNumber++
		}
		return unicode.IsSpace(r) == false
	}

	for index := 0; index < len(src); {
		slice := src[index:]
		spaceCount := strings.IndexFunc(slice, isNotSpace)
		if spaceCount < 0 {
			break // eof
		}
		index += spaceCount
		slice = src[index:]

		if strings.HasPrefix(slice, "//") {
			for src[index] != '\n' {
				index++
			}
			continue
		}

		if token, lexemeLen := match(slice); lexemeLen >= 0 {
			token.LineNumber = lineNumber
			tokens = append(tokens, token)
			index += lexemeLen
		} else {
			return tokens, fmt.Errorf("Error lexing '%s'", slice)
		}
	}

	return tokens, nil
}

func match(src string) (token, int) {
	for _, m := range matchers {
		if loc := m.regexp.FindStringIndex(src); loc != nil {
			lexeme := src[loc[0]:loc[1]]
			tok := token{
				Type:   m.lookup(lexeme),
				Lexeme: lexeme,
			}
			return tok, loc[1]
		}
	}
	return emptyToken, -1
}

type tokenLookup func(lexeme string) tokenType

type matcher struct {
	regexp *regexp.Regexp
	lookup tokenLookup
}

func makeMatcher(pattern string, lookup tokenLookup) matcher {
	return matcher{
		regexp: regexp.MustCompile("^" + pattern),
		lookup: lookup,
	}
}

var keywordTokens = map[string]tokenType{
	"and":    ttAnd,
	"or":     ttOr,
	"not":    ttNot,
	"for":    ttFor,
	"in":     ttIn,
	"yield":  ttYield,
	"if":     ttIf,
	"else":   ttElse,
	"true":   ttTrue,
	"false":  ttFalse,
	"log":    ttLog,
	"blt":    ttBlt,
	"commit": ttCommit,
	"fn":     ttFn,
	"return": ttReturn,
	"nil":    ttNil,
}

func lookupKeyword(lexeme string) tokenType {
	if token, ok := keywordTokens[lexeme]; ok {
		return token
	}
	return ttIdent
}

var matchers = []matcher{
	makeMatcher(`\(`, func(string) tokenType { return ttLParen }),
	makeMatcher(`\)`, func(string) tokenType { return ttRParen }),
	makeMatcher(`\+`, func(string) tokenType { return ttPlus }),
	makeMatcher(`\->`, func(string) tokenType { return ttArrow }),
	makeMatcher(`\-`, func(string) tokenType { return ttMinus }),
	makeMatcher(`\*`, func(string) tokenType { return ttStar }),
	makeMatcher(`/`, func(string) tokenType { return ttSlash }),
	makeMatcher(`%`, func(string) tokenType { return ttPercent }),
	makeMatcher(`,`, func(string) tokenType { return ttComma }),
	makeMatcher(`\:=`, func(string) tokenType { return ttColonEq }),
	makeMatcher(`\:`, func(string) tokenType { return ttColon }),
	makeMatcher(`\[`, func(string) tokenType { return ttLBracket }),
	makeMatcher(`\]`, func(string) tokenType { return ttRBracket }),
	makeMatcher(`==`, func(string) tokenType { return ttEqEq }),
	makeMatcher(`=`, func(string) tokenType { return ttEq }),
	makeMatcher(`!=`, func(string) tokenType { return ttNeq }),
	makeMatcher(`>=`, func(string) tokenType { return ttGe }),
	makeMatcher(`>`, func(string) tokenType { return ttGt }),
	makeMatcher(`<=`, func(string) tokenType { return ttLe }),
	makeMatcher(`<`, func(string) tokenType { return ttLt }),
	makeMatcher(`@`, func(string) tokenType { return ttAt }),
	makeMatcher(`\.\.`, func(string) tokenType { return ttDotDot }),
	makeMatcher(`\.`, func(string) tokenType { return ttDot }),
	makeMatcher(`;`, func(string) tokenType { return ttSemicolon }),
	makeMatcher(`{`, func(string) tokenType { return ttLBrace }),
	makeMatcher(`}`, func(string) tokenType { return ttRBrace }),
	makeMatcher(`\?`, func(string) tokenType { return ttQMark }),
	makeMatcher(`\|`, func(string) tokenType { return ttPipe }),
	makeMatcher(`#[0-9a-fA-F]{6}(\:[0-9a-fA-F]{2})?`, func(string) tokenType { return ttColor }),
	makeMatcher(`".*?"`, func(string) tokenType { return ttString }),
	makeMatcher(`[0-9]+(\.[0-9]+)?\b`, func(string) tokenType { return ttNumber }),
	makeMatcher(`[a-zA-Z_][a-zA-Z0-9_]*\b`, lookupKeyword),
}
