package formula

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name    string
		args    args
		want    Formula
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Ok",
			args: args{e: "1"},
			want: Formula{decimal(1)},
		}, {
			name: "Ok",
			args: args{e: "1+1"},
			want: Formula{decimal(1), decimal(1), binary('+')},
		}, {
			name: "Ok",
			args: args{e: "-1"},
			want: Formula{decimal(1), unary('-')},
		}, {
			name: "Ok",
			args: args{e: "-(1)"},
			want: Formula{decimal(1), unary('-')},
		}, {
			name: "Ok",
			args: args{e: "A(B())"},
			want: Formula{function("B"), function("A'")},
		}, {
			name: "Ok",
			args: args{e: "A()"},
			want: Formula{function("A")},
		}, {
			name: "Ok",
			args: args{e: "-A()"},
			want: Formula{function("A"), unary('-')},
		}, {
			name: "Ok",
			args: args{e: "-A(1,-1)"},
			want: Formula{decimal(1), decimal(1), unary('-'), function("A''"), unary('-')},
		}, {
			name: "Ok",
			args: args{e: "A(B(C()))"},
			want: Formula{function("C"), function("B'"), function("A'")},
		}, {
			name: "Ok",
			args: args{e: "A(B(),C())"},
			want: Formula{function("B"), function("C"), function("A''")},
		}, {
			name: "Ok",
			args: args{e: "A(B(),C())+D()"},
			want: Formula{function("B"), function("C"), function("A''"), function("D"), binary('+')},
		}, {
			name: "Ok",
			args: args{e: "A(B(),C(D()))"},
			want: Formula{function("B"), function("D"), function("C'"), function("A''")},
		}, {
			name:    "Empty",
			args:    args{},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Empty",
			args:    args{"+"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"1*"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"(1+2/3)*"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"F(,)"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"F(2,)"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"F((2,))"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"1*F(2,,1+F())"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"1*F(2,2+)"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Few",
			args:    args{"(1+2)+x/"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Comma",
			args:    args{","},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Comma",
			args:    args{"1,2"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Comma",
			args:    args{"(1,2)"},
			want:    nil,
			wantErr: true,
		}, {
			name:    "Comma",
			args:    args{"1+(1,2)"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() got = %v error = %v, wantErr %v", got, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNew() {
	f, err := New("1+Sin(2*Pi*x)/2")
	if err != nil {
		log.Fatal(err)
	}
	var b []byte
	b, err = json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}
	var j Formula
	err = json.Unmarshal(b, &j)
	if err != nil {
		log.Fatal(err)
	}
	var v interface{}
	v, err = j.Evaluate(Bind{"Sin'": math.Sin, "Pi": math.Pi, "x": .25})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f)
	fmt.Println(j)
	fmt.Println(v)
	// Output:
	// [1 2 Pi * x * Sin' 2 / +]
	// [1 2 Pi * x * Sin' 2 / +]
	// 1.5
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := New("1+A(y,B(3*C()/2-D(x*x)))")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFormula_Evaluate(b *testing.B) {
	formula := Formula{
		decimal(1), variable("y"),
		decimal(3), function("C"), binary(0x2a),
		decimal(2), binary(0x2f), variable("x"),
		variable("x"), binary(0x2a),
		function("D'"), binary(0x2d),
		function("B'"), function("A''"), binary(0x2b)}
	bind := Bind{
		"x":   3,
		"y":   2,
		"A''": func(int, int) int { return 1 },
		"B'":  func(float64) float64 { return 1 },
		"C":   func() int64 { return 1 },
		"D'":  func(int64) int64 { return 1 },
	}
	for i := 0; i < b.N; i++ {
		result, err := formula.Evaluate(bind)
		if err != nil {
			b.Error(err)
		}
		switch v := result.(type) {
		case int64:
			if v == 2 {
				continue
			}
		}
		b.Errorf("wrong result %v", result)
	}
}

func TestFormula_Evaluate(t *testing.T) {
	type args struct {
		r Getter
	}
	tests := []struct {
		name    string
		f       string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			f:       "36:3*(8-6)/6",
			args:    args{},
			want:    int64(4),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(t, tt.f, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() got = %v::%T, want %v::%T", got, got, tt.want, tt.want)
			}
		})
	}
}

func Calculate(t *testing.T, e string, g Getter) (interface{}, error) {
	t.Helper()
	f, err := New(e)
	if err != nil {
		return nil, err
	}
	return f.Evaluate(g)
}
