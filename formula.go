package formula

import (
	"encoding/json"
	"fmt"
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

type Getter interface {
	Get(string) (interface{}, bool)
}

type stacker interface {
	pop() (token, bool)
	push(token)
}

type token interface {
	evaluate(Getter, stacker) (token, error)
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
	defer q.len(n)
	return (*q)[n], true
}

func (q *queue) push(t token) {
	*q = append(*q, t)
}

func (q *queue) len(n int) {
	*q = (*q)[:n]
}

type Formula queue

func (q queue) validate() bool {
	var p int
	for _, t := range q {
		switch v := t.(type) {
		case decimal, number, variable:
			p++
		case binary:
			p--
		case unary:
		case function:
			p -= strings.Count(string(v), "'")
			p++
		default:
			return false
		}
	}
	return p == 1
}

// Evaluate calculating the value of an expression using the
// Dijkstra Shunting Yard algorithm
func (f Formula) Evaluate(r Getter) (interface{}, error) {
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
	switch v.(type) {
	case decimal, number:
		return v, nil
	}
	return nil, ErrIllegalToken
}

type Variable map[string]interface{}

func (m Variable) Get(n string) (v interface{}, ok bool) {
	v, ok = m[n]
	return
}

func New(e string) (Formula, error) {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(e))
	var p, q queue
	var u, f int
	var a []int
	var t rune
	var w string
	for {
		switch t {
		case 0:
			a = append(a, 0)
		case scanner.Comment:
			u--
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
				p.push(function(w))
				f++
				a = append(a, 0)
			} else {
				q.push(variable(w))
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
			p.len(n)
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
			p.len(n)
			if a[f] == 0 {
				a[f]++
			}
			a[f]++
			u = 0
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
					if a[f] == 0 && u > 1 {
						a[f]++
					}
					v += function(strings.Repeat("'", a[f]))
					q.push(v)
					n--
					a = a[:f]
					f--
				}
			}
			p.len(n)
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
			if q.validate() {
				return Formula(q), nil
			}
			return nil, ErrFewOperands
		default:
			return nil, ErrIllegalToken
		}
		u++
		t = s.Scan()
		w = s.TokenText()
	}
}

type Token struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

func (t Token) tokenize() (token, error) {
	switch t.Type {
	case "decimal":
		v, err := strconv.ParseInt(t.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return decimal(v), nil
	case "number":
		v, err := strconv.ParseFloat(t.Value, 64)
		if err != nil {
			return nil, err
		}
		return number(v), nil
	case "function":
		return function(t.Value), nil
	case "variable":
		return variable(t.Value), nil
	case "unary":
		switch t.Value[0] {
		case '+', '-':
			return unary(t.Value[0]), nil
		}
	case "binary":
		switch t.Value[0] {
		case '+', '-', '*', '/':
			return binary(t.Value[0]), nil
		}
	}
	return nil, ErrIllegalToken
}

func (f *Formula) UnmarshalJSON(b []byte) error {
	var t []Token
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	var q queue
	for _, j := range t {
		v, err := j.tokenize()
		if v == nil {
			return err
		}
		q = append(q, v)
	}
	if q != nil && q.validate() {
		*f = Formula(q)
		return nil
	}
	return ErrFewOperands
}
