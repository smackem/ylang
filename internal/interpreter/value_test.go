package interpreter

import (
	"github.com/smackem/ylang/internal/lang"
	"image"
	"reflect"
	"testing"
)

func TestNumber_compare(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1 cmp 1",
			n:    1,
			args: args{other: number(1)},
			want: number(0),
		},
		{
			name: "1 cmp 2",
			n:    1,
			args: args{other: number(2)},
			want: number(-1),
		},
		{
			name: "1 cmp 'x'",
			n:    1,
			args: args{other: str("x")},
			want: nil,
		},
		{
			name: "2 cmp 1",
			n:    2,
			args: args{other: number(1)},
			want: number(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.compare(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_add(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "120+43.5",
			n:    120.0,
			args: args{other: number(43.5)},
			want: number(163.5),
		},
		{
			name: "120+1;2",
			n:    120.0,
			args: args{other: point{1, 2}},
			want: point{121, 122},
		},
		{
			name: "120+rgb(1,2,3)",
			n:    120.0,
			args: args{other: color(lang.NewRgba(1, 2, 3, 255))},
			want: color(lang.NewRgba(121, 122, 123, 255)),
		},
		{
			name:    "120+'x'",
			n:       120.0,
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_sub(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "125.5-25.5",
			n:    125.5,
			args: args{other: number(25.5)},
			want: number(100),
		},
		{
			name:    "120-'x'",
			n:       120.0,
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.sub(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_mul(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "12*5.5",
			n:    12,
			args: args{other: number(5.5)},
			want: number(66),
		},
		{
			name:    "1*x",
			n:       120.0,
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.mul(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_div(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "12/6",
			n:    12,
			args: args{other: number(6)},
			want: number(2),
		},
		{
			name:    "1/x",
			n:       1.0,
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_mod(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "5%3",
			n:    5,
			args: args{other: number(3)},
			want: number(2),
		},
		{
			name:    "1%x",
			n:       1.0,
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.mod(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_neg(t *testing.T) {
	tests := []struct {
		name    string
		n       number
		want    value
		wantErr bool
	}{
		{
			name: "-123.125",
			n:    123.125,
			want: number(-123.125),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.neg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_add(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		s       str
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "abc+def",
			s:    "abc",
			args: args{other: str("def")},
			want: str("abcdef"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("str.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("str.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_add(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2+1;2",
			p:    point{1, 2},
			args: args{other: point{1, 2}},
			want: point{2, 4},
		},
		{
			name: "1;2+100",
			p:    point{1, 2},
			args: args{other: number(100)},
			want: point{101, 102},
		},
		{
			name:    "1;2+'x'",
			p:       point{1, 2},
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_sub(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2-1;2",
			p:    point{1, 2},
			args: args{other: point{1, 2}},
			want: point{0, 0},
		},
		{
			name: "1;2-1",
			p:    point{1, 2},
			args: args{other: number(1)},
			want: point{0, 1},
		},
		{
			name:    "1;2-'x'",
			p:       point{1, 2},
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.sub(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_mul(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2*3;2",
			p:    point{1, 2},
			args: args{other: point{3, 2}},
			want: point{3, 4},
		},
		{
			name: "1;2*50",
			p:    point{1, 2},
			args: args{other: number(50)},
			want: point{50, 100},
		},
		{
			name:    "1;2*'x'",
			p:       point{1, 2},
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.mul(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_div(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2/3;2",
			p:    point{10, 6},
			args: args{other: point{5, 3}},
			want: point{2, 2},
		},
		{
			name: "12;60/6",
			p:    point{12, 60},
			args: args{other: number(6)},
			want: point{2, 10},
		},
		{
			name:    "1;2/'x'",
			p:       point{1, 2},
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_mod(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "5;9%4;5",
			p:    point{5, 9},
			args: args{other: point{4, 5}},
			want: point{1, 4},
		},
		{
			name: "5;9%3",
			p:    point{5, 9},
			args: args{other: number(3)},
			want: point{2, 0},
		},
		{
			name:    "1;2%'x'",
			p:       point{1, 2},
			args:    args{other: str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.mod(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_in(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2 in rect(0;0, 10, 10)",
			p:    point{1, 2},
			args: args{other: rect{image.Point{0, 0}, image.Point{10, 10}}},
			want: boolean(true),
		},
		{
			name: "10;12 in rect(0;0, 10, 10)",
			p:    point{10, 12},
			args: args{other: rect{image.Point{0, 0}, image.Point{10, 10}}},
			want: boolean(false),
		},
		{
			name:    "10;12 in number",
			p:       point{10, 12},
			args:    args{other: number(1)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.in(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.in() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("point.in() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_neg(t *testing.T) {
	tests := []struct {
		name    string
		p       point
		want    value
		wantErr bool
	}{
		{
			name: "-1;2",
			p:    point{1, 2},
			want: point{-1, -2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("point.neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.neg() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: mock Bitmap interface
func TestPosition_at(t *testing.T) {
	type args struct {
		bitmap BitmapContext
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		// TODO: add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.at(tt.args.bitmap)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.at() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.at() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_property(t *testing.T) {
	type args struct {
		ident string
	}
	tests := []struct {
		name    string
		p       point
		args    args
		want    value
		wantErr bool
	}{
		// TODO: add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.property(tt.args.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.property() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.property() = %v, want %v", got, tt.want)
			}
		})
	}
}
