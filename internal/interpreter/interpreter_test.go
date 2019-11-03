package interpreter

import (
	"github.com/smackem/ylang/internal/lang"
	"github.com/smackem/ylang/internal/lexer"
	"github.com/smackem/ylang/internal/parser"
	"image"
	"reflect"
	"testing"
)

func compileAndInterpret(src string) (scope, error) {
	tokens, err := lexer.Lex(src)
	if err != nil {
		return nil, err
	}
	program, err := parser.Parse(tokens, true)
	if err != nil {
		return nil, err
	}
	ir := newInterpreter(nil)
	err = ir.visitStmtList(program.Stmts)
	if err != nil {
		return nil, err
	}
	topScope := ir.idents[1]
	delete(topScope, lastRectIdent)
	return topScope, nil
}

func Test_interpret(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    scope
		wantErr bool
	}{
		{
			name: "simple",
			src:  "x := 1",
			want: scope{
				"x": number(1),
			},
		},
		{
			name: "multi_declaration",
			src: `x := 1
			      x := 2`,
			want: scope{
				"x": number(2),
			},
		},
		{
			name: "rect",
			src:  "x := rect(1, 2, 3, 4)",
			want: scope{
				"x": rect{Min: image.Point{1, 2}, Max: image.Point{4, 6}},
			},
		},
		{
			name: "colors",
			src: `c1 := rgb(1,2,3)
				  c2 := rgba(1,2,3,4)
				  r := c1.r
				  g := c1.g
				  b := c1.b
				  a := c1.a`,
			want: scope{
				"c1": color(lang.NewRgba(1, 2, 3, 255)),
				"c2": color(lang.NewRgba(1, 2, 3, 4)),
				"r":  number(1),
				"g":  number(2),
				"b":  number(3),
				"a":  number(255),
			},
		},
		{
			name: "scolors",
			src: `c1 := rgb01(1,2,3)
				  c2 := rgba01(1,2,3,4)
				  scr := c1.r01
				  scg := c1.g01
				  scb := c1.b01
				  sca := c1.a01`,
			want: scope{
				"c1":  color(lang.NewSrgba(1, 2, 3, 1)),
				"c2":  color(lang.NewSrgba(1, 2, 3, 4)),
				"scr": number(1),
				"scg": number(2),
				"scb": number(3),
				"sca": number(1),
			},
		},
		{
			name: "color_literal",
			src:  "c1 := #011223:f0",
			want: scope{
				"c1": color(lang.NewRgba(0x01, 0x12, 0x23, 0xf0)),
			},
		},
		{
			name: "kernel",
			src:  "k := |1 2 3 4|",
			want: scope{
				"k": kernel{width: 2, height: 2, values: []lang.Number{lang.Number(1), lang.Number(2), lang.Number(3), lang.Number(4)}},
			},
		},
		{
			name: "if_else",
			src: `x := 0
			      if true { x = 1 } else { x = 2 }`,
			want: scope{
				"x": number(1),
			},
		},
		{
			name: "if_else_2",
			src: `x := 0
			      if false { x = 1 } else { x = 2 }`,
			want: scope{
				"x": number(2),
			},
		},
		{
			name: "for",
			src: `p := 0
			      for pos in rect(0,0,1,1) { p = pos }`,
			want: scope{
				"p": point{0, 0},
			},
		},
		{
			name: "for_2",
			src: `p := 0
			      for pos in rect(0,0,2,2) { p = pos }`,
			want: scope{
				"p": point{1, 1},
			},
		},
		{
			name: "kernel_index",
			src:  "n := |1 2 3 4|[1]",
			want: scope{
				"n": number(2),
			},
		},
		{
			name: "kernel_index_2",
			src:  "n := |1 2 3 4|[0;1]",
			want: scope{
				"n": number(3),
			},
		},
		{
			name: "sort_kernel",
			src:  "k := sort(|4 1 3 2|)",
			want: scope{
				"k": kernel{
					width:  2,
					height: 2,
					values: []lang.Number{1, 2, 3, 4},
				},
			},
		},
		{
			name: "min_max",
			src: `min := min(|4 1 3 2|)
			      max := max(|4 1 3 2|)`,
			want: scope{
				"min": number(1),
				"max": number(4),
			},
		},
		{
			name: "list_func",
			src:  "l := list(4, 123)",
			want: scope{
				"l": list{
					elements: []Value{number(123), number(123), number(123), number(123)},
				},
			},
		},
		{
			name: "kernel_func",
			src:  "k := kernel(2, 3, 1)",
			want: scope{
				"k": kernel{
					width:  2,
					height: 3,
					values: []lang.Number{1, 1, 1, 1, 1, 1},
				},
			},
		},
		{
			name: "line_func",
			src:  "l := line(0;0, 100;100)",
			want: scope{
				"l": line{
					point1: point{0, 0},
					point2: point{100, 100},
				},
			},
		},
		{
			name: "line_props",
			src: `p1 := line(0;1, 100;101).p1
			      p2 := line(0;1, 100;101).p2`,
			want: scope{
				"p1": point{0, 1},
				"p2": point{100, 101},
			},
		},
		{
			name: "polygon_func",
			src:  "p := polygon(0;0, 100;0, 100;100, 0;0)",
			want: scope{
				"p": polygon{
					vertices: []point{
						point{0, 0},
						point{100, 0},
						point{100, 100},
					},
				},
			},
		},
		{
			name: "polygon_bounds",
			src:  "b := polygon(0;0, 100;0, 100;100, 0;0).bounds",
			want: scope{
				"b": rect{
					Min: image.Point{0, 0},
					Max: image.Point{100, 100},
				},
			},
		},
		{
			name: "polygon_vertices",
			src:  "vs := polygon(0;0, 100;0, 100;100).vertices",
			want: scope{
				"vs": list{
					elements: []Value{
						point{0, 0},
						point{100, 0},
						point{100, 100},
					},
				},
			},
		},
		{
			name: "polygon_vertices_count",
			src:  "c := polygon(0;0, 100;0, 100;100).vertices.count",
			want: scope{
				"c": number(3),
			},
		},
		{
			name: "indexed_assign_kernel",
			src: `k := |1 2 3 4|
				  k[0] = 0`,
			want: scope{
				"k": kernel{width: 2, height: 2, values: []lang.Number{0, 2, 3, 4}},
			},
		},
		{
			name: "indexed_assign_kernel_neg",
			src: `k := |1 2 3 4|
				  k[-1] = 0`,
			want: scope{
				"k": kernel{width: 2, height: 2, values: []lang.Number{1, 2, 3, 0}},
			},
		},
		{
			name: "hashmap",
			src:  `m := {a: 1, b: 2, c: 3}`,
			want: scope{
				"m": hashMap{str("a"): number(1), str("b"): number(2), str("c"): number(3)},
			},
		},
		{
			name: "hashmap_index",
			src: `m := {a: 1, b: 2, c: 3,}
				  a := m.a
				  b := m.b
				  c := m["c"]`,
			want: scope{
				"m": hashMap{str("a"): number(1), str("b"): number(2), str("c"): number(3)},
				"a": number(1),
				"b": number(2),
				"c": number(3),
			},
		},
		{
			name: "hashmap_index_2",
			src: `m := {a: 1}
				  a := m.a
				  a1 := m["a"]
				  b := m["b"]`,
			want: scope{
				"m":  hashMap{str("a"): number(1)},
				"a":  number(1),
				"a1": number(1),
				"b":  nilval{},
			},
		},
		{
			name: "hashmap_indexed_assign",
			src: `m := {}
				  m["a"] = 123`,
			want: scope{
				"m": hashMap{str("a"): number(123)},
			},
		},
		{
			name: "list",
			src:  `l := [1, 2, 3]`,
			want: scope{
				"l": list{
					elements: []Value{number(1), number(2), number(3)},
				},
			},
		},
		{
			name: "list_empty",
			src:  `l := []`,
			want: scope{
				"l": list{
					elements: []Value{},
				},
			},
		},
		{
			name: "list_index",
			src: `l := [1, 2, 3]
				  v := l[0]`,
			want: scope{
				"l": list{
					elements: []Value{number(1), number(2), number(3)},
				},
				"v": number(1),
			},
		},
		{
			name: "list_index_neg",
			src: `l := [1, 2, 3]
				  v := l[-1]`,
			want: scope{
				"l": list{
					elements: []Value{number(1), number(2), number(3)},
				},
				"v": number(3),
			},
		},
		{
			name: "list_index_range",
			src: `l := [1, 2, 3, 4]
				  s1 := l[0..2]
				  s2 := l[2..-1]
				  s3 := l[0..0]
				  s4 := l[l.count-2 .. l.count-1]`,
			want: scope{
				"l": list{
					elements: []Value{number(1), number(2), number(3), number(4)},
				},
				"s1": list{
					elements: []Value{number(1), number(2), number(3)},
				},
				"s2": list{
					elements: []Value{number(3), number(4)},
				},
				"s3": list{
					elements: []Value{number(1)},
				},
				"s4": list{
					elements: []Value{number(3), number(4)},
				},
			},
		},
		{
			name: "list_index_assign",
			src: `l := [0]
				  l[0] = 1`,
			want: scope{
				"l": list{
					elements: []Value{number(1)},
				},
			},
		},
		{
			name: "list_concat_scalars",
			src: `l := []
				  l = l :: 1 :: 2`,
			want: scope{
				"l": list{
					elements: []Value{number(1), number(2)},
				},
			},
		},
		{
			name: "list_concat_list",
			src: `l := [1, 2]
				  l = l :: [3, 4]`,
			want: scope{
				"l": list{
					elements: []Value{number(1), number(2), number(3), number(4)},
				},
			},
		},
		{
			name: "function_call",
			src: `f := fn() -> 123
			      ret := f()`,
			want: scope{
				"f": function{
					parameterNames: nil,
					body: []parser.Statement{
						parser.ReturnStmt{
							StmtBase: parser.StmtBase{},
							Result:   lang.Number(123),
						},
					},
					closure: []scope{},
				},
				"ret": number(123),
			},
		},
		{
			name: "function_call_with_param",
			src: `f := fn(x) -> x
			      ret := f(5)`,
			want: scope{
				"f": function{
					parameterNames: []string{"x"},
					body: []parser.Statement{
						parser.ReturnStmt{
							StmtBase: parser.StmtBase{},
							Result:   parser.IdentExpr("x"),
						},
					},
					closure: []scope{},
				},
				"ret": number(5),
			},
		},
		{
			name: "scopes",
			src: `x := 1
				  y := 1
				  if true {
					 x = 2
					 y := 2
				  }`,
			want: scope{
				"x": number(2),
				"y": number(1),
			},
		},
		{
			name: "compare_func",
			src: `a := compare(10, 10)
				  b := compare(10, 11)`,
			want: scope{
				"a": number(0),
				"b": number(-1),
			},
		},
		{
			name: "sort_list_fn",
			src: `ls1 := [150;10, 12;102, 200;23, 1;404]
				  ls2 := sort(ls1, fn(a, b) -> compare(a.x, b.x))`,
			want: scope{
				"ls1": list{
					elements: []Value{
						point{150, 10},
						point{12, 102},
						point{200, 23},
						point{1, 404},
					},
				},
				"ls2": list{
					elements: []Value{
						point{1, 404},
						point{12, 102},
						point{150, 10},
						point{200, 23},
					},
				},
			},
		},
		{
			name: "pipeline",
			src:  `a := 1 | $ + 1 | $ + 2 | $ + 3 | "a" + $`,
			want: scope{
				"a": str("a7"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compileAndInterpret(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("compileAndInterpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compileAndInterpret() =\n%#v\nwant\n%#v", got, tt.want)
			}
		})
	}
}

func Test_newInterpreter(t *testing.T) {
	t.Run("scopeCount", func(t *testing.T) {
		if got := newInterpreter(nil); len(got.idents) != initialScopeCount {
			t.Errorf("interpreter initial scope count = %v, want %v", len(got.idents), initialScopeCount)
		}
	})
}

func Test_validateArguments(t *testing.T) {
	numberType := reflect.TypeOf(number(0))
	boolType := reflect.TypeOf(boolean(false))
	numberSliceType := reflect.TypeOf([]number{})
	type args struct {
		arguments []Value
		params    []reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: false,
		},
		{
			name: "two_numbers_ok",
			args: args{
				arguments: []Value{number(1), number(2)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: false,
		},
		{
			name: "two_numbers_any_ok",
			args: args{
				arguments: []Value{number(1), number(2)},
				params:    []reflect.Type{numberType, valueType},
			},
			wantErr: false,
		},
		{
			name: "two_numbers_error_toofew",
			args: args{
				arguments: []Value{number(1)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: true,
		},
		{
			name: "two_numbers_error_toomany",
			args: args{
				arguments: []Value{number(1), number(2), number(3)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: true,
		},
		{
			name: "type_mismatch",
			args: args{
				arguments: []Value{number(1), boolean(false)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: true,
		},
		{
			name: "mixed_types",
			args: args{
				arguments: []Value{number(1), boolean(false)},
				params:    []reflect.Type{numberType, boolType},
			},
			wantErr: false,
		},
		{
			name: "varargs_ok",
			args: args{
				arguments: []Value{number(1), number(2), number(3)},
				params:    []reflect.Type{numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_empty_ok",
			args: args{
				arguments: []Value{},
				params:    []reflect.Type{numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_ok",
			args: args{
				arguments: []Value{boolean(true), number(1), number(2)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_single_ok",
			args: args{
				arguments: []Value{boolean(true), number(1)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_empty_ok",
			args: args{
				arguments: []Value{boolean(true)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_empty_error",
			args: args{
				arguments: []Value{},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: true,
		},
		{
			name: "varargs_trailing_type_mismatch",
			args: args{
				arguments: []Value{boolean(true), number(1), number(2), boolean(false)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateArguments(tt.args.arguments, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("validateArguments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hasMatchingType(t *testing.T) {
	type args struct {
		v   Value
		typ reflect.Type
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number_any",
			args: args{
				v:   number(1),
				typ: valueType,
			},
			want: true,
		},
		{
			name: "point_any",
			args: args{
				v:   point{},
				typ: valueType,
			},
			want: true,
		},
		{
			name: "number_number",
			args: args{
				v:   number(1),
				typ: numberType,
			},
			want: true,
		},
		{
			name: "number_point",
			args: args{
				v:   number(1),
				typ: pointType,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasMatchingType(tt.args.v, tt.args.typ); got != tt.want {
				t.Errorf("hasMatchingType() = %v, want %v", got, tt.want)
			}
		})
	}
}
