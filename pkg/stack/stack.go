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
