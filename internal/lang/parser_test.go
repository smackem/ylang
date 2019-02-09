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
		{
			name:    "error_constant_assignment",
			src:     "NUM := 1 NUM = 2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lex(tt.src)
			if _, err := parse(tokens, true); (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parse_ast(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    []statement
		wantErr bool
	}{
		{
			name: "declaration",
			src:  "x := 1",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					rhs:      Number(1),
				},
			},
		},
		{
			name: "pixel_assign",
			src:  "x := sort(a(1)) @p = 2",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					rhs: invokeExpr{
						funcName: "sort",
						parameters: []expression{
							invokeExpr{
								funcName: "a",
								parameters: []expression{
									Number(1),
								},
							},
						},
					},
				},
				pixelAssignStmt{
					stmtBase: stmtBase{},
					lhs:      identExpr("p"),
					rhs:      Number(2),
				},
			},
		},
		{
			name: "multiple_statements",
			src:  "log(1) commit",
			want: []statement{
				logStmt{
					stmtBase: stmtBase{},
					parameters: []expression{
						Number(1),
					},
				},
				commitStmt{
					stmtBase: stmtBase{},
					rect:     nil,
				},
			},
		},
		{
			name: "parameter_list",
			src:  "log(1, 2, 3)",
			want: []statement{
				logStmt{
					stmtBase: stmtBase{},
					parameters: []expression{
						Number(1),
						Number(2),
						Number(3),
					},
				},
			},
		},
		{
			name: "blt",
			src:  "blt(BOUNDS) log(1)",
			want: []statement{
				bltStmt{
					stmtBase: stmtBase{},
					rect:     identExpr("BOUNDS"),
				},
				logStmt{
					stmtBase:   stmtBase{},
					parameters: []expression{Number(1)},
				},
			},
		},
		{
			name: "blt_with_rect",
			src:  "blt(IMAGE) log(1)",
			want: []statement{
				bltStmt{
					stmtBase: stmtBase{},
					rect:     identExpr("IMAGE"),
				},
				logStmt{
					stmtBase:   stmtBase{},
					parameters: []expression{Number(1)},
				},
			},
		},
		{
			name: "molecules",
			src:  "x := @(1;2).r",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "x",
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
		{
			name: "molecules_2",
			src:  "x := k[1;2].m",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					rhs: memberExpr{
						member: "m",
						recvr: indexExpr{
							recvr: identExpr("k"),
							index: posExpr{
								x: Number(1),
								y: Number(2),
							},
						},
					},
				},
			},
		},
		{
			name: "color_literal",
			src:  "log(#ffee44:0f)",
			want: []statement{
				logStmt{
					stmtBase:   stmtBase{},
					parameters: []expression{NewRgba(0xff, 0xee, 0x44, 0x0f)},
				},
			},
		},
		{
			name: "functions",
			src:  "log(sort(map_b(1)))",
			want: []statement{
				logStmt{
					parameters: []expression{
						invokeExpr{
							funcName: "sort",
							parameters: []expression{
								invokeExpr{
									funcName: "map_b",
									parameters: []expression{
										Number(1),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "if",
			src:  "if true { log(1) }",
			want: []statement{
				ifStmt{
					stmtBase: stmtBase{},
					cond:     Bool(true),
					trueStmts: []statement{
						logStmt{
							stmtBase:   stmtBase{},
							parameters: []expression{Number(1)},
						},
					},
					falseStmts: nil,
				},
			},
		},
		{
			name: "if_else",
			src:  "if true { log(1) } else { log(2) }",
			want: []statement{
				ifStmt{
					stmtBase: stmtBase{},
					cond:     Bool(true),
					trueStmts: []statement{
						logStmt{
							stmtBase:   stmtBase{},
							parameters: []expression{Number(1)},
						},
					},
					falseStmts: []statement{
						logStmt{
							stmtBase:   stmtBase{},
							parameters: []expression{Number(2)},
						},
					},
				},
			},
		},
		{
			name: "if_elseif",
			src:  "if true { log(1) } else if false { log(2) }",
			want: []statement{
				ifStmt{
					stmtBase: stmtBase{},
					cond:     Bool(true),
					trueStmts: []statement{
						logStmt{
							stmtBase:   stmtBase{},
							parameters: []expression{Number(1)},
						},
					},
					falseStmts: []statement{
						ifStmt{
							stmtBase: stmtBase{},
							cond:     Bool(false),
							trueStmts: []statement{
								logStmt{
									stmtBase:   stmtBase{},
									parameters: []expression{Number(2)},
								},
							},
							falseStmts: nil,
						},
					},
				},
			},
		},
		{
			name: "if_elseif_else",
			src:  "if true { log(1) } else if false { log(2) } else { log(3) }",
			want: []statement{
				ifStmt{
					stmtBase: stmtBase{},
					cond:     Bool(true),
					trueStmts: []statement{
						logStmt{
							stmtBase:   stmtBase{},
							parameters: []expression{Number(1)},
						},
					},
					falseStmts: []statement{
						ifStmt{
							stmtBase: stmtBase{},
							cond:     Bool(false),
							trueStmts: []statement{
								logStmt{
									stmtBase:   stmtBase{},
									parameters: []expression{Number(2)},
								},
							},
							falseStmts: []statement{
								logStmt{
									stmtBase:   stmtBase{},
									parameters: []expression{Number(3)},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, _ := lex(tt.src)
			got, err := parse(tokens, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.stmts, tt.want) {
				t.Errorf("parse() = %#v, want %#v", got.stmts, tt.want)
			}
		})
	}
}
