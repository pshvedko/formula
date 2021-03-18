package formula

import "fmt"

type number float64

func (t number) Plus(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t + number(v), nil
	case number:
		return t + v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) Minus(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t - number(v), nil
	case number:
		return t - v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) Divide(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t / number(v), nil
	case number:
		return t / v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) Multiply(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t * number(v), nil
	case number:
		return t * v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) Invert() (token, error) {
	return -t, nil
}

func (t number) Float64() float64 {
	return float64(t)
}

func (t number) Int64() int64 {
	return int64(t)
}

func (t number) evaluate(Resolver, queue) (token, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) value() (Valuer, error) {
	return t, nil
}

func (t number) String() string {
	return fmt.Sprintf("number[%v]", float64(t))
}
