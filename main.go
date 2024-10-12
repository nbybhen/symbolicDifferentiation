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

	_ Node = (*Literal)(nil)
	_ Node = (*Variable)(nil)

	_ Node = (*Sin)(nil)
	_ Node = (*Cos)(nil)
	_ Node = (*Tan)(nil)
)

type Node interface {
	Diff() Node
}

// Basic operations (+ * - / ^)

type Add struct {
	left  Node
	right Node
}

func (a Add) Diff() Node {
	return Add{a.left.Diff(), a.right.Diff()}
}

type Subtract struct {
	left  Node
	right Node
}

func (s Subtract) Diff() Node {
	return Subtract{s.left.Diff(), s.right.Diff()}
}

type Multiply struct {
	left  Node
	right Node
}

func (m Multiply) Diff() Node {
	return Multiply{m.left.Diff(), m.right.Diff()}
}

type Divide struct {
	left  Node
	right Node
}

func (d Divide) Diff() Node {
	return Divide{d.left.Diff(), d.right.Diff()}
}

type Log struct {
	inner Node
}

// d/dx(ln(x)) = 1/ln(x) * x'

func (l Log) Diff() Node {
	return Multiply{Divide{Literal{1}, Log{l.inner}}, l.inner.Diff()}
}

type Exponent struct {
	base Node
	exp  Node
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

type Variable struct {
	name string
}

func (v Variable) Diff() Node {
	return Literal{1}
}

// Trig functions

type Sin struct {
	inner Node
}

func (s Sin) Diff() Node {
	return Multiply{Cos{s.inner}, s.inner.Diff()}
}

type Cos struct {
	inner Node
}

func (c Cos) Diff() Node {
	return Multiply{Multiply{Literal{-1}, Sin{c.inner}}, c.inner.Diff()}
}

type Tan struct {
	inner Node
}

func (t Tan) Diff() Node {
	return Divide{t.inner.Diff(), Exponent{Cos{t.inner}, Literal{2}}}
}

func main() {
	test := Add{Literal{3}, Literal{2}}
	fmt.Println("Test: ", test)

	res := test.Diff()
	fmt.Println("Result: ", res)
}
