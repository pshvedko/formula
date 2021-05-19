package formula

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

var (
	ErrFewOperands    = fmt.Errorf("few operands")
	ErrUnaryOperator  = fmt.Errorf("illegal unary operator")
	ErrBinaryOperator = fmt.Errorf("illegal binary operator")
	ErrUnopenedComma  = fmt.Errorf("missing comma or opening bracket")
	ErrUnopened       = fmt.Errorf("missing opening bracket")
	ErrIllegalToken   = fmt.Errorf("illegal token")
	ErrUnclosed       = fmt.Errorf("missing closing bracket")
	ErrUndefinedFunc  = fmt.Errorf("undefined function")
	ErrUndefinedVar   = fmt.Errorf("undefined variable")
	ErrDivisionByZero = fmt.Errorf("division by zero")
)

type Evaluator interface {
	Evaluate(r Getter) (Valuer, error)
}

type Getter interface {
	Get(string) (interface{}, bool)
}

type Valuer interface {
	Float64() float64
	Int64() int64
}

type stacker interface {
	pop() (token, bool)
	push(token)
}

type token interface {
	evaluate(Getter, stacker) (token, error)
	value() (Valuer, error)
}

type Token struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t Token) tokenize() token {
	switch t.Type {
	case "decimal":
		v, err := strconv.ParseInt(t.Value, 10, 64)
		if err != nil {
			return nil
		}
		return decimal(v)
	case "number":
		v, err := strconv.ParseFloat(t.Value, 64)
		if err != nil {
			return nil
		}
		return number(v)
	case "function":
		return function(t.Value)
	case "variable":
		return variable(t.Value)
	case "unary":
		switch t.Value[0] {
		case '+', '-':
			return unary(t.Value[0])
		}
	case "binary":
		switch t.Value[0] {
		case '+', '-', '*', '/':
			return binary(t.Value[0])
		}
	default:
	}
	panic(t)
	return nil
}

type calculator interface {
	plus(token) (token, error)
	minus(token) (token, error)
	divide(token) (token, error)
	multiply(token) (token, error)
	invert() (token, error)
}

type queue []token

func (q *queue) pop() (token, bool) {
	n := len(*q)
	if n == 0 {
		return nil, false
	}
	n--
	e := (*q)[n]
	*q = (*q)[:n]
	return e, true
}

func (q *queue) push(t token) {
	*q = append(*q, t)
}

type formula queue

func (f formula) validate(r Getter) (Evaluator, error) {
	var q formula
	for _, t := range f {
		switch t.(type) {
		case unary:
			t = unary('+')
		case binary:
			t = binary('*')
		case decimal:
			t = decimal(1)
		case number:
			t = number(1)
		case function:
		case variable:
		default:
			return nil, ErrIllegalToken
		}
		q = append(q, t)
	}
	if len(q) == 0 {
		return nil, ErrFewOperands
	}
	v, err := q.Evaluate(r)
	if err != nil {
		return nil, err
	}
	if v.Float64() != 1 {
		return nil, ErrIllegalToken
	}
	return f, nil
}

func (f formula) Evaluate(r Getter) (Valuer, error) {
	var q queue
	for _, t := range f {
		v, err := t.evaluate(r, &q)
		if err != nil {
			return nil, err
		}
		q.push(v)
	}
	v, ok := q.pop()
	if !ok {
		return nil, ErrFewOperands
	}
	return v.value()
}

type Bind map[string]interface{}

func (m Bind) Get(n string) (v interface{}, ok bool) {
	v, ok = m[n]
	return
}

func New(e string) (Evaluator, error) {
	m := Bind{}
	s := scanner.Scanner{}
	s.Init(strings.NewReader(e))
	var p, q queue
	var d, b []int
	var u, f int
	var t rune = scanner.Comment
	var w string
	for {
		switch t {
		case scanner.Comment:
		case scanner.Int:
			v, err := strconv.ParseInt(w, 0, 64)
			if err != nil {
				return nil, err
			}
			q.push(decimal(v))
		case scanner.Float:
			v, err := strconv.ParseFloat(w, 64)
			if err != nil {
				return nil, err
			}
			q.push(number(v))
		case scanner.Ident:
			u++
			t = s.Scan()
			if t == '(' {
				d = append(d, len(q))
				b = append(b, 0)
				p.push(function(w))
				f++
			} else {
				q.push(variable(w))
				m[w] = 1
			}
			w = s.TokenText()
			continue
		case '+', '-', '/', '*':
			n := len(p)
			for u > 1 && n > 0 {
				n--
				switch v := p[n].(type) {
				case unary:
					q.push(v)
					continue
				case binary:
					if v.less(t) {
						break
					}
					q.push(v)
					continue
				}
				n++
				break
			}
			p = p[:n]
			if u == 1 {
				switch t {
				case '+', '-':
					p.push(unary(t))
				default:
					return nil, ErrUnaryOperator
				}
			} else {
				switch t {
				case '+', '-', '/', '*':
					p.push(binary(t))
				default:
					return nil, ErrBinaryOperator
				}
			}
			u = 0
		case ',':
			n := len(p)
			for {
				if n == 0 {
					return nil, ErrUnopenedComma
				} else {
					n--
				}
				switch v := p[n].(type) {
				case bracket:
					n++
				case binary:
					q.push(v)
					continue
				}
				break
			}
			p = p[:n]
			u = 0
			b[f-1]++
		case '(':
			p.push(bracket(t))
			u = 0
		case ')':
			n := len(p)
			for {
				if n == 0 {
					return nil, ErrUnopened
				}
				n--
				switch v := p[n].(type) {
				case bracket:
				case binary, unary:
					q.push(v)
					continue
				default:
					return nil, ErrIllegalToken
				}
				break
			}
			if n > 0 {
				if v, ok := p[n-1].(function); ok {
					f--
					if d[f] < len(q) {
						b[f]++
					}
					o := make([]reflect.Type, b[f])
					for i := range o {
						o[i] = reflect.TypeOf(0)
					}
					m[string(v)] = reflect.MakeFunc(reflect.FuncOf(o, []reflect.Type{reflect.TypeOf(0)}, false),
						func(i []reflect.Value) []reflect.Value {
							return []reflect.Value{reflect.ValueOf(1)}
						}).Interface()
					d = d[:f]
					b = b[:f]
					q.push(v)
					n--
				}
			}
			p = p[:n]
		case scanner.EOF:
			n := len(p)
			for n > 0 {
				n--
				switch v := p[n].(type) {
				case binary, unary:
					q.push(v)
				case bracket:
					return nil, ErrUnclosed
				default:
					return nil, ErrIllegalToken
				}
			}
			return formula(q).validate(m)
		default:
			return nil, ErrIllegalToken
		}
		u++
		t = s.Scan()
		w = s.TokenText()
	}
}

func UnmarshalJSON(b []byte) (Evaluator, error) {
	var t []Token
	err := json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}
	var q queue
	for _, j := range t {
		v := j.tokenize()
		if v == nil {
			return nil, io.EOF
		}
		q = append(q, v)
	}
	return formula(q), nil
}
