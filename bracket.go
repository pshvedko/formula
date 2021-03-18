package formula

import "fmt"

type bracket byte

func (t bracket) evaluate(f Resolver, q queue) (token, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t bracket) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t bracket) String() string {
	return fmt.Sprintf("bracket[%v]", byte(t))
}
