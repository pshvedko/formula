package formula

import (
	"encoding/json"
	"fmt"
)

type binary byte

func (t binary) evaluate(_ Getter, q stacker) (token, error) {
	var a, b token
	var ok bool
	b, ok = q.pop()
	if !ok {
		return nil, ErrFewOperands
	}
	a, ok = q.pop()
	if !ok {
		return nil, ErrFewOperands
	}
	switch c := a.(type) {
	case calculator:
		switch t {
		case '+':
			return c.plus(b)
		case '-':
			return c.minus(b)
		case '*':
			return c.multiply(b)
		case '/':
			return c.divide(b)
		}
	}
	return nil, fmt.Errorf("FIXME5") // FIXME
}

func (t binary) less(a rune) bool {
	switch t {
	case '*', '/':
		return false
	}
	switch a {
	case '+', '-':
		return false
	}
	return true
}

func (t binary) String() string {
	return string(t)
}

func (t binary) MarshalJSON() ([]byte, error) {
	return json.Marshal(Token{Type: "binary", Value: t.String()})
}
