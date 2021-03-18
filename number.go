package formula

import "fmt"

type number float64

func (t number) plus(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t + number(v), nil
	case number:
		return t + v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) minus(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t - number(v), nil
	case number:
		return t - v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) divide(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		if v == 0 {
			break
		}
		return t / number(v), nil
	case number:
		if v == 0 {
			break
		}
		return t / v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) multiply(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t * number(v), nil
	case number:
		return t * v, nil
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t number) invert() (token, error) {
	return -t, nil
}

func (t number) evaluate(Resolver, stacker) (token, error) {
	return t, nil
}

func (t number) value() (Valuer, error) {
	return t, nil
}

func (t number) Float64() float64 {
	return float64(t)
}

func (t number) Int64() int64 {
	return int64(t)
}
