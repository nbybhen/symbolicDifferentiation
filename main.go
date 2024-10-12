package main

import "fmt"

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
	_ Node = (*Variable)(nil)

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

type Variable struct {
	name string
}

func (v Variable) Diff() Node {
	return Literal{1}
}

func (v Variable) Simplify() Node {
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

func main() {
	//test := Add{Literal{3}, Literal{2}}
	test := Add{Variable{"x"}, Literal{2}}
	fmt.Printf("Test: %#v\n", test)

	tmp := test.Diff()
	fmt.Printf("Post-diff: %#v\n", tmp)

	tmp = tmp.Simplify()
	fmt.Printf("Post-simplify: %#v\n", tmp)
}
