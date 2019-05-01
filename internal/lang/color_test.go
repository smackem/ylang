package lang

import (
	"reflect"
	"testing"
)

func Test_color(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    scope
		wantErr bool
	}{
		{
			name: "compare_true",
			src: `eq1 := #ffcc44:80 == #ffcc44:80
				  eq2 := #ffcc44 == #ffcc44:ff
				  neq := #ffcc44:80 != #ffcc33
				  ge := #ffcc44 >= #ffcc44
				  le := #ffcc44 <= #ffcc44`,
			want: scope{
				"eq1": boolean(true),
				"eq2": boolean(true),
				"neq": boolean(true),
				"ge":  boolean(true),
				"le":  boolean(true),
			},
		},
		{
			name: "compare_false",
			src: `eq := #c0a010 == #c1a111
				  neq := #c0a010 != #c0a010
				  gt := #c0a010 > #c0a000
				  ge := #c0a010 >= #c0a000
				  lt := #c0a010 < #c0a000
				  le := #c0a010 <= #c0a000
				  invalid1 := #c0a010 == "abc"
				  invalid2 := #c0a010 > 100;200`,
			want: scope{
				"eq":       boolean(false),
				"neq":      boolean(false),
				"gt":       boolean(false),
				"ge":       boolean(false),
				"lt":       boolean(false),
				"le":       boolean(false),
				"invalid1": boolean(false),
				"invalid2": boolean(false),
			},
		},
		{
			name: "ctors",
			src: `c1 := rgb(255, 0, 128)
				  c2 := rgb01(0.5, 1.0, 0.0)
				  c3 := rgba(255, 0, 128, 64)
				  c4 := rgba01(0.5, 1.0, 0.0, 0.25)`,
			want: scope{
				"c1": NewRgba(255, 0, 128, 255),
				"c2": NewSrgba(0.5, 1.0, 0.0, 1.0),
				"c3": NewRgba(255, 0, 128, 64),
				"c4": NewSrgba(0.5, 1.0, 0.0, 0.25),
			},
		},
		{
			name: "add",
			src: `c1 := rgb(1, 2, 3) + rgb(1, 2, 3)
				  c2 := rgba(1, 2, 3, 4) + rgba(1, 2, 3, 4)
				  c3 := rgb(1, 2, 3) + 100
				  c4 := 100 + rgb(1, 2, 3)`,
			want: scope{
				"c1": NewRgba(2, 4, 6, 255),
				"c2": NewRgba(2, 4, 6, 4),
				"c3": NewRgba(101, 102, 103, 255),
				"c4": NewRgba(101, 102, 103, 255),
			},
		},
		{
			name: "sub",
			src: `c1 := rgb(1, 2, 3) - rgb(1, 2, 3)
				  c2 := rgba(1, 2, 3, 4) - rgba(1, 2, 3, 4)
				  c3 := rgb(101, 102, 103) - 100
				  c4 := 100 - rgb(1, 2, 3)`,
			want: scope{
				"c1": NewRgba(0, 0, 0, 255),
				"c2": NewRgba(0, 0, 0, 4),
				"c3": NewRgba(1, 2, 3, 255),
				"c4": NewRgba(99, 98, 97, 255),
			},
		},
		{
			name: "mul",
			src: `c1 := rgb01(0.5, 0.5, 0.5) * rgb(10, 20, 30)
				  c2 := rgba01(0.5, 0.5, 0.5, 0.5) * rgb(10, 20, 30)
				  c3 := rgb(1, 2, 3) * 4
				  c4 := 4 * rgb(1, 2, 3)`,
			want: scope{
				"c1": NewRgba(5, 10, 15, 255),
				"c2": NewRgba(5, 10, 15, 127.5),
				"c3": NewRgba(4, 8, 12, 255),
				"c4": NewRgba(4, 8, 12, 255),
			},
		},
		{
			name: "div",
			src: `c1 := rgb01(1, 1, 1) / rgb(127.5, 127.5, 127.5)
				  c2 := rgb(20, 40, 80) / 4
				  c3 := 60 / rgb(2, 4, 6)`,
			want: scope{
				"c1": NewSrgba(2, 2, 2, 1),
				"c2": NewRgba(5, 10, 20, 255),
				"c3": NewRgba(30, 15, 10, 255),
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
