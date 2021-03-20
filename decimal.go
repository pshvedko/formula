package formula

import "fmt"

type decimal int64

func (t decimal) plus(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) + v, nil
	case decimal:
		return t + v, nil
	}
	return nil, fmt.Errorf("FIXME9") // FIXME
}

func (t decimal) minus(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) - v, nil
	case decimal:
		return t - v, nil
	}
	return nil, fmt.Errorf("FIXME11") // FIXME
}

func (t decimal) multiply(a token) (token, error) {
	switch v := a.(type) {
	case number:
		return number(t) * v, nil
	case decimal:
		return t * v, nil
	}
	return nil, fmt.Errorf("FIXME12") // FIXME
}

func (t decimal) divide(a token) (token, error) {
	switch v := a.(type) {
	case number:
		if v == 0 {
			return nil, ErrDivisionByZero
		}
		return number(t) / v, nil
	case decimal:
		if v == 0 {
			return nil, ErrDivisionByZero
		}
		return t / v, nil
	}
	return nil, fmt.Errorf("FIXME12 %v", a) // FIXME
}

func (t decimal) invert() (token, error) {
	return -t, nil
}

func (t decimal) evaluate(Resolver, stacker) (token, error) {
	return t, nil
}

func (t decimal) value() (Valuer, error) {
	return t, nil
}

func (t decimal) Float64() float64 {
	return float64(t)
}

func (t decimal) Int64() int64 {
	return int64(t)
}
