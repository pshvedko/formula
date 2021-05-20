package formula

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
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
	f, err := New("-1.5+Sin(2*Pi*x)/2+2+0*Rand()+0*Pow(1,2)")
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}
	j, err := UnmarshalJSON(b)
	if err != nil {
		log.Fatal(err)
	}
	v, err := j.Evaluate(Bind{"Pow''": math.Pow, "Sin'": math.Sin, "Rand": rand.Float64, "Pi": math.Pi, "x": .25})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(f)
	fmt.Println(j)
	fmt.Println(v)
	// Output:
	// [1.5 - 2 Pi * x * Sin' 2 / + 2 + 0 Rand * + 0 1 2 Pow'' * +]
	// [1.5 - 2 Pi * x * Sin' 2 / + 2 + 0 Rand * + 0 1 2 Pow'' * +]
	// 1
}
