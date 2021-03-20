package formula

import (
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
		want    Evaluator
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Empty",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty",
			args:    args{"+"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Few",
			args:    args{"1*"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Few",
			args:    args{"(1+2/3)*"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Few",
			args:    args{"1*Pow(2,,1+Pi())"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Few",
			args:    args{"1*Pow(2,2+)"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Few",
			args:    args{"x/"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %+#v, want %+#v", got, tt.want)
			}
		})
	}
}

func ExampleNew() {
	f, err := New("-1.5+Sin(2*Pi*x)/2+2")
	if err != nil {
		log.Fatal(err)
	}
	v, err := f.Evaluate(Bind{"Sin": math.Sin, "Pi": math.Pi, "x": .25})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
	// Output:
	// 1
}
