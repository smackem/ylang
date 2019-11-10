package interpreter

import (
	"github.com/smackem/ylang/internal/lang"
	"image"
	"reflect"
	"testing"
)

func TestNumber_compare(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "1 cmp 1",
			n:    1,
			args: args{other: Number(1)},
			want: Number(0),
		},
		{
			name: "1 cmp 2",
			n:    1,
			args: args{other: Number(2)},
			want: Number(-1),
		},
		{
			name: "1 cmp 'x'",
			n:    1,
			args: args{other: Str("x")},
			want: nil,
		},
		{
			name: "2 cmp 1",
			n:    2,
			args: args{other: Number(1)},
			want: Number(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Compare(tt.args.other)
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
		other Value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "120+43.5",
			n:    120.0,
			args: args{other: Number(43.5)},
			want: Number(163.5),
		},
		{
			name: "120+1;2",
			n:    120.0,
			args: args{other: Point{1, 2}},
			want: Point{121, 122},
		},
		{
			name: "120+rgb(1,2,3)",
			n:    120.0,
			args: args{other: Color(lang.NewRgba(1, 2, 3, 255))},
			want: Color(lang.NewRgba(121, 122, 123, 255)),
		},
		{
			name:    "120+'x'",
			n:       120.0,
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_sub(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "125.5-25.5",
			n:    125.5,
			args: args{other: Number(25.5)},
			want: Number(100),
		},
		{
			name:    "120-'x'",
			n:       120.0,
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Sub(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_mul(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "12*5.5",
			n:    12,
			args: args{other: Number(5.5)},
			want: Number(66),
		},
		{
			name:    "1*x",
			n:       120.0,
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Mul(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_div(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "12/6",
			n:    12,
			args: args{other: Number(6)},
			want: Number(2),
		},
		{
			name:    "1/x",
			n:       1.0,
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_mod(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "5%3",
			n:    5,
			args: args{other: Number(3)},
			want: Number(2),
		},
		{
			name:    "1%x",
			n:       1.0,
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Mod(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_neg(t *testing.T) {
	tests := []struct {
		name    string
		n       Number
		want    Value
		wantErr bool
	}{
		{
			name: "-123.125",
			n:    123.125,
			want: Number(-123.125),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Number.Neg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString_add(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		s       Str
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "abc+def",
			s:    "abc",
			args: args{other: Str("def")},
			want: Str("abcdef"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("str.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("str.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_add(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "1;2+1;2",
			p:    Point{1, 2},
			args: args{other: Point{1, 2}},
			want: Point{2, 4},
		},
		{
			name: "1;2+100",
			p:    Point{1, 2},
			args: args{other: Number(100)},
			want: Point{101, 102},
		},
		{
			name:    "1;2+'x'",
			p:       Point{1, 2},
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_sub(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "1;2-1;2",
			p:    Point{1, 2},
			args: args{other: Point{1, 2}},
			want: Point{0, 0},
		},
		{
			name: "1;2-1",
			p:    Point{1, 2},
			args: args{other: Number(1)},
			want: Point{0, 1},
		},
		{
			name:    "1;2-'x'",
			p:       Point{1, 2},
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Sub(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_mul(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "1;2*3;2",
			p:    Point{1, 2},
			args: args{other: Point{3, 2}},
			want: Point{3, 4},
		},
		{
			name: "1;2*50",
			p:    Point{1, 2},
			args: args{other: Number(50)},
			want: Point{50, 100},
		},
		{
			name:    "1;2*'x'",
			p:       Point{1, 2},
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Mul(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_div(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "1;2/3;2",
			p:    Point{10, 6},
			args: args{other: Point{5, 3}},
			want: Point{2, 2},
		},
		{
			name: "12;60/6",
			p:    Point{12, 60},
			args: args{other: Number(6)},
			want: Point{2, 10},
		},
		{
			name:    "1;2/'x'",
			p:       Point{1, 2},
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_mod(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "5;9%4;5",
			p:    Point{5, 9},
			args: args{other: Point{4, 5}},
			want: Point{1, 4},
		},
		{
			name: "5;9%3",
			p:    Point{5, 9},
			args: args{other: Number(3)},
			want: Point{2, 0},
		},
		{
			name:    "1;2%'x'",
			p:       Point{1, 2},
			args:    args{other: Str("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Mod(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_in(t *testing.T) {
	type args struct {
		other Value
	}
	tests := []struct {
		name    string
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		{
			name: "1;2 In rect(0;0, 10, 10)",
			p:    Point{1, 2},
			args: args{other: Rect{image.Point{0, 0}, image.Point{10, 10}}},
			want: Boolean(true),
		},
		{
			name: "10;12 In rect(0;0, 10, 10)",
			p:    Point{10, 12},
			args: args{other: Rect{image.Point{0, 0}, image.Point{10, 10}}},
			want: Boolean(false),
		},
		{
			name:    "10;12 In number",
			p:       Point{10, 12},
			args:    args{other: Number(1)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.In(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.In() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("point.In() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_neg(t *testing.T) {
	tests := []struct {
		name    string
		p       Point
		want    Value
		wantErr bool
	}{
		{
			name: "-1;2",
			p:    Point{1, 2},
			want: Point{-1, -2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Neg() = %v, want %v", got, tt.want)
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
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.At(tt.args.bitmap)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.At() = %v, want %v", got, tt.want)
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
		p       Point
		args    args
		want    Value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Property(tt.args.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("point.Property() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("point.Property() = %v, want %v", got, tt.want)
			}
		})
	}
}
