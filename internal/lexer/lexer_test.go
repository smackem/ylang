package lexer

import (
	"github.com/smackem/ylang/internal/lang"
	"reflect"
	"testing"
)

func Test_lex(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    []Token
		wantErr bool
	}{
		{
			name: "declaration",
			src:  "x := 1",
			want: []Token{
				Token{
					Type:       TTIdent,
					Lexeme:     "x",
					LineNumber: 1,
				},
				Token{
					Type:       TTColonEq,
					Lexeme:     ":=",
					LineNumber: 1,
				},
				Token{
					Type:       TTNumber,
					Lexeme:     "1",
					LineNumber: 1,
				},
			},
		},
		{
			name: "linenumbers",
			src:  "a\nb\nc",
			want: []Token{
				Token{
					Type:       TTIdent,
					Lexeme:     "a",
					LineNumber: 1,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "b",
					LineNumber: 2,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "c",
					LineNumber: 3,
				},
			},
		},
		{
			name: "numbers",
			src:  "12.5 999",
			want: []Token{
				Token{
					Type:       TTNumber,
					Lexeme:     "12.5",
					LineNumber: 1,
				},
				Token{
					Type:       TTNumber,
					Lexeme:     "999",
					LineNumber: 1,
				},
			},
		},
		{
			name: "comment",
			src:  "x\ny//comment\nz",
			want: []Token{
				Token{
					Type:       TTIdent,
					Lexeme:     "x",
					LineNumber: 1,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "y",
					LineNumber: 2,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "z",
					LineNumber: 3,
				},
			},
		},
		{
			name: "strings",
			src:  "\"\" \"hepp\"",
			want: []Token{
				Token{
					Type:       TTString,
					Lexeme:     "\"\"",
					LineNumber: 1,
				},
				Token{
					Type:       TTString,
					Lexeme:     "\"hepp\"",
					LineNumber: 1,
				},
			},
		},
		{
			name: "colors",
			src:  "#1a2b3c #1F2E3D:c0",
			want: []Token{
				Token{
					Type:       TTColor,
					Lexeme:     "#1a2b3c",
					LineNumber: 1,
				},
				Token{
					Type:       TTColor,
					Lexeme:     "#1F2E3D:c0",
					LineNumber: 1,
				},
			},
		},
		{
			name: "pipeline",
			src:  "1 | $ + 2",
			want: []Token{
				Token{
					Type:       TTNumber,
					Lexeme:     "1",
					LineNumber: 1,
				},
				Token{
					Type:       TTPipe,
					Lexeme:     "|",
					LineNumber: 1,
				},
				Token{
					Type:       TTDollar,
					Lexeme:     "$",
					LineNumber: 1,
				},
				Token{
					Type:       TTPlus,
					Lexeme:     "+",
					LineNumber: 1,
				},
				Token{
					Type:       TTNumber,
					Lexeme:     "2",
					LineNumber: 1,
				},
			},
		},
		{
			name: "for",
			src:  "for pos in IMAGE {\n    @pos = -@pos\n}",
			want: []Token{
				Token{
					Type:       TTFor,
					Lexeme:     "for",
					LineNumber: 1,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "pos",
					LineNumber: 1,
				},
				Token{
					Type:       TTIn,
					Lexeme:     "in",
					LineNumber: 1,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "IMAGE",
					LineNumber: 1,
				},
				Token{
					Type:       TTLBrace,
					Lexeme:     "{",
					LineNumber: 1,
				},
				Token{
					Type:       TTAt,
					Lexeme:     "@",
					LineNumber: 2,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "pos",
					LineNumber: 2,
				},
				Token{
					Type:       TTEq,
					Lexeme:     "=",
					LineNumber: 2,
				},
				Token{
					Type:       TTMinus,
					Lexeme:     "-",
					LineNumber: 2,
				},
				Token{
					Type:       TTAt,
					Lexeme:     "@",
					LineNumber: 2,
				},
				Token{
					Type:       TTIdent,
					Lexeme:     "pos",
					LineNumber: 2,
				},
				Token{
					Type:       TTRBrace,
					Lexeme:     "}",
					LineNumber: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Lex(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_token_parseColor(t *testing.T) {
	tests := []struct {
		name string
		t    Token
		want lang.Color
	}{
		{
			name: "rgb",
			t: Token{
				Type:   TTColor,
				Lexeme: "#1a2b3c",
			},
			want: lang.NewRgba(lang.Number(0x1a), lang.Number(0x2b), lang.Number(0x3c), lang.Number(255)),
		},
		{
			name: "rgb",
			t: Token{
				Type:   TTColor,
				Lexeme: "#1F2E3D:c0",
			},
			want: lang.NewRgba(lang.Number(0x1f), lang.Number(0x2e), lang.Number(0x3d), lang.Number(0xc0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.ParseColor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Token.parseColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
