package stack

type Stack struct {
	Index  int
	Values []int
}

func NewStack(size int) *Stack {
	stack := &Stack{
		Index:  0,
		Values: make([]int, size),
	}

	return stack
}

func (s *Stack) Clear() {
	s.Index = 0
}

func (s *Stack) Push(b int) {
	s.Values[s.Index] = b
	s.Index++
}

func (s *Stack) Pop() int {
	s.Index--
	return s.Values[s.Index]
}

func (s *Stack) Peek() int {
	return s.Values[s.Index-1]
}

func (s *Stack) PeekN(index int) int {
	return s.Values[s.Index-index]
}

func (s *Stack) Dup() {
	s.Push(s.Peek())
}

func (s *Stack) Swap() {
	top := s.Pop()
	next := s.Pop()
	s.Push(top)
	s.Push(next)
}

func (s *Stack) Rotate() {
	outer := s.Pop()
	inner := s.Pop()
	value := s.Pop()
	s.Push(outer)
	s.Push(inner)
	s.Push(value)
}
