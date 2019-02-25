package lang

import "testing"

func Test_hashMap_indexAssign(t *testing.T) {
	type args struct {
		key value
		val value
	}
	tests := []struct {
		name string
		h    hashMap
		args args
	}{
		{
			name: "hashMap_set",
			h:    make(hashMap),
			args: args{
				key: str("key"),
				val: Number(100),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.indexAssign(tt.args.key, tt.args.val)
			if got := tt.h[tt.args.key]; got != tt.args.val {
				t.Errorf("hashMap[%v] = %v, want %v", tt.args.key, got, tt.args.val)
			}
		})
	}
}
