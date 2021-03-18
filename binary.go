package formula

import "fmt"

type binary byte

func (t binary) evaluate(_ Resolver, q stacker) (token, error) {
	b := q.pop()
	a := q.pop()
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
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t binary) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
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
