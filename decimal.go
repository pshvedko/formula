package formula

import "fmt"

type decimal int64

func (t decimal) Plus(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) + v, nil
	case decimal:
		return t + v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t decimal) Minus(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) - v, nil
	case decimal:
		return t - v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t decimal) Multiply(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) * v, nil
	case decimal:
		return t * v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t decimal) Divide(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) / v, nil
	case decimal:
		return t / v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t decimal) Invert() (token, error) {
	return -t, nil
}

func (t decimal) Float64() float64 {
	return float64(t)
}

func (t decimal) Int64() int64 {
	return int64(t)
}

func (t decimal) evaluate(Resolver, queue) (token, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t decimal) value() (Valuer, error) {
	return t, nil
}

func (t decimal) String() string {
	return fmt.Sprintf("int[%v]", int64(t))
}
