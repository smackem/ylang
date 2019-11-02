package lexer

import (
	"fmt"
	"github.com/smackem/ylang/internal/lang"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// TokenType represents the type of Token
type TokenType int

// The types of tokens that can be lexed
const (
	TTNumber TokenType = iota
	TTIdent
	TTString
	TTColor
	TTAnd
	TTOr
	TTColonEq
	TTEq
	TTEqEq
	TTNeq
	TTGt
	TTGe
	TTLt
	TTLe
	TTPlus
	TTMinus
	TTStar
	TTSlash
	TTPercent
	TTLParen
	TTRParen
	TTNot
	TTAt
	TTComma
	TTDot
	TTSemicolon
	TTFor
	TTIn
	TTYield
	TTDotDot
	TTIf
	TTElse
	TTLBrace
	TTRBrace
	TTLBracket
	TTRBracket
	TTTrue
	TTFalse
	TTQMark
	TTColon
	TTLog
	TTFn
	TTReturn
	TTArrow
	TTNil
	TTPipe
	TTColonColon
	TTWhile
	TTDollar
	TTEOF
)

var EmptyToken = Token{Type: TTEOF, Lexeme: ""}
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
	"fn",
	"return",
	"->",
	"nil",
	"|",
	"::",
	"while",
	"$",
	"eof",
}

func TokenTypeName(tt TokenType) string {
	return tokenTypeNames[tt]
}

// Token represents a lexed Token
type Token struct {
	Type       TokenType
	Lexeme     string
	LineNumber int
}

func (t Token) String() string {
	return t.Lexeme
}

func (t Token) ParseNumber() lang.Number {
	if t.Type == TTNumber {
		if n, err := strconv.ParseFloat(t.Lexeme, 64); err == nil {
			return lang.Number(n)
		}
	}

	panic(fmt.Sprintf("error converting %s to number", t.Lexeme))
}

func (t Token) ParseString() string {
	if t.Type == TTString {
		return strings.Trim(t.Lexeme, `"`)
	}

	panic(fmt.Sprintf("error converting %s to string", t.Lexeme))
}

func (t Token) ParseColor() lang.Color {
	if t.Type == TTColor {
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
			return lang.NewRgba(lang.Number(r), lang.Number(g), lang.Number(b), lang.Number(a))
		}
	}

	panic(fmt.Sprintf("error converting %s to color", t.Lexeme))
}

// Lex walks the specified string and returns an array of lexed Tokens
// or an non-nil error if the input could not be lexed.
func Lex(src string) ([]Token, error) {
	tokens := []Token{}
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

func match(src string) (Token, int) {
	for _, m := range matchers {
		if loc := m.regexp.FindStringIndex(src); loc != nil {
			lexeme := src[loc[0]:loc[1]]
			tok := Token{
				Type:   m.lookup(lexeme),
				Lexeme: lexeme,
			}
			return tok, loc[1]
		}
	}
	return EmptyToken, -1
}

type tokenLookup func(lexeme string) TokenType

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

var keywordTokens = map[string]TokenType{
	"and":    TTAnd,
	"or":     TTOr,
	"not":    TTNot,
	"for":    TTFor,
	"in":     TTIn,
	"yield":  TTYield,
	"if":     TTIf,
	"else":   TTElse,
	"true":   TTTrue,
	"false":  TTFalse,
	"log":    TTLog,
	"fn":     TTFn,
	"return": TTReturn,
	"nil":    TTNil,
	"while":  TTWhile,
}

func lookupKeyword(lexeme string) TokenType {
	if token, ok := keywordTokens[lexeme]; ok {
		return token
	}
	return TTIdent
}

var matchers = []matcher{
	makeMatcher(`\(`, func(string) TokenType { return TTLParen }),
	makeMatcher(`\)`, func(string) TokenType { return TTRParen }),
	makeMatcher(`\+`, func(string) TokenType { return TTPlus }),
	makeMatcher(`\->`, func(string) TokenType { return TTArrow }),
	makeMatcher(`\-`, func(string) TokenType { return TTMinus }),
	makeMatcher(`\*`, func(string) TokenType { return TTStar }),
	makeMatcher(`/`, func(string) TokenType { return TTSlash }),
	makeMatcher(`%`, func(string) TokenType { return TTPercent }),
	makeMatcher(`,`, func(string) TokenType { return TTComma }),
	makeMatcher(`\:=`, func(string) TokenType { return TTColonEq }),
	makeMatcher(`\:\:`, func(string) TokenType { return TTColonColon }),
	makeMatcher(`\:`, func(string) TokenType { return TTColon }),
	makeMatcher(`\[`, func(string) TokenType { return TTLBracket }),
	makeMatcher(`\]`, func(string) TokenType { return TTRBracket }),
	makeMatcher(`==`, func(string) TokenType { return TTEqEq }),
	makeMatcher(`=`, func(string) TokenType { return TTEq }),
	makeMatcher(`!=`, func(string) TokenType { return TTNeq }),
	makeMatcher(`>=`, func(string) TokenType { return TTGe }),
	makeMatcher(`>`, func(string) TokenType { return TTGt }),
	makeMatcher(`<=`, func(string) TokenType { return TTLe }),
	makeMatcher(`<`, func(string) TokenType { return TTLt }),
	makeMatcher(`@`, func(string) TokenType { return TTAt }),
	makeMatcher(`\.\.`, func(string) TokenType { return TTDotDot }),
	makeMatcher(`\.`, func(string) TokenType { return TTDot }),
	makeMatcher(`;`, func(string) TokenType { return TTSemicolon }),
	makeMatcher(`{`, func(string) TokenType { return TTLBrace }),
	makeMatcher(`}`, func(string) TokenType { return TTRBrace }),
	makeMatcher(`\?`, func(string) TokenType { return TTQMark }),
	makeMatcher(`\|`, func(string) TokenType { return TTPipe }),
	makeMatcher(`\$`, func(string) TokenType { return TTDollar }),
	makeMatcher(`#[0-9a-fA-F]{6}(\:[0-9a-fA-F]{2})?`, func(string) TokenType { return TTColor }),
	makeMatcher(`".*?"`, func(string) TokenType { return TTString }),
	makeMatcher(`[0-9]+(\.[0-9]+)?\b`, func(string) TokenType { return TTNumber }),
	makeMatcher(`[a-zA-Z_][a-zA-Z0-9_]*\b`, lookupKeyword),
}
