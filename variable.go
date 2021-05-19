package formula

import (
	"encoding/json"
	"fmt"
)

type variable string

func (t variable) evaluate(f Getter, _ stacker) (token, error) {
	r, ok := f.Get(string(t))
	if !ok {
		return nil, ErrUndefinedVar
	}
	switch v := r.(type) {
	case int:
		return decimal(v), nil
	case int8:
		return decimal(v), nil
	case int16:
		return decimal(v), nil
	case int32:
		return decimal(v), nil
	case int64:
		return decimal(v), nil
	case float32:
		return number(v), nil
	case float64:
		return number(v), nil
	default:
		return nil, fmt.Errorf("FIXME44") // FIXME
	}
}

func (t variable) value() (Valuer, error) {
	return nil, fmt.Errorf("FIXME45") // FIXME
}

func (t variable) MarshalJSON() ([]byte, error) {
	return json.Marshal(Token{Type: "variable", Value: t.String()})
}

func (t variable) String() string {
	return string(t)
}
