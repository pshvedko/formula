package formula

import "fmt"

type bracket byte

func (t bracket) evaluate(Resolver, stacker) (token, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}

func (t bracket) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME") // FIXME
}
