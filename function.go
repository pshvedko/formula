package formula

import (
	"fmt"
	"reflect"
)

type function string

func (t function) evaluate(e Resolver, q stacker) (token, error) {
	f, ok := e.Resolve(string(t))
	if !ok {
		return nil, fmt.Errorf("FIXME") // FIXME
	}
	r := reflect.ValueOf(f)
	j := r.Type()
	if r.Kind() == reflect.Func && !j.IsVariadic() && j.NumOut() == 1 {
		i := j.NumIn()
		p := make([]reflect.Value, i)
		for i > 0 {
			i--
			a := j.In(i)
			v := q.pop()
			switch b := v.(type) {
			case Valuer:
				switch a.Kind() {
				case reflect.Float32, reflect.Float64:
					p[i] = reflect.ValueOf(b.Float64()).Convert(a)
					continue
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					p[i] = reflect.ValueOf(b.Int64()).Convert(a)
					continue
				}
			}
			return nil, fmt.Errorf("FIXME") // FIXME
		}
		v := r.Call(p)[0]
		switch v.Kind() {
		case reflect.Float32, reflect.Float64:
			return number(v.Float()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return decimal(v.Int()), nil
		}
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t function) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}
