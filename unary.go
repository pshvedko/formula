package formula

import "fmt"

type unary binary

func (t unary) evaluate(f Resolver, q queue) (token, error) {
	a := q.pop()
	switch v := a.(type) {
	case calculator:
		switch t {
		case '-':
			return v.Invert()
		case '+':
			return a, nil
		}
	}
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t unary) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t unary) String() string {
	return fmt.Sprintf("unary[%c]", rune(t))
}
