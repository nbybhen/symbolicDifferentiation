package main

import (
	"fmt"
)

// Ensures that interfaces were implemented properly
var (
	_ Node = (*Add)(nil)
	_ Node = (*Subtract)(nil)
	_ Node = (*Multiply)(nil)
	_ Node = (*Divide)(nil)
	_ Node = (*Exponent)(nil)
	_ Node = (*Log)(nil)
	_ Node = (*Exp)(nil)

	_ Node = (*Literal)(nil)
	_ Node = (*Identifier)(nil)

	_ Node = (*Sin)(nil)
	_ Node = (*Cos)(nil)
	_ Node = (*Tan)(nil)
)

type Node interface {
	Diff() Node
	Simplify() Node
}

// Basic operations (+ * - / ^ ln)

type Add struct {
	left  Node
	right Node
}

func (a Add) Diff() Node {
	return Add{a.left.Diff(), a.right.Diff()}
}

func (a Add) Simplify() Node {
	switch l := a.left.Simplify().(type) {
	case Literal:
		switch r := a.right.Simplify().(type) {
		case Literal:
			return Literal{l.value + r.value}
		default:
			if l.value == 0 {
				return r
			}
		}
	default:
		switch r := a.right.Simplify().(type) {
		case Literal:
			if r.value == 0 {
				return l
			}
		}
	}
	return a
}

type Subtract struct {
	left  Node
	right Node
}

func (s Subtract) Diff() Node {
	return Subtract{s.left.Diff(), s.right.Diff()}
}

func (s Subtract) Simplify() Node {
	switch l := s.left.Simplify().(type) {
	case Literal:
		switch r := s.right.Simplify().(type) {
		case Literal:
			return Literal{l.value - r.value}
		default:
			if l.value == 0 {
				return r
			}
		}
	default:
		switch r := s.right.Simplify().(type) {
		case Literal:
			if r.value == 0 {
				return l
			}
		}
	}
	return s
}

type Multiply struct {
	left  Node
	right Node
}

func (m Multiply) Diff() Node {
	return Multiply{m.left.Diff(), m.right.Diff()}
}

func (m Multiply) Simplify() Node {
	switch l := m.left.Simplify().(type) {
	case Literal:
		switch r := m.right.Simplify().(type) {
		case Literal:
			return Literal{l.value * r.value}
		default:
			if l.value == 0 {
				return Literal{0}
			} else if l.value == 1 {
				return r
			}
		}
	default:
		switch r := m.right.Simplify().(type) {
		case Literal:
			if r.value == 0 {
				return Literal{0}
			} else if r.value == 1 {
				return l
			}
		}
	}
	return m
}

type Divide struct {
	left  Node
	right Node
}

func (d Divide) Simplify() Node {
	switch l := d.left.Simplify().(type) {
	case Literal:
		switch r := d.right.Simplify().(type) {
		case Literal:
			return Literal{l.value / r.value}
		default:
			if l.value == 0 {
				return Literal{0}
			}
		}
	default:
		switch r := d.right.Simplify().(type) {
		case Literal:
			if r.value == 0 {
				panic("Division by zero.")
			} else if r.value == 1 {
				return l
			}
		}
	}
	return d
}

func (d Divide) Diff() Node {
	return Divide{d.left.Diff(), d.right.Diff()}
}

type Log struct {
	inner Node
}

func (l Log) Simplify() Node {
	return Log{l.inner.Simplify()}
}

type Exp struct {
	exp Node
}

func (e Exp) Diff() Node {
	return Multiply{e, e.exp.Diff()}
}

func (e Exp) Simplify() Node {
	return Exp{e.exp.Simplify()}
}

// d/dx(ln(x)) = 1/ln(x) * x'

func (l Log) Diff() Node {
	return Multiply{Divide{Literal{1}, Log{l.inner}}, l.inner.Diff()}
}

type Exponent struct {
	base Node
	exp  Node
}

func (e Exponent) Simplify() Node {
	switch exp := e.exp.(type) {
	case Literal:
		if exp.value == 0 {
			return Literal{1}
		} else if exp.value == 1 {
			return e.base
		}
	}
	return e
}

// d/dx(base^exp) = base^(exp - 1) * (exp * d/dx(base) + base * ln(base) * d/dx(exp))

func (e Exponent) Diff() Node {
	return Multiply{Exponent{e.base, Subtract{e.exp, Literal{1}}}, Add{Multiply{e.exp, e.base.Diff()}, Multiply{Multiply{e.base, Log{e.base}}, e.exp.Diff()}}}
}

// Container types

type Literal struct {
	value int
}

func (l Literal) Diff() Node {
	return Literal{0}
}

func (l Literal) Simplify() Node {
	return l
}

type Identifier struct {
	name string
}

func (v Identifier) Diff() Node {
	return Literal{1}
}

func (v Identifier) Simplify() Node {
	return v
}

// Trig functions

type Sin struct {
	inner Node
}

func (s Sin) Diff() Node {
	return Multiply{Cos{s.inner}, s.inner.Diff()}
}

func (s Sin) Simplify() Node {
	return Sin{s.inner.Simplify()}
}

type Cos struct {
	inner Node
}

func (c Cos) Diff() Node {
	return Multiply{Multiply{Literal{-1}, Sin{c.inner}}, c.inner.Diff()}
}

func (c Cos) Simplify() Node {
	return Cos{c.inner.Simplify()}
}

type Tan struct {
	inner Node
}

func (t Tan) Diff() Node {
	return Divide{t.inner.Diff(), Exponent{Cos{t.inner}, Literal{2}}}
}

func (t Tan) Simplify() Node {
	return Tan{t.inner.Simplify()}
}

type Operator struct {
	value string
}

func (o Operator) Diff() Node {
	return o
}

func (o Operator) Simplify() Node {
	return o
}

type Stack struct {
	data []Node
}

func (s *Stack) Push(item Node) {
	s.data = append(s.data, item)
}

func (s *Stack) Pop() Node {
	ret := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return ret
}

func (s *Stack) Clear() {
	s.data = nil
}

type LeftParen struct {
	value string
}

func (l LeftParen) Diff() Node {
	return l
}

func (l LeftParen) Simplify() Node {
	return l
}

type RightParen struct {
	value string
}

func (r RightParen) Diff() Node {
	return r
}

func (r RightParen) Simplify() Node {
	return r
}

func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {
	stack := Stack{}
	tests := [4]string{"(+ 2 x)", "(* (+ x 3) 5)", "(cos (+ x 1))", "(^ x 2)"}
	operations := map[string]Node{"+": Add{}, "-": Subtract{}, "*": Multiply{}, "/": Divide{}}

	for _, exp := range tests {
		for i := len(exp) - 1; i >= 0; i-- {
			if exp[i] > 'A' && exp[i] < 'z' {
				str := ""
				// Collects string
				for exp[i] > 'A' && exp[i] < 'z' {
					str += string(exp[i])
					i--
				}
				stack.Push(Identifier{Reverse(str)})
			} else if val, ok := operations[string(exp[i])]; ok {
				r := stack.Pop()
				l := stack.Pop()
				switch val.(type) {
				case Add:
					stack.Push(Add{l, r})
				case Multiply:
					stack.Push(Multiply{l, r})
				case Divide:
					stack.Push(Divide{l, r})
				case Subtract:
					stack.Push(Subtract{l, r})
				}
			} else if string(exp[i]) == "(" {
				//stack.Push(LeftParen{"("})
			} else if string(exp[i]) == ")" {
				//stack.Push(RightParen{")"})
			} else if exp[i] > '0' && exp[i] < '9' {
				stack.Push(Literal{int(exp[i] - '0')})
			}
		}
		fmt.Printf("From %s to: %#v\n", exp, stack)
		stack = Stack{}
	}
}
