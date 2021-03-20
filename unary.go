package formula

import "fmt"

type unary binary

func (t unary) evaluate(_ Resolver, q stacker) (token, error) {
	a, ok := q.pop()
	if !ok {
		return nil, ErrFewOperands
	}
	switch v := a.(type) {
	case calculator:
		switch t {
		case '-':
			return v.invert()
		case '+':
			return a, nil
		}
	}
	return nil, fmt.Errorf("FIXME33") // FIXME
}

func (t unary) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME34") // FIXME
}

func (t unary) GoString() string {
	return string(t)
}
