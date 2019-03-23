package lang

import (
	"image"
	"reflect"
	"testing"
)

func compileAndInterpret(src string) (scope, error) {
	tokens, err := lex(src)
	if err != nil {
		return nil, err
	}
	program, err := parse(tokens, true)
	if err != nil {
		return nil, err
	}
	ir := newInterpreter(nil)
	err = ir.visitStmtList(program.stmts)
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
				"x": Number(1),
			},
		},
		{
			name: "multi_declaration",
			src: `x := 1
			      x := 2`,
			want: scope{
				"x": Number(2),
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
				"c1": NewRgba(1, 2, 3, 255),
				"c2": NewRgba(1, 2, 3, 4),
				"r":  Number(1),
				"g":  Number(2),
				"b":  Number(3),
				"a":  Number(255),
			},
		},
		{
			name: "scolors",
			src: `c1 := srgb(1,2,3)
				  c2 := srgba(1,2,3,4)
				  scr := c1.scr
				  scg := c1.scg
				  scb := c1.scb
				  sca := c1.sca`,
			want: scope{
				"c1":  NewSrgba(1, 2, 3, 1),
				"c2":  NewSrgba(1, 2, 3, 4),
				"scr": Number(1),
				"scg": Number(2),
				"scb": Number(3),
				"sca": Number(1),
			},
		},
		{
			name: "color_literal",
			src:  "c1 := #011223:f0",
			want: scope{
				"c1": NewRgba(0x01, 0x12, 0x23, 0xf0),
			},
		},
		{
			name: "kernel",
			src:  "k := |1 2 3 4|",
			want: scope{
				"k": kernel{width: 2, height: 2, values: []Number{Number(1), Number(2), Number(3), Number(4)}},
			},
		},
		{
			name: "if_else",
			src: `x := 0
			      if true { x = 1 } else { x = 2 }`,
			want: scope{
				"x": Number(1),
			},
		},
		{
			name: "if_else_2",
			src: `x := 0
			      if false { x = 1 } else { x = 2 }`,
			want: scope{
				"x": Number(2),
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
				"n": Number(2),
			},
		},
		{
			name: "kernel_index_2",
			src:  "n := |1 2 3 4|[0;1]",
			want: scope{
				"n": Number(3),
			},
		},
		{
			name: "sort_kernel",
			src:  "k := sort_kernel(|4 1 3 2|)",
			want: scope{
				"k": kernel{
					width:  2,
					height: 2,
					values: []Number{1, 2, 3, 4},
				},
			},
		},
		{
			name: "min_max",
			src: `min := min_kernel(|4 1 3 2|)
			      max := max_kernel(|4 1 3 2|)`,
			want: scope{
				"min": Number(1),
				"max": Number(4),
			},
		},
		{
			name: "list_func",
			src:  "l := list(4, 123)",
			want: scope{
				"l": list{
					elements: []value{Number(123), Number(123), Number(123), Number(123)},
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
					values: []Number{1, 1, 1, 1, 1, 1},
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
					elements: []value{
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
				"c": Number(3),
			},
		},
		{
			name: "indexed_assign_kernel",
			src: `k := |1 2 3 4|
				  k[0] = 0`,
			want: scope{
				"k": kernel{width: 2, height: 2, values: []Number{Number(0), Number(2), Number(3), Number(4)}},
			},
		},
		{
			name: "indexed_assign_kernel_neg",
			src: `k := |1 2 3 4|
				  k[-1] = 0`,
			want: scope{
				"k": kernel{width: 2, height: 2, values: []Number{Number(1), Number(2), Number(3), Number(0)}},
			},
		},
		{
			name: "hashmap",
			src:  `m := {a: 1, b: 2, c: 3}`,
			want: scope{
				"m": hashMap{str("a"): Number(1), str("b"): Number(2), str("c"): Number(3)},
			},
		},
		{
			name: "hashmap_index",
			src: `m := {a: 1, b: 2, c: 3,}
				  a := m.a
				  b := m.b
				  c := m["c"]`,
			want: scope{
				"m": hashMap{str("a"): Number(1), str("b"): Number(2), str("c"): Number(3)},
				"a": Number(1),
				"b": Number(2),
				"c": Number(3),
			},
		},
		{
			name: "hashmap_index_2",
			src: `m := {a: 1}
				  a := m.a
				  a1 := m["a"]
				  b := m["b"]`,
			want: scope{
				"m":  hashMap{str("a"): Number(1)},
				"a":  Number(1),
				"a1": Number(1),
				"b":  nilval{},
			},
		},
		{
			name: "hashmap_indexed_assign",
			src: `m := {}
				  m["a"] = 123`,
			want: scope{
				"m": hashMap{str("a"): Number(123)},
			},
		},
		{
			name: "list",
			src:  `l := [1, 2, 3]`,
			want: scope{
				"l": list{
					elements: []value{Number(1), Number(2), Number(3)},
				},
			},
		},
		{
			name: "list_empty",
			src:  `l := []`,
			want: scope{
				"l": list{
					elements: []value{},
				},
			},
		},
		{
			name: "list_index",
			src: `l := [1, 2, 3]
				  v := l[0]`,
			want: scope{
				"l": list{
					elements: []value{Number(1), Number(2), Number(3)},
				},
				"v": Number(1),
			},
		},
		{
			name: "list_index_neg",
			src: `l := [1, 2, 3]
				  v := l[-1]`,
			want: scope{
				"l": list{
					elements: []value{Number(1), Number(2), Number(3)},
				},
				"v": Number(3),
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
					elements: []value{Number(1), Number(2), Number(3), Number(4)},
				},
				"s1": list{
					elements: []value{Number(1), Number(2), Number(3)},
				},
				"s2": list{
					elements: []value{Number(3), Number(4)},
				},
				"s3": list{
					elements: []value{Number(1)},
				},
				"s4": list{
					elements: []value{Number(3), Number(4)},
				},
			},
		},
		{
			name: "list_index_assign",
			src: `l := [0]
				  l[0] = 1`,
			want: scope{
				"l": list{
					elements: []value{Number(1)},
				},
			},
		},
		{
			name: "list_concat_scalars",
			src: `l := []
				  l = l :: 1 :: 2`,
			want: scope{
				"l": list{
					elements: []value{Number(1), Number(2)},
				},
			},
		},
		{
			name: "list_concat_list",
			src: `l := [1, 2]
				  l = l :: [3, 4]`,
			want: scope{
				"l": list{
					elements: []value{Number(1), Number(2), Number(3), Number(4)},
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
					body: []statement{
						returnStmt{
							stmtBase: stmtBase{},
							result:   Number(123),
						},
					},
					closure: []scope{},
				},
				"ret": Number(123),
			},
		},
		{
			name: "function_call_with_param",
			src: `f := fn(x) -> x
			      ret := f(5)`,
			want: scope{
				"f": function{
					parameterNames: []string{"x"},
					body: []statement{
						returnStmt{
							stmtBase: stmtBase{},
							result:   identExpr("x"),
						},
					},
					closure: []scope{},
				},
				"ret": Number(5),
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
				"x": Number(2),
				"y": Number(1),
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
	numberType := reflect.TypeOf(Number(0))
	boolType := reflect.TypeOf(boolean(false))
	numberSliceType := reflect.TypeOf([]Number{})
	type args struct {
		arguments []value
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
				arguments: []value{Number(1), Number(2)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: false,
		},
		{
			name: "two_numbers_error_toofew",
			args: args{
				arguments: []value{Number(1)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: true,
		},
		{
			name: "two_numbers_error_toomany",
			args: args{
				arguments: []value{Number(1), Number(2), Number(3)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: true,
		},
		{
			name: "type_mismatch",
			args: args{
				arguments: []value{Number(1), boolean(false)},
				params:    []reflect.Type{numberType, numberType},
			},
			wantErr: true,
		},
		{
			name: "mixed_types",
			args: args{
				arguments: []value{Number(1), boolean(false)},
				params:    []reflect.Type{numberType, boolType},
			},
			wantErr: false,
		},
		{
			name: "varargs_ok",
			args: args{
				arguments: []value{Number(1), Number(2), Number(3)},
				params:    []reflect.Type{numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_empty_ok",
			args: args{
				arguments: []value{},
				params:    []reflect.Type{numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_ok",
			args: args{
				arguments: []value{boolean(true), Number(1), Number(2)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_single_ok",
			args: args{
				arguments: []value{boolean(true), Number(1)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_empty_ok",
			args: args{
				arguments: []value{boolean(true)},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: false,
		},
		{
			name: "varargs_trailing_empty_error",
			args: args{
				arguments: []value{},
				params:    []reflect.Type{boolType, numberSliceType},
			},
			wantErr: true,
		},
		{
			name: "varargs_trailing_type_mismatch",
			args: args{
				arguments: []value{boolean(true), Number(1), Number(2), boolean(false)},
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
