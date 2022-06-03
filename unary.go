package formula

import (
	"encoding/json"
	"fmt"
)

type unary binary

func (t unary) evaluate(_ Getter, q stacker) (token, error) {
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

func (t unary) String() string {
	return string(t)
}

func (t unary) MarshalJSON() ([]byte, error) {
	return json.Marshal(Token{Type: "unary", Value: t.String()})
}
