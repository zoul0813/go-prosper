package stack

type Stack struct {
	Index     int
	Values    []int
	stackSize int
}

func NewStack(size int) *Stack {
	stack := &Stack{
		Index:     0,
		Values:    make([]int, size),
		stackSize: size,
	}

	return stack
}

func (s *Stack) Clear() {
	s.Index = 0
}

func (s *Stack) Put(i int, b int) {
	s.Values[i] = b
}

func (s *Stack) Get(i int) int {
	return s.Values[i]
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

func (s *Stack) Allot(size int) (bool, error) {
	s.Index += size
	// fmt.Printf("allot:\n")
	// fmt.Printf("    size: %v\n", size)
	// fmt.Printf("     idx: %v\n", s.Index)
	// fmt.Printf("     cap: %v\n", cap(s.Values))
	// fmt.Printf("     len: %v\n", len(s.Values))
	// if we need more space, increase stack by stackSize
	if s.Index > cap(s.Values)-1 {
		grow := s.stackSize // amount to grow

		// grow until we've added enough to cover the new allot size
		for (cap(s.Values) + grow - 1) < s.Index {
			grow += s.stackSize
		}
		// fmt.Printf("    grow: %v\n", grow)
		// make a new slice
		new := make([]int, len(s.Values), cap(s.Values)+grow)
		// copy the old slice into the new slice
		copy(new, s.Values)

		// update values
		s.Values = new[0:cap(new)]
		// fmt.Printf("  resize:\n")
		// fmt.Printf("    size: %v\n", size)
		// fmt.Printf("     idx: %v\n", s.Index)
		// fmt.Printf("     cap: %v\n", cap(s.Values))
		// fmt.Printf("     len: %v\n", len(s.Values))
		// fmt.Printf("     new: %v\n", cap(new))
		// fmt.Printf("    nlen: %v\n", len(new))
	}
	return true, nil
}
