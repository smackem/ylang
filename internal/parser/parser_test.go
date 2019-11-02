package parser

import (
	"github.com/smackem/ylang/internal/lang"
	"github.com/smackem/ylang/internal/lexer"
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
			tokens, _ := lexer.Lex(tt.src)
			if _, err := Parse(tokens, true); (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parse_ast(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    []Statement
		wantErr bool
	}{
		{
			name: "declaration",
			src:  "x := 1",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Rhs:      lang.Number(1),
				},
			},
		},
		{
			name: "pixel_assign",
			src:  "x := sort(a(1)) @p = 2",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Rhs: InvokeExpr{
						FuncName: "sort",
						Args: []Expression{
							InvokeExpr{
								FuncName: "a",
								Args: []Expression{
									lang.Number(1),
								},
							},
						},
					},
				},
				PixelAssignStmt{
					StmtBase: StmtBase{},
					Lhs:      IdentExpr("p"),
					Rhs:      lang.Number(2),
				},
			},
		},
		{
			name: "indexedAssign",
			src:  "x[1] = 2",
			want: []Statement{
				IndexedAssignStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Index:    lang.Number(1),
					Rhs:      lang.Number(2),
				},
			},
		},
		{
			name: "functionDecl",
			src:  "f := fn(x) { return 1 }",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "f",
					Rhs: FunctionExpr{
						ParameterNames: []string{"x"},
						Body: []Statement{
							ReturnStmt{
								StmtBase: StmtBase{},
								Result:   lang.Number(1),
							},
						},
					},
				},
			},
		},
		{
			name: "functionDecl_Lambda",
			src:  "f := fn() -> 5",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "f",
					Rhs: FunctionExpr{
						ParameterNames: nil,
						Body: []Statement{
							ReturnStmt{
								StmtBase: StmtBase{},
								Result:   lang.Number(5),
							},
						},
					},
				},
			},
		},
		{
			name: "return",
			src:  "return 100",
			want: []Statement{
				ReturnStmt{
					StmtBase: StmtBase{},
					Result:   lang.Number(100),
				},
			},
		},
		{
			name: "multiple_statements",
			src:  "log(1) log(2)",
			want: []Statement{
				LogStmt{
					StmtBase: StmtBase{},
					Args: []Expression{
						lang.Number(1),
					},
				},
				LogStmt{
					StmtBase: StmtBase{},
					Args: []Expression{
						lang.Number(2),
					},
				},
			},
		},
		{
			name: "parameter_list",
			src:  "log(1, 2, 3)",
			want: []Statement{
				LogStmt{
					StmtBase: StmtBase{},
					Args: []Expression{
						lang.Number(1),
						lang.Number(2),
						lang.Number(3),
					},
				},
			},
		},
		{
			name: "list_empty",
			src:  "l := []",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "l",
					Rhs:      ListExpr{},
				},
			},
		},
		{
			name: "list",
			src:  "l := [1,2,3]",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "l",
					Rhs: ListExpr{
						Elements: []Expression{lang.Number(1), lang.Number(2), lang.Number(3)},
					},
				},
			},
		},
		{
			name: "list_trailing_comma",
			src:  "l := [1,]",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "l",
					Rhs: ListExpr{
						Elements: []Expression{lang.Number(1)},
					},
				},
			},
		},
		{
			name: "molecules",
			src:  "x := @(1;2).r",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Rhs: MemberExpr{
						Member: "r",
						Recvr: AtExpr{
							Inner: PosExpr{
								X: lang.Number(1),
								Y: lang.Number(2),
							},
						},
					},
				},
			},
		},
		{
			name: "molecules_2",
			src:  "x := k[1;2].m",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Rhs: MemberExpr{
						Member: "m",
						Recvr: IndexExpr{
							Recvr: IdentExpr("k"),
							Index: PosExpr{
								X: lang.Number(1),
								Y: lang.Number(2),
							},
						},
					},
				},
			},
		},
		{
			name: "concat",
			src:  "l := [1] :: 2 :: 3",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "l",
					Rhs: ConcatExpr{
						Left: ConcatExpr{
							Left: ListExpr{
								Elements: []Expression{lang.Number(1)},
							},
							Right: lang.Number(2),
						},
						Right: lang.Number(3),
					},
				},
			},
		},
		{
			name: "indexRange",
			src:  "s := ls[1..10]",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "s",
					Rhs: IndexRangeExpr{
						Recvr: IdentExpr("ls"),
						Lower: lang.Number(1),
						Upper: lang.Number(10),
					},
				},
			},
		},
		{
			name: "term",
			src:  "x := 1 + 3 - 2",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Rhs: SubExpr{
						Left: AddExpr{
							Left:  lang.Number(1),
							Right: lang.Number(3),
						},
						Right: lang.Number(2),
					},
				},
			},
		},
		{
			name: "product",
			src:  "x := 16 / 4 * 2",
			want: []Statement{
				DeclStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Rhs: MulExpr{
						Left: DivExpr{
							Left:  lang.Number(16),
							Right: lang.Number(4),
						},
						Right: lang.Number(2),
					},
				},
			},
		},
		{
			name: "members",
			src:  "log(x.member1.member2.member3)",
			want: []Statement{
				LogStmt{
					StmtBase: StmtBase{},
					Args: []Expression{
						MemberExpr{
							Member: "member3",
							Recvr: MemberExpr{
								Member: "member2",
								Recvr: MemberExpr{
									Member: "member1",
									Recvr:  IdentExpr("x"),
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
			want: []Statement{
				LogStmt{
					StmtBase: StmtBase{},
					Args:     []Expression{lang.NewRgba(0xff, 0xee, 0x44, 0x0f)},
				},
			},
		},
		{
			name: "function_without_params",
			src:  "log(fun())",
			want: []Statement{
				LogStmt{
					Args: []Expression{
						InvokeExpr{
							FuncName: "fun",
							Args:     nil,
						},
					},
				},
			},
		},
		{
			name: "functions",
			src:  "log(sort(map_b(1)))",
			want: []Statement{
				LogStmt{
					Args: []Expression{
						InvokeExpr{
							FuncName: "sort",
							Args: []Expression{
								InvokeExpr{
									FuncName: "map_b",
									Args: []Expression{
										lang.Number(1),
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
			want: []Statement{
				IfStmt{
					StmtBase: StmtBase{},
					Cond:     lang.TrueVal,
					TrueStmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
					FalseStmts: nil,
				},
			},
		},
		{
			name: "for_range",
			src:  "for x in 0..10 { log(1) }",
			want: []Statement{
				ForRangeStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Lower:    lang.Number(0),
					Step:     lang.Number(1),
					Upper:    lang.Number(10),
					Stmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
				},
			},
		},
		{
			name: "for_range_with_step",
			src:  "for x in 0..2..10 { log(1) }",
			want: []Statement{
				ForRangeStmt{
					StmtBase: StmtBase{},
					Ident:    "x",
					Lower:    lang.Number(0),
					Step:     lang.Number(2),
					Upper:    lang.Number(10),
					Stmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
				},
			},
		},
		{
			name: "for",
			src:  "for x in coll { log(1) }",
			want: []Statement{
				ForStmt{
					StmtBase:   StmtBase{},
					Ident:      "x",
					Collection: IdentExpr("coll"),
					Stmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
				},
			},
		},
		{
			name: "if_else",
			src:  "if true { log(1) } else { log(2) }",
			want: []Statement{
				IfStmt{
					StmtBase: StmtBase{},
					Cond:     lang.TrueVal,
					TrueStmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
					FalseStmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(2)},
						},
					},
				},
			},
		},
		{
			name: "if_elseif",
			src:  "if true { log(1) } else if false { log(2) }",
			want: []Statement{
				IfStmt{
					StmtBase: StmtBase{},
					Cond:     lang.TrueVal,
					TrueStmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
					FalseStmts: []Statement{
						IfStmt{
							StmtBase: StmtBase{},
							Cond:     lang.FalseVal,
							TrueStmts: []Statement{
								LogStmt{
									StmtBase: StmtBase{},
									Args:     []Expression{lang.Number(2)},
								},
							},
							FalseStmts: nil,
						},
					},
				},
			},
		},
		{
			name: "if_elseif_else",
			src:  "if true { log(1) } else if false { log(2) } else { log(3) }",
			want: []Statement{
				IfStmt{
					StmtBase: StmtBase{},
					Cond:     lang.TrueVal,
					TrueStmts: []Statement{
						LogStmt{
							StmtBase: StmtBase{},
							Args:     []Expression{lang.Number(1)},
						},
					},
					FalseStmts: []Statement{
						IfStmt{
							StmtBase: StmtBase{},
							Cond:     lang.FalseVal,
							TrueStmts: []Statement{
								LogStmt{
									StmtBase: StmtBase{},
									Args:     []Expression{lang.Number(2)},
								},
							},
							FalseStmts: []Statement{
								LogStmt{
									StmtBase: StmtBase{},
									Args:     []Expression{lang.Number(3)},
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
			want: []Statement{
				WhileStmt{
					Cond: lang.TrueVal,
					Stmts: []Statement{
						LogStmt{
							Args: []Expression{lang.Number(1)},
						},
					},
				},
			},
		},
		{
			name: "while_2",
			src:  "while x < y { log(100) }",
			want: []Statement{
				WhileStmt{
					Cond: LtExpr{
						Left:  IdentExpr("x"),
						Right: IdentExpr("y"),
					},
					Stmts: []Statement{
						LogStmt{
							Args: []Expression{lang.Number(100)},
						},
					},
				},
			},
		},
		{
			name: "pipeline",
			src:  "log(1 | $+1 | $+2 | $+3)",
			want: []Statement{
				LogStmt{
					Args: []Expression{
						PipelineExpr{
							Left: lang.Number(1),
							Right: PipelineExpr{
								Left: AddExpr{
									Left:  IdentExpr("$"),
									Right: lang.Number(1),
								},
								Right: PipelineExpr{
									Left: AddExpr{
										Left:  IdentExpr("$"),
										Right: lang.Number(2),
									},
									Right: AddExpr{
										Left:  IdentExpr("$"),
										Right: lang.Number(3),
									},
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
			tokens, _ := lexer.Lex(tt.src)
			got, err := Parse(tokens, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Stmts, tt.want) {
				t.Errorf("parse() =\n%#v\nwant\n%#v", got.Stmts, tt.want)
			}
		})
	}
}
