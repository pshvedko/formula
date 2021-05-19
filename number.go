package formula

import (
	"encoding/json"
	"fmt"
)

type number float64

func (t number) plus(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t + number(v), nil
	case number:
		return t + v, nil
	}
	return nil, fmt.Errorf("FIXME22") // FIXME
}

func (t number) minus(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t - number(v), nil
	case number:
		return t - v, nil
	}
	return nil, fmt.Errorf("FIXME23") // FIXME
}

func (t number) divide(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		if v == 0 {
			return nil, ErrDivisionByZero
		}
		return t / number(v), nil
	case number:
		if v == 0 {
			return nil, ErrDivisionByZero
		}
		return t / v, nil
	}
	return nil, fmt.Errorf("FIXME24") // FIXME
}

func (t number) multiply(a token) (token, error) {
	switch v := a.(type) {
	case decimal:
		return t * number(v), nil
	case number:
		return t * v, nil
	}
	return nil, fmt.Errorf("FIXME25") // FIXME
}

func (t number) invert() (token, error) {
	return -t, nil
}

func (t number) evaluate(Getter, stacker) (token, error) {
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

func (t number) MarshalJSON() ([]byte, error) {
	return json.Marshal(Token{Type: "number", Value: t.String()})
}

func (t number) String() string {
	return fmt.Sprint(t.Float64())
}
