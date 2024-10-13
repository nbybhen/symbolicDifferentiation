package main

import (
	"fmt"
	"math"
)

// Ensures that interfaces were implemented properly
var (
	_ Node = (*Add)(nil)
	_ Node = (*Subtract)(nil)
	_ Node = (*Multiply)(nil)
	_ Node = (*Divide)(nil)
	_ Node = (*Pow)(nil)
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
	PrefixString() string
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
	l, r := a.left.Simplify(), a.right.Simplify()
	switch l := l.(type) {
	case Literal:
		switch r := r.(type) {
		case Literal:
			return Literal{l.value + r.value}
		default:
			if l.value == 0 {
				return r
			}
		}
	default:
		switch r := r.(type) {
		case Literal:
			if r.value == 0 {
				return l
			}
		}
	}
	return Add{l, r}
}

func (a Add) PrefixString() string {
	return fmt.Sprintf("(+ %s %s)", a.left.PrefixString(), a.right.PrefixString())
}

type Subtract struct {
	left  Node
	right Node
}

func (s Subtract) Diff() Node {
	return Subtract{s.left.Diff(), s.right.Diff()}
}

func (s Subtract) Simplify() Node {
	l, r := s.left.Simplify(), s.right.Simplify()
	switch l := l.(type) {
	case Literal:
		switch r := r.(type) {
		case Literal:
			return Literal{l.value - r.value}
		default:
			if l.value == 0 {
				return r
			}
		}
	default:
		switch r := r.(type) {
		case Literal:
			if r.value == 0 {
				return l
			}
		}
	}
	return Subtract{l, r}
}

func (s Subtract) PrefixString() string {
	return fmt.Sprintf("(- %s %s)", s.left.PrefixString(), s.right.PrefixString())
}

type Multiply struct {
	left  Node
	right Node
}

func (m Multiply) Diff() Node {
	return Add{Multiply{m.left.Diff(), m.right}, Multiply{m.left, m.right.Diff()}}
}

func (m Multiply) Simplify() Node {
	l, r := m.left.Simplify(), m.right.Simplify()
	switch l := l.(type) {
	case Literal:
		switch r := r.(type) {
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
		switch r := r.(type) {
		case Literal:
			if r.value == 0 {
				return Literal{0}
			} else if r.value == 1 {
				return l
			}
		}
	}
	return Multiply{l, r}
}

func (m Multiply) PrefixString() string {
	return fmt.Sprintf("(* %s %s)", m.left.PrefixString(), m.right.PrefixString())
}

type Divide struct {
	left  Node
	right Node
}

func (d Divide) Diff() Node {
	fmt.Println("Left: ", d.left, "Right: ", d.right)
	return Divide{Subtract{Multiply{d.left.Diff(), d.right}, Multiply{d.left, d.right.Diff()}}, Pow{d.right, Literal{2}}}
}

func (d Divide) Simplify() Node {
	l, r := d.left.Simplify(), d.right.Simplify()
	switch l := l.(type) {
	case Literal:
		switch r := r.(type) {
		case Literal:
			return Literal{l.value / r.value}
		default:
			if l.value == 0 {
				return Literal{0}
			}
		}
	default:
		switch r := r.(type) {
		case Literal:
			if r.value == 0 {
				panic("Division by zero.")
			} else if r.value == 1 {
				return l
			}
		}
	}
	return Divide{l, r}
}

func (d Divide) PrefixString() string {
	return fmt.Sprintf("(/ %s %s)", d.left.PrefixString(), d.right.PrefixString())
}

type Log struct {
	inner Node
}

// d/dx(ln(x)) = 1/ln(x) * x'

func (l Log) Diff() Node {
	return Divide{l.inner.Diff(), l.inner}
}

func (l Log) Simplify() Node {
	return Log{l.inner.Simplify()}
}

func (l Log) PrefixString() string {
	return fmt.Sprintf("(log %s)", l.inner.PrefixString())
}

type Exp struct {
	exp Node
}

func (e Exp) Diff() Node {
	return Multiply{Exp{e.exp}, e.exp.Diff()}
}

func (e Exp) Simplify() Node {
	return Exp{e.exp.Simplify()}
}

func (e Exp) PrefixString() string {
	return fmt.Sprintf("(exp %s)", e.exp.PrefixString())
}

type Pow struct {
	base Node
	exp  Node
}

func (e Pow) Simplify() Node {
	exp := e.exp.Simplify()
	base := e.base.Simplify()
	switch exp := exp.(type) {
	case Literal:
		if exp.value == 0 {
			return Literal{1}
		} else if exp.value == 1 {
			return base
		}
		switch base := base.(type) {
		case Literal:
			return Literal{math.Pow(base.value, exp.value)}
		}
	}
	return Pow{base, exp}
}

// d/dx(base^exp) = (exp * d/dx(base) + base * ln(base) * d/dx(exp)) * base^(exp - 1)

func (e Pow) Diff() Node {
	return Multiply{
		Add{
			Multiply{
				e.exp,
				e.base.Diff(),
			},
			Multiply{
				Multiply{
					e.base,
					Log{e.base},
				},
				e.exp.Diff()},
		},
		Pow{
			e.base,
			Subtract{
				e.exp,
				Literal{1},
			},
		},
	}
}

func (e Pow) PrefixString() string {
	return fmt.Sprintf("(^ %s %s)", e.base.PrefixString(), e.exp.PrefixString())
}

// Container types

type Literal struct {
	value float64
}

func (l Literal) Diff() Node {
	return Literal{0}
}

func (l Literal) Simplify() Node {
	return l
}

func (l Literal) PrefixString() string {
	if l.value == float64(int(l.value)) {
		return fmt.Sprintf("%d", int(l.value))
	} else {
		return fmt.Sprintf("%f", l.value)
	}
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

func (v Identifier) PrefixString() string {
	return v.name
}

// Trig functions

type Sin struct {
	inner Node
}

func (s Sin) Diff() Node {
	return Multiply{s.inner.Diff(), Cos{s.inner}}
}

func (s Sin) Simplify() Node {
	return Sin{s.inner.Simplify()}
}

func (s Sin) PrefixString() string {
	return fmt.Sprintf("(sin %s)", s.inner.PrefixString())
}

type Cos struct {
	inner Node
}

func (c Cos) Diff() Node {
	return Multiply{c.inner.Diff(), Multiply{Literal{-1}, Sin{c.inner}}}
}

func (c Cos) Simplify() Node {
	return Cos{c.inner.Simplify()}
}

func (c Cos) PrefixString() string {
	return fmt.Sprintf("(cos %s)", c.inner.PrefixString())
}

type Tan struct {
	inner Node
}

func (t Tan) Diff() Node {
	return Multiply{t.inner.Diff(), Add{Literal{1}, Pow{Tan{t.inner}, Literal{2}}}}
}

func (t Tan) Simplify() Node {
	return Tan{t.inner.Simplify()}
}

func (t Tan) PrefixString() string {
	return fmt.Sprintf("(tan %s)", t.inner.PrefixString())
}
