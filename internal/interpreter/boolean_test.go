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
			src: `eq := false == true
				  neq := false != false
				  gt := false > true
				  ge := false >= true
				  lt := true < false
				  le := true <= false
				  invalid1 := true == "abc"
				  invalid2 := false > 100;200`,
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
