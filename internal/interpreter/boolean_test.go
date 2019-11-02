package interpreter

import (
	"reflect"
	"testing"
)

func Test_boolean(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		want    scope
		wantErr bool
	}{
		{
			name: "compare_true",
			src: `eq := true == true
				  neq := true != false
				  gt := true > false
				  ge1 := true >= false
				  ge2 := true >= true
				  lt := false < true
				  le1 := false <= false
				  le2 := false <= true`,
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
			src: `eq := false == true
				  neq := false != false
				  gt := false > true
				  ge := false >= true
				  lt := true < false
				  le := true <= false
				  invalid1 := true == "abc"
				  invalid2 := false > 100;200`,
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
