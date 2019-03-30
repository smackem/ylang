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
			src:     "k := |1 2 3 4 5 6 7 8 9|",
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
						args: []expression{
							invokeExpr{
								funcName: "a",
								args: []expression{
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
			name: "indexedAssign",
			src:  "x[1] = 2",
			want: []statement{
				indexedAssignStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					index:    Number(1),
					rhs:      Number(2),
				},
			},
		},
		{
			name: "functionDecl",
			src:  "f := fn(x) { return 1 }",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "f",
					rhs: functionExpr{
						parameterNames: []string{"x"},
						body: []statement{
							returnStmt{
								stmtBase: stmtBase{},
								result:   Number(1),
							},
						},
					},
				},
			},
		},
		{
			name: "functionDecl_Lambda",
			src:  "f := fn() -> 5",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "f",
					rhs: functionExpr{
						parameterNames: nil,
						body: []statement{
							returnStmt{
								stmtBase: stmtBase{},
								result:   Number(5),
							},
						},
					},
				},
			},
		},
		{
			name: "return",
			src:  "return 100",
			want: []statement{
				returnStmt{
					stmtBase: stmtBase{},
					result:   Number(100),
				},
			},
		},
		{
			name: "multiple_statements",
			src:  "log(1) log(2)",
			want: []statement{
				logStmt{
					stmtBase: stmtBase{},
					args: []expression{
						Number(1),
					},
				},
				logStmt{
					stmtBase: stmtBase{},
					args: []expression{
						Number(2),
					},
				},
			},
		},
		{
			name: "parameter_list",
			src:  "log(1, 2, 3)",
			want: []statement{
				logStmt{
					stmtBase: stmtBase{},
					args: []expression{
						Number(1),
						Number(2),
						Number(3),
					},
				},
			},
		},
		{
			name: "list_empty",
			src:  "l := []",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "l",
					rhs:      listExpr{},
				},
			},
		},
		{
			name: "list",
			src:  "l := [1,2,3]",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "l",
					rhs: listExpr{
						elements: []expression{Number(1), Number(2), Number(3)},
					},
				},
			},
		},
		{
			name: "list_trailing_comma",
			src:  "l := [1,]",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "l",
					rhs: listExpr{
						elements: []expression{Number(1)},
					},
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
			name: "concat",
			src:  "l := [1] :: 2 :: 3",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "l",
					rhs: concatExpr{
						left: concatExpr{
							left: listExpr{
								elements: []expression{Number(1)},
							},
							right: Number(2),
						},
						right: Number(3),
					},
				},
			},
		},
		{
			name: "indexRange",
			src:  "s := ls[1..10]",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "s",
					rhs: indexRangeExpr{
						recvr: identExpr("ls"),
						lower: Number(1),
						upper: Number(10),
					},
				},
			},
		},
		{
			name: "term",
			src:  "x := 1 + 3 - 2",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					rhs: subExpr{
						left: addExpr{
							left:  Number(1),
							right: Number(3),
						},
						right: Number(2),
					},
				},
			},
		},
		{
			name: "product",
			src:  "x := 16 / 4 * 2",
			want: []statement{
				declStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					rhs: mulExpr{
						left: divExpr{
							left:  Number(16),
							right: Number(4),
						},
						right: Number(2),
					},
				},
			},
		},
		{
			name: "members",
			src:  "log(x.member1.member2.member3)",
			want: []statement{
				logStmt{
					stmtBase: stmtBase{},
					args: []expression{
						memberExpr{
							member: "member3",
							recvr: memberExpr{
								member: "member2",
								recvr: memberExpr{
									member: "member1",
									recvr:  identExpr("x"),
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
			want: []statement{
				logStmt{
					stmtBase: stmtBase{},
					args:     []expression{NewRgba(0xff, 0xee, 0x44, 0x0f)},
				},
			},
		},
		{
			name: "function_without_params",
			src:  "log(fun())",
			want: []statement{
				logStmt{
					args: []expression{
						invokeExpr{
							funcName: "fun",
							args:     nil,
						},
					},
				},
			},
		},
		{
			name: "functions",
			src:  "log(sort(map_b(1)))",
			want: []statement{
				logStmt{
					args: []expression{
						invokeExpr{
							funcName: "sort",
							args: []expression{
								invokeExpr{
									funcName: "map_b",
									args: []expression{
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
					cond:     boolean(true),
					trueStmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
					falseStmts: nil,
				},
			},
		},
		{
			name: "for_range",
			src:  "for x in 0..10 { log(1) }",
			want: []statement{
				forRangeStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					lower:    Number(0),
					step:     Number(1),
					upper:    Number(10),
					stmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
				},
			},
		},
		{
			name: "for_range_with_step",
			src:  "for x in 0..2..10 { log(1) }",
			want: []statement{
				forRangeStmt{
					stmtBase: stmtBase{},
					ident:    "x",
					lower:    Number(0),
					step:     Number(2),
					upper:    Number(10),
					stmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
				},
			},
		},
		{
			name: "for",
			src:  "for x in coll { log(1) }",
			want: []statement{
				forStmt{
					stmtBase:   stmtBase{},
					ident:      "x",
					collection: identExpr("coll"),
					stmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
				},
			},
		},
		{
			name: "if_else",
			src:  "if true { log(1) } else { log(2) }",
			want: []statement{
				ifStmt{
					stmtBase: stmtBase{},
					cond:     boolean(true),
					trueStmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
					falseStmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(2)},
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
					cond:     boolean(true),
					trueStmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
					falseStmts: []statement{
						ifStmt{
							stmtBase: stmtBase{},
							cond:     boolean(false),
							trueStmts: []statement{
								logStmt{
									stmtBase: stmtBase{},
									args:     []expression{Number(2)},
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
					cond:     boolean(true),
					trueStmts: []statement{
						logStmt{
							stmtBase: stmtBase{},
							args:     []expression{Number(1)},
						},
					},
					falseStmts: []statement{
						ifStmt{
							stmtBase: stmtBase{},
							cond:     boolean(false),
							trueStmts: []statement{
								logStmt{
									stmtBase: stmtBase{},
									args:     []expression{Number(2)},
								},
							},
							falseStmts: []statement{
								logStmt{
									stmtBase: stmtBase{},
									args:     []expression{Number(3)},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "while",
			src:  "while true { log(1) }",
			want: []statement{
				whileStmt{
					cond: boolean(true),
					stmts: []statement{
						logStmt{
							args: []expression{Number(1)},
						},
					},
				},
			},
		},
		{
			name: "while_2",
			src:  "while x < y { log(100) }",
			want: []statement{
				whileStmt{
					cond: ltExpr{
						left:  identExpr("x"),
						right: identExpr("y"),
					},
					stmts: []statement{
						logStmt{
							args: []expression{Number(100)},
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
				t.Errorf("parse() =\n%#v\nwant\n%#v", got.stmts, tt.want)
			}
		})
	}
}
