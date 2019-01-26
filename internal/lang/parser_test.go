package lang

import (
	"reflect"
	"testing"
)

func Test_parse_syntax(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{
			name:    "decl",
			src:     "x := 1",
			wantErr: false,
		},
		{
			name:    "assignment",
			src:     "x = 2.0",
			wantErr: false,
		},
		{
			name:    "if",
			src:     "if x > 0 { log(2.0) }",
			wantErr: false,
		},
		{
			name:    "ifelse",
			src:     "if x > 0 { log(2.0) } else { log(3.0) }",
			wantErr: false,
		},
		{
			name:    "for",
			src:     "for pos in IMAGE { @pos = @pos }",
			wantErr: false,
		},
		{
			name:    "empty",
			src:     "",
			wantErr: false,
		},
		{
			name:    "double_operator",
			src:     "if x * > 3 { log(1) }",
			wantErr: true,
		},
		{
			name:    "missing_rbrace",
			src:     "if x > 3 { log(1)",
			wantErr: true,
		},
		{
			name:    "multiple_stmts",
			src:     "x := 1 x = 2 x = 3",
			wantErr: false,
		},
		{
			name:    "parameter_list",
			src:     "log(1, 2, 3, 4)",
			wantErr: false,
		},
		{
			name:    "kernel_expr",
			src:     "k := [1 2 3 4 5 6 7 8 9]",
			wantErr: false,
		},
		{
			name:    "multiple_stmt_lists",
			src:     "if true { log(1) log(2) } if false { log(3) }",
			wantErr: false,
		},
		{
			name:    "functions",
			src:     "r := rect(1, 2, 3, 4) c := rgb(255, 254, 254)",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lex(tt.src)
			if _, err := parse(tokens); (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parse_ast(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    Program
		wantErr bool
	}{
		{
			name: "declaration",
			src:  "x := 1",
			want: Program{
				[]statement{
					declStmt{
						ident: "x",
						rhs:   Number(1),
					},
				},
			},
		},
		{
			name: "multiple_statements",
			src:  "log(1) blt",
			want: Program{
				[]statement{
					logStmt{
						parameters: []expression{
							Number(1),
						},
					},
					bltStmt{rect: nil},
				},
			},
		},
		{
			name: "parameter_list",
			src:  "log(1, 2, 3)",
			want: Program{
				[]statement{
					logStmt{
						parameters: []expression{
							Number(1),
							Number(2),
							Number(3),
						},
					},
				},
			},
		},
		{
			name: "blt",
			src:  "blt log(1)",
			want: Program{
				[]statement{
					bltStmt{},
					logStmt{
						parameters: []expression{Number(1)},
					},
				},
			},
		},
		{
			name: "blt_with_rect",
			src:  "blt(IMAGE) log(1)",
			want: Program{
				[]statement{
					bltStmt{
						rect: identExpr("IMAGE"),
					},
					logStmt{
						parameters: []expression{Number(1)},
					},
				},
			},
		},
		{
			name: "molecules",
			src:  "x := @(1;2).r",
			want: Program{
				[]statement{
					declStmt{
						ident: "x",
						rhs: memberExpr{
							member: "r",
							recvr: atExpr{
								inner: posExpr{
									x: Number(1),
									y: Number(2),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "color_literal",
			src:  "log(#ffee44:0f)",
			want: Program{
				[]statement{
					logStmt{
						parameters: []expression{NewRgba(0xff, 0xee, 0x44, 0x0f)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lex(tt.src)
			got, err := parse(tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
