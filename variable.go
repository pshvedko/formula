package formula

import "fmt"

type variable string

func (t variable) evaluate(f Resolver, q queue) (token, error) {
	r := f.Resolve(string(t))
	switch v := r.(type) {
	case int:
		return decimal(v), nil
	case int8:
		return decimal(v), nil
	case int16:
		return decimal(v), nil
	case int32:
		return decimal(v), nil
	case int64:
		return decimal(v), nil
	case float32:
		return number(v), nil
	case float64:
		return number(v), nil
	default:
		return nil, fmt.Errorf("FIXME") // FIXME
	}
}

func (t variable) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}
