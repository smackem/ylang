package interpreter

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
				"eq":  Boolean(true),
				"neq": Boolean(true),
				"gt":  Boolean(true),
				"ge1": Boolean(true),
				"ge2": Boolean(true),
				"lt":  Boolean(true),
				"le1": Boolean(true),
				"le2": Boolean(true),
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
				"eq":       Boolean(false),
				"neq":      Boolean(false),
				"gt":       Boolean(false),
				"ge":       Boolean(false),
				"lt":       Boolean(false),
				"le":       Boolean(false),
				"invalid1": Boolean(false),
				"invalid2": Boolean(false),
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
