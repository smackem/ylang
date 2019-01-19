package lang

import (
	"image"
	"reflect"
	"testing"
)

func TestNumber_equals(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1==1",
			n:    1,
			args: args{other: Number(1)},
			want: Bool(true),
		},
		{
			name: "1==2",
			n:    1,
			args: args{other: Number(2)},
			want: Bool(false),
		},
		{
			name:    "1=='x'",
			n:       1,
			args:    args{other: String("x")},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.equals(tt.args.other)
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

func TestNumber_greaterThan(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "2>1",
			n:    2,
			args: args{other: Number(1)},
			want: Bool(true),
		},
		{
			name: "1>2",
			n:    1,
			args: args{other: Number(2)},
			want: Bool(false),
		},
		{
			name:    "1>'x'",
			n:       1,
			args:    args{other: String("x")},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.greaterThan(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.greaterThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.greaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_greaterThanOrEqual(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1>=1",
			n:    1,
			args: args{other: Number(1)},
			want: Bool(true),
		},
		{
			name: "2>=1",
			n:    2,
			args: args{other: Number(1)},
			want: Bool(true),
		},
		{
			name: "1>=2",
			n:    1,
			args: args{other: Number(2)},
			want: Bool(false),
		},
		{
			name:    "1>='x'",
			n:       1,
			args:    args{other: String("x")},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.greaterThanOrEqual(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.greaterThanOrEqual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.greaterThanOrEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_lessThan(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "2<3",
			n:    2,
			args: args{other: Number(3)},
			want: Bool(true),
		},
		{
			name: "2<1",
			n:    2,
			args: args{other: Number(1)},
			want: Bool(false),
		},
		{
			name:    "1<'x'",
			n:       1,
			args:    args{other: String("x")},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.lessThan(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.lessThan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.lessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_lessThanOrEqual(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		n       Number
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "2<=2",
			n:    2,
			args: args{other: Number(2)},
			want: Bool(true),
		},
		{
			name: "2<=4",
			n:    2,
			args: args{other: Number(4)},
			want: Bool(true),
		},
		{
			name: "2<=1",
			n:    2,
			args: args{other: Number(1)},
			want: Bool(false),
		},
		{
			name:    "1<='x'",
			n:       1,
			args:    args{other: String("x")},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.lessThanOrEqual(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.lessThanOrEqual() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.lessThanOrEqual() = %v, want %v", got, tt.want)
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
		n       Number
		args    args
		want    value
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
			args: args{other: Position{1, 2}},
			want: Position{121, 122},
		},
		{
			name: "120+rgb(1,2,3)",
			n:    120.0,
			args: args{other: NewRgba(1, 2, 3, 255)},
			want: NewRgba(121, 122, 123, 255),
		},
		{
			name:    "120+'x'",
			n:       120.0,
			args:    args{other: String("x")},
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
		n       Number
		args    args
		want    value
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
			args:    args{other: String("x")},
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
		n       Number
		args    args
		want    value
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
			args:    args{other: String("x")},
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
		n       Number
		args    args
		want    value
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
			args:    args{other: String("x")},
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
		n       Number
		args    args
		want    value
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
			args:    args{other: String("x")},
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
		n       Number
		want    value
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

func TestString_equals(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		s       String
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "abc==abc",
			s:    "abc",
			args: args{other: String("abc")},
			want: Bool(true),
		},
		{
			name: "abc==def",
			s:    "abc",
			args: args{other: String("def")},
			want: Bool(false),
		},
		{
			name:    "abc==1",
			s:       "abc",
			args:    args{other: Number(1)},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.equals(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("String.equals() = %v, want %v", got, tt.want)
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
		s       String
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "abc+def",
			s:    "abc",
			args: args{other: String("def")},
			want: String("abcdef"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("String.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_equals(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2==1;2",
			p:    Position{1, 2},
			args: args{other: Position{1, 2}},
			want: Bool(true),
		},
		{
			name: "1;2==2;3",
			p:    Position{1, 2},
			args: args{other: Position{2, 3}},
			want: Bool(false),
		},
		{
			name:    "1;2==100",
			p:       Position{1, 2},
			args:    args{other: Number(100)},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.equals(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Position.equals() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2+1;2",
			p:    Position{1, 2},
			args: args{other: Position{1, 2}},
			want: Position{2, 4},
		},
		{
			name: "1;2+100",
			p:    Position{1, 2},
			args: args{other: Number(100)},
			want: Position{101, 102},
		},
		{
			name:    "1;2+'x'",
			p:       Position{1, 2},
			args:    args{other: String("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.add() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2-1;2",
			p:    Position{1, 2},
			args: args{other: Position{1, 2}},
			want: Position{0, 0},
		},
		{
			name: "1;2-1",
			p:    Position{1, 2},
			args: args{other: Number(1)},
			want: Position{0, 1},
		},
		{
			name:    "1;2-'x'",
			p:       Position{1, 2},
			args:    args{other: String("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.sub(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.sub() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2*3;2",
			p:    Position{1, 2},
			args: args{other: Position{3, 2}},
			want: Position{3, 4},
		},
		{
			name: "1;2*50",
			p:    Position{1, 2},
			args: args{other: Number(50)},
			want: Position{50, 100},
		},
		{
			name:    "1;2*'x'",
			p:       Position{1, 2},
			args:    args{other: String("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.mul(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.mul() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2/3;2",
			p:    Position{10, 6},
			args: args{other: Position{5, 3}},
			want: Position{2, 2},
		},
		{
			name: "12;60/6",
			p:    Position{12, 60},
			args: args{other: Number(6)},
			want: Position{2, 10},
		},
		{
			name:    "1;2/'x'",
			p:       Position{1, 2},
			args:    args{other: String("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.div() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "5;9%4;5",
			p:    Position{5, 9},
			args: args{other: Position{4, 5}},
			want: Position{1, 4},
		},
		{
			name: "5;9%3",
			p:    Position{5, 9},
			args: args{other: Number(3)},
			want: Position{2, 0},
		},
		{
			name:    "1;2%'x'",
			p:       Position{1, 2},
			args:    args{other: String("x")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.mod(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.mod() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		{
			name: "1;2 in rect(0;0, 10, 10)",
			p:    Position{1, 2},
			args: args{other: Rect{image.Point{0, 0}, image.Point{10, 10}}},
			want: Bool(true),
		},
		{
			name: "10;12 in rect(0;0, 10, 10)",
			p:    Position{10, 12},
			args: args{other: Rect{image.Point{0, 0}, image.Point{10, 10}}},
			want: Bool(false),
		},
		{
			name:    "10;12 in rect(0;0, 10, 10)",
			p:       Position{10, 12},
			args:    args{other: Number(1)},
			want:    Bool(false),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.in(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.in() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Position.in() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPosition_neg(t *testing.T) {
	tests := []struct {
		name    string
		p       Position
		want    value
		wantErr bool
	}{
		{
			name: "-1;2",
			p:    Position{1, 2},
			want: Position{-1, -2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.neg() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: mock Bitmap interface
func TestPosition_at(t *testing.T) {
	type args struct {
		bitmap Bitmap
	}
	tests := []struct {
		name    string
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.at(tt.args.bitmap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.at() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.at() = %v, want %v", got, tt.want)
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
		p       Position
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.property(tt.args.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("Position.property() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Position.property() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_equals(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.equals(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Color.equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_add(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.add(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_sub(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.sub(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_mul(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.mul(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_div(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.div(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_mod(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.mod(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.mod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_neg(t *testing.T) {
	tests := []struct {
		name    string
		c       Color
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.neg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_property(t *testing.T) {
	type args struct {
		ident string
	}
	tests := []struct {
		name    string
		c       Color
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.property(tt.args.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("Color.property() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Color.property() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_equals(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		rect    Rect
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rect.equals(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rect.equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rect.equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_in(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		rect    Rect
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rect.in(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rect.in() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rect.in() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRect_property(t *testing.T) {
	type args struct {
		ident string
	}
	tests := []struct {
		name    string
		rect    Rect
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.rect.property(tt.args.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rect.property() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rect.property() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kernel_equals(t *testing.T) {
	type args struct {
		other value
	}
	tests := []struct {
		name    string
		k       kernel
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.equals(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("kernel.equals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("kernel.equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kernel_property(t *testing.T) {
	type args struct {
		ident string
	}
	tests := []struct {
		name    string
		k       kernel
		args    args
		want    value
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.k.property(tt.args.ident)
			if (err != nil) != tt.wantErr {
				t.Errorf("kernel.property() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("kernel.property() = %v, want %v", got, tt.want)
			}
		})
	}
}
