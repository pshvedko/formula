package formula

import "fmt"

type calculator interface {
	Plus(token) (token, error)
	Minus(token) (token, error)
	Divide(token) (token, error)
	Multiply(token) (token, error)
	Invert() (token, error)
}

type binary rune

func (t binary) evaluate(_ Resolver, q queue) (token, error) {
	b := q.pop()
	a := q.pop()
	switch c := a.(type) {
	case calculator:
		switch t {
		case '+':
			return c.Plus(b)
		case '-':
			return c.Minus(b)
		case '*':
			return c.Multiply(b)
		case '/':
			return c.Divide(b)
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

func (t binary) String() string {
	return fmt.Sprintf("binary[%c]", rune(t))
}
