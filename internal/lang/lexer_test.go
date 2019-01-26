package lang

import (
	"reflect"
	"testing"
)

func Test_lex(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    []token
		wantErr bool
	}{
		{
			name: "declaration",
			src:  "x := 1",
			want: []token{
				token{
					Type:       ttIdent,
					Lexeme:     "x",
					LineNumber: 1,
				},
				token{
					Type:       ttColonEq,
					Lexeme:     ":=",
					LineNumber: 1,
				},
				token{
					Type:       ttNumber,
					Lexeme:     "1",
					LineNumber: 1,
				},
			},
		},
		{
			name: "linenumbers",
			src:  "a\nb\nc",
			want: []token{
				token{
					Type:       ttIdent,
					Lexeme:     "a",
					LineNumber: 1,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "b",
					LineNumber: 2,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "c",
					LineNumber: 3,
				},
			},
		},
		{
			name: "numbers",
			src:  "12.5 999",
			want: []token{
				token{
					Type:       ttNumber,
					Lexeme:     "12.5",
					LineNumber: 1,
				},
				token{
					Type:       ttNumber,
					Lexeme:     "999",
					LineNumber: 1,
				},
			},
		},
		{
			name: "comment",
			src:  "x\ny//comment\nz",
			want: []token{
				token{
					Type:       ttIdent,
					Lexeme:     "x",
					LineNumber: 1,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "y",
					LineNumber: 2,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "z",
					LineNumber: 3,
				},
			},
		},
		{
			name: "strings",
			src:  "\"\" \"hepp\"",
			want: []token{
				token{
					Type:       ttString,
					Lexeme:     "\"\"",
					LineNumber: 1,
				},
				token{
					Type:       ttString,
					Lexeme:     "\"hepp\"",
					LineNumber: 1,
				},
			},
		},
		{
			name: "colors",
			src:  "#1a2b3c #1F2E3D:c0",
			want: []token{
				token{
					Type:       ttColor,
					Lexeme:     "#1a2b3c",
					LineNumber: 1,
				},
				token{
					Type:       ttColor,
					Lexeme:     "#1F2E3D:c0",
					LineNumber: 1,
				},
			},
		},
		{
			name: "for",
			src:  "for pos in IMAGE {\n    @pos = -@pos\n}",
			want: []token{
				token{
					Type:       ttFor,
					Lexeme:     "for",
					LineNumber: 1,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "pos",
					LineNumber: 1,
				},
				token{
					Type:       ttIn,
					Lexeme:     "in",
					LineNumber: 1,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "IMAGE",
					LineNumber: 1,
				},
				token{
					Type:       ttLBrace,
					Lexeme:     "{",
					LineNumber: 1,
				},
				token{
					Type:       ttAt,
					Lexeme:     "@",
					LineNumber: 2,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "pos",
					LineNumber: 2,
				},
				token{
					Type:       ttEq,
					Lexeme:     "=",
					LineNumber: 2,
				},
				token{
					Type:       ttMinus,
					Lexeme:     "-",
					LineNumber: 2,
				},
				token{
					Type:       ttAt,
					Lexeme:     "@",
					LineNumber: 2,
				},
				token{
					Type:       ttIdent,
					Lexeme:     "pos",
					LineNumber: 2,
				},
				token{
					Type:       ttRBrace,
					Lexeme:     "}",
					LineNumber: 3,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lex(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("lex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_token_parseColor(t *testing.T) {
	tests := []struct {
		name string
		t    token
		want Color
	}{
		{
			name: "rgb",
			t: token{
				Type:   ttColor,
				Lexeme: "#1a2b3c",
			},
			want: NewRgba(Number(0x1a), Number(0x2b), Number(0x3c), Number(255)),
		},
		{
			name: "rgb",
			t: token{
				Type:   ttColor,
				Lexeme: "#1F2E3D:c0",
			},
			want: NewRgba(Number(0x1f), Number(0x2e), Number(0x3d), Number(0xc0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.parseColor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("token.parseColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
