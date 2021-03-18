package formula

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Evaluator interface {
	Evaluate(r Resolver) (Valuer, error)
}

type Resolver interface {
	Resolve(string) interface{}
}

type Valuer interface {
	Float64() float64
	Int64() int64
}

type token interface {
	evaluate(Resolver, queue) (token, error)
	value() (Valuer, error)
}

type calculator interface {
	plus(token) (token, error)
	minus(token) (token, error)
	divide(token) (token, error)
	multiply(token) (token, error)
	invert() (token, error)
}

type queue []token

func (s *queue) pop() token {
	n := len(*s)
	n--
	e := (*s)[n]
	*s = (*s)[:n]
	return e
}

func (s *queue) push(t token) {
	*s = append(*s, t)
}

type formula []token

func (f formula) Evaluate(r Resolver) (Valuer, error) {
	q := queue{}
	for _, t := range f {
		v, err := t.evaluate(r, q)
		if err != nil {
			return nil, err
		}
		q.push(v)
	}
	return q.pop().value()
}

func New(e string) (Evaluator, error) {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(e))
	var p, q queue
	var u int
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
				p.push(function(w))
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
			p = p[:n]
			if u == 1 {
				switch t {
				case '+', '-':
					p.push(unary(t))
				default:
					return nil, fmt.Errorf("illegal unary operator %c", t)
				}
			} else {
				switch t {
				case '+', '-', '/', '*':
					p.push(binary(t))
				default:
					return nil, fmt.Errorf("illegal binary operator %c", t)
				}
			}
			u = 0
		case ',':
			n := len(p)
			for {
				if n == 0 {
					return nil, fmt.Errorf("missing comma or opening bracket")
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
		case '(':
			p.push(bracket(t))
			u = 0
		case ')':
			n := len(p)
			for {
				if n == 0 {
					return nil, fmt.Errorf("missing opening bracket")
				} else {
					n--
				}
				switch v := p[n].(type) {
				case bracket:
				case binary, unary:
					q.push(v)
					continue
				default:
					return nil, fmt.Errorf("illegal token %v", v)
				}
				break
			}
			if n > 0 {
				if v, ok := p[n-1].(function); ok {
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
					return nil, fmt.Errorf("missing closing bracket")
				default:
					return nil, fmt.Errorf("illegal token %v", v)
				}
			}
			return formula(q), nil
		default:
			return nil, fmt.Errorf("illegal token %v", scanner.TokenString(t))
		}
		u++
		t = s.Scan()
		w = s.TokenText()
	}
}
