package formula

import (
	"fmt"
	"reflect"
)

type function string

func (t function) evaluate(e Resolver, q queue) (token, error) {
	f := e.Resolve(string(t))
	r := reflect.ValueOf(f)
	j := r.Type()
	if r.Kind() == reflect.Func && !j.IsVariadic() {
		if j.NumOut() == 1 {
			i := j.NumIn()
			p := make([]reflect.Value, i)
			for i > 0 {
				i--
				a := j.In(i)
				switch a.Kind() {
				case reflect.Float32, reflect.Float64:
					p[i] = reflect.ValueOf(q.pop().(Valuer).Float64()).Convert(a)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					p[i] = reflect.ValueOf(q.pop().(Valuer).Int64()).Convert(a)
				default:
					return nil, fmt.Errorf("FIXME") // FIXME
				}
			}
			v := r.Call(p)[0]
			switch v.Kind() {
			case reflect.Float32, reflect.Float64:
				return number(v.Float()), nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return decimal(v.Int()), nil
			default:
				return nil, fmt.Errorf("FIXME") // FIXME
			}
		}
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t function) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t function) String() string {
	return fmt.Sprintf("function[%v]", string(t))
}
