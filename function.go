package formula

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type function string

func (t function) evaluate(e Getter, q stacker) (token, error) {
	f, ok := e.Get(string(t))
	if !ok {
		return nil, ErrUndefinedFunc
	}
	r := reflect.ValueOf(f)
	j := r.Type()
	if r.Kind() == reflect.Func && !j.IsVariadic() && j.NumOut() == 1 {
		i := j.NumIn()
		p := make([]reflect.Value, i)
		var v token
		for i > 0 {
			i--
			a := j.In(i)
			v, ok = q.pop()
			if !ok {
				return nil, ErrFewOperands
			}
			switch b := v.(type) {
			case decimal, number:
				switch a.Kind() {
				case reflect.Float32, reflect.Float64:
					p[i] = reflect.ValueOf(b).Convert(a)
					continue
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					p[i] = reflect.ValueOf(b).Convert(a)
					continue
				}
			}
			return nil, fmt.Errorf("FIXME3") // FIXME
		}
		r = r.Call(p)[0]
		switch r.Kind() {
		case reflect.Float32, reflect.Float64:
			return number(r.Float()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return decimal(r.Int()), nil
		}
	}
	return nil, fmt.Errorf("FIXME2 <%v> <%v> <%v> ", r, r.Kind() == reflect.Func, r.Type()) // FIXME
}

func (t function) MarshalJSON() ([]byte, error) {
	return json.Marshal(Token{Type: "function", Value: t.String()})
}

func (t function) String() string {
	return string(t)
}
