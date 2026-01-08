package memory

import (
	"github.com/zoul0813/go-prosper/pkg/dict"
	"github.com/zoul0813/go-prosper/pkg/stack"
)

type Memory struct {
	stack       *stack.Stack
	returnStack *stack.Stack
	dataSpace   *stack.Stack
	dict        *dict.Dict
	numberBase  int
}

func NewMemory(stackSize int, dataSize int) *Memory {
	return &Memory{
		stack:       stack.NewStack(stackSize),
		returnStack: stack.NewStack(stackSize),
		dataSpace:   stack.NewStack(dataSize),
		dict:        dict.NewDict(),
		numberBase:  10,
	}
}

func (m *Memory) Recover() {
	m.stack.Index = 0
	m.returnStack.Index = 0
}

func (m Memory) Stack() *stack.Stack {
	return m.stack
}
func (m Memory) ReturnStack() *stack.Stack {
	return m.returnStack
}

func (m Memory) DataSpace() *stack.Stack {
	return m.dataSpace
}

func (m Memory) Dict() *dict.Dict {
	return m.dict
}

func (m Memory) NumberBase() int {
	return m.numberBase
}

func (m *Memory) SetNumberBase(base int) {
	m.numberBase = base
}

// func (m *Memory) Allot(size int) error {
// 	// ( n -- )
// 	// reserve n units of data space
// 	m.dataSpace.Index += size
// 	return nil
// }

// func (m *Memory) Cell(x int) error {
// 	// , ( x -- )
// 	// reserve one cell and store x in the cell
// 	m.dataSpace.Push(x)
// 	return nil
// }

// func (m *Memory) Char(x byte) error {
// 	// C, ( char -- )
// 	// reserve one character and store x in the space
// 	m.dataSpace.Push(int(x))
// 	return nil
// }

// func (m *Memory) Align() error {
// 	// Implementors may define these words as no-ops on systems for which they aren't functional.
// 	return nil
// }

// func (m *Memory) Aligned() error {
// 	// ( addr -- a-addr )
// 	// return addr
// 	return nil
// }
