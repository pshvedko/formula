package formula

import "fmt"

type bracket byte

func (t bracket) evaluate(Getter, stacker) (token, error) {
	return nil, fmt.Errorf("FIXME7") // FIXME
}

func (t bracket) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME8") // FIXME
}
