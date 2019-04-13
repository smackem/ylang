package lang

import (
	"reflect"
	"testing"
)

func Test_number(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    scope
		wantErr bool
	}{
		{
			name: "compare_true",
			src: `eq := 1 == 1
				  neq := 1 != 2
				  gt := 10 > 5
				  ge1 := 10 >= 9
				  ge2 := 10 >= 10
				  lt := 1 < 10
				  le1 := 1 <= 1
				  le2 := 1 <= 2`,
			want: scope{
				"eq":  boolean(true),
				"neq": boolean(true),
				"gt":  boolean(true),
				"ge1": boolean(true),
				"ge2": boolean(true),
				"lt":  boolean(true),
				"le1": boolean(true),
				"le2": boolean(true),
			},
		},
		{
			name: "compare_false",
			src: `eq := 1 == 2
				  neq := 1 != 1
				  gt := 10 > 20
				  ge := 10 >= 11
				  lt := 1 < 0
				  le := 1 <= 0
				  invalid1 := 1 == "abc"
				  invalid2 := 1 > 100;200`,
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
