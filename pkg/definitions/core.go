package definitions

import (
	"fmt"
	"math"
	"time"
)

var Core = map[string]DefinitionFn{
	// arithmetic
	"+": func(m memory) (bool, error) {
		// add
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left + right)
		return false, nil
	},
	"-": func(m memory) (bool, error) {
		// sub
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left - right)
		return false, nil
	},
	"*": func(m memory) (bool, error) {
		// mul
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left * right)
		return false, nil
	},
	"/": func(m memory) (bool, error) {
		// div
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left / right)
		return false, nil
	},
	"negate": func(m memory) (bool, error) {
		m.Stack().Push(-m.Stack().Pop())
		return false, nil
	},
	"abs": func(m memory) (bool, error) {
		abs := int(math.Abs(float64(m.Stack().Pop())))
		m.Stack().Push(abs)
		return false, nil
	},
	"min": func(m memory) (bool, error) {
		n1 := float64(m.Stack().Pop())
		n2 := float64(m.Stack().Pop())
		min := int(math.Min(n1, n2))
		m.Stack().Push(min)
		return false, nil
	},
	"max": func(m memory) (bool, error) {
		n1 := float64(m.Stack().Pop())
		n2 := float64(m.Stack().Pop())
		max := int(math.Max(n1, n2))
		m.Stack().Push(max)
		return false, nil
	},

	// stack
	"depth": func(m memory) (bool, error) {
		m.Stack().Push(m.Stack().Index)
		return false, nil
	},
	"dup": func(m memory) (bool, error) {
		m.Stack().Dup()
		return false, nil
	},
	"drop": func(m memory) (bool, error) {
		m.Stack().Pop()
		return false, nil
	},
	"swap": func(m memory) (bool, error) {
		m.Stack().Swap()
		return false, nil
	},
	"over": func(m memory) (bool, error) {
		value := m.Stack().PeekN(2)
		m.Stack().Push(value)
		return false, nil
	},
	"rot": func(m memory) (bool, error) {
		m.Stack().Rotate()
		return false, nil
	},
	">r": func(m memory) (bool, error) {
		m.ReturnStack().Push(m.Stack().Pop())
		return false, nil
	},
	"r>": func(m memory) (bool, error) {
		m.Stack().Push(m.ReturnStack().Pop())
		return false, nil
	},
	"clear": func(m memory) (bool, error) {
		m.Stack().Clear()
		return false, nil
	},

	// data space
	"here": func(m memory) (bool, error) {
		m.Stack().Push(m.DataSpace().Index)
		return false, nil
	},
	",": func(m memory) (bool, error) {
		m.DataSpace().Push(m.Stack().Pop())
		return true, nil
	},
	"!": func(m memory) (bool, error) {
		addr := m.Stack().Pop()
		value := m.Stack().Pop()
		m.DataSpace().Put(addr, value)
		return false, nil
	},
	"@": func(m memory) (bool, error) {
		addr := m.Stack().Pop()
		value := m.DataSpace().Get(addr)
		m.Stack().Push(value)
		return false, nil
	},
	"cells": func(m memory) (bool, error) {
		n := m.Stack().Pop()
		// todo: make this work on a byte level?
		// m.Stack().Push(n * (strconv.IntSize / 8))
		m.Stack().Push(n * 1)
		return false, nil
	},
	"allot": func(m memory) (bool, error) {
		size := m.Stack().Pop()
		m.DataSpace().Index += size
		// todo: error check overflow?
		return false, nil
	},

	// time
	"ms": func(m memory) (bool, error) {
		ms := m.Stack().Pop()
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return false, nil
	},
	"utime": func(m memory) (bool, error) {
		ms := time.Now().UnixNano()
		us := int(ms / 1000)

		m.Stack().Push(us)
		m.Stack().Push(0) // todo: why does gforth do this?
		return false, nil
	},

	// output
	"(emit)": func(m memory) (bool, error) {
		i := m.Stack().Pop()
		fmt.Printf("%d ", i)
		return false, nil
	},
	"(showstack)": func(m memory) (bool, error) {
		fmt.Printf("<%d> %v ", m.Stack().Index, m.Stack().Values[:m.Stack().Index])
		return false, nil
	},

	// branching
	"branch": func(m memory) (bool, error) {
		// unconditional branch
		return false, notYetImplemented("")
	},
	"branch?": func(m memory) (bool, error) {
		// condition branch
		return false, notYetImplemented("")
	},
}

func notYetImplemented(word string) error {
	return fmt.Errorf("error: not yet implemented: %v", word)
}
