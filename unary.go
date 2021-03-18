package formula

import "fmt"

type unary binary

func (t unary) evaluate(_ Resolver, q stacker) (token, error) {
	a := q.pop()
	switch v := a.(type) {
	case calculator:
		switch t {
		case '-':
			return v.invert()
		case '+':
			return a, nil
		}
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t unary) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}
