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
	program, err := parse(tokens, false)
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
			src:  "x := 1 x := 2",
			want: scope{
				"x": Number(2),
			},
		},
		{
			name: "rect",
			src:  "x := rect(1, 2, 3, 4)",
			want: scope{
				"x": Rect{Min: image.Point{1, 2}, Max: image.Point{4, 6}},
			},
		},
		{
			name: "colors",
			src:  "c1 := rgb(1,2,3) c2 := rgba(1,2,3,4) r := c1.r g := c1.g b := c1.b a := c1.a",
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
			src:  "c1 := srgb(1,2,3) c2 := srgba(1,2,3,4) scr := c1.scr scg := c1.scg scb := c1.scb sca := c1.sca",
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
			src:  "k := [1 2 3 4]",
			want: scope{
				"k": kernel{width: 2, height: 2, values: []Number{Number(1), Number(2), Number(3), Number(4)}},
			},
		},
		{
			name: "if_else",
			src:  "x := 0 if true { x = 1 } else { x = 2 }",
			want: scope{
				"x": Number(1),
			},
		},
		{
			name: "if_else_2",
			src:  "x := 0 if false { x = 1 } else { x = 2 }",
			want: scope{
				"x": Number(2),
			},
		},
		{
			name: "for",
			src:  "p := 0 for pos in rect(0,0,1,1) { p = pos }",
			want: scope{
				"p": Position{0, 0},
			},
		},
		{
			name: "for_2",
			src:  "p := 0 for pos in rect(0,0,2,2) { p = pos }",
			want: scope{
				"p": Position{1, 1},
			},
		},
		{
			name: "kernel_index",
			src:  "n := [1 2 3 4][1]",
			want: scope{
				"n": Number(2),
			},
		},
		{
			name: "kernel_index_2",
			src:  "n := [1 2 3 4][0;1]",
			want: scope{
				"n": Number(3),
			},
		},
		{
			name: "sort",
			src:  "k := sort([4 1 3 2])",
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
			src:  "min := min([4 1 3 2]) max := max([4 1 3 2])",
			want: scope{
				"min": Number(1),
				"max": Number(4),
			},
		},
		{
			name: "list",
			src:  "l := list(4, 123)",
			want: scope{
				"l": kernel{
					width:  0,
					height: 0,
					values: []Number{123, 123, 123, 123},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compileAndInterpret(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("compileAndInterpret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compileAndInterpret() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
