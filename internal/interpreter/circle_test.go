package interpreter

import (
	"reflect"
	"testing"
)

func Test_circle(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    scope
		wantErr bool
	}{
		{
			name: "compare_true",
			src: `eq := circle(10;10, 5) == circle(10;10, 5)
				  neq := circle(10;10, 5) != circle(10;10, 6)
				  ge := circle(0;0, 0) >= circle(0;0, 0)
				  le := circle(0;0, 0) <= circle(0;0, 0)`,
			want: scope{
				"eq":  boolean(true),
				"neq": boolean(true),
				"ge":  boolean(true),
				"le":  boolean(true),
			},
		},
		{
			name: "compare_false",
			src: `eq := circle(10;10, 5) == circle(11;10, 5)
				  neq := circle(10;10, 5) != circle(10;10, 5)
				  gt := circle(0;0, 0) > circle(0;0, 0)
				  ge := circle(0;0, 0) >= circle(0;0, 1)
				  lt := circle(0;0, 0) < circle(0;0, 0)
				  le := circle(0;0, 0) <= circle(0;0, 1)
				  invalid1 := circle(0;0, 0) == "abc"
				  invalid2 := circle(0;0, 0) > 100;200`,
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
