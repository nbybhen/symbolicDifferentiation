package main

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

func (s *Stack) Parse(val Node) {
	switch val.(type) {
	case Add:
		r, l := s.Pop(), s.Pop()
		s.Push(Add{r, l})
	case Multiply:
		r, l := s.Pop(), s.Pop()
		s.Push(Multiply{r, l})
	case Divide:
		r, l := s.Pop(), s.Pop()
		s.Push(Divide{r, l})
	case Subtract:
		r, l := s.Pop(), s.Pop()
		s.Push(Subtract{r, l})
	case Cos:
		i := s.Pop()
		s.Push(Cos{i})
	case Tan:
		i := s.Pop()
		s.Push(Tan{i})
	case Sin:
		i := s.Pop()
		s.Push(Sin{i})
	case Pow:
		base := s.Pop()
		exp := s.Pop()
		s.Push(Pow{base, exp})
	case Exp:
		exp := s.Pop()
		s.Push(Exp{exp})
	case Log:
		i := s.Pop()
		s.Push(Log{i})
	}
}
