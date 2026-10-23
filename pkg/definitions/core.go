package definitions

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

var Core = map[string]DefinitionFn{
	// arithmetic
	"+": func(m memory) (int, error) {
		// add
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left + right)
		return 0, nil
	},
	"-": func(m memory) (int, error) {
		// sub
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left - right)
		return 0, nil
	},
	"*": func(m memory) (int, error) {
		// mul
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left * right)
		return 0, nil
	},
	"/": func(m memory) (int, error) {
		// div
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left / right)
		return 0, nil
	},
	"mod": func(m memory) (int, error) {
		// div
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left % right)
		return 0, nil
	},
	"/mod": func(m memory) (int, error) {
		// div
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		quotient := left / right
		modulus := left % right
		m.Stack().Push(modulus)
		m.Stack().Push(quotient)
		return 0, nil
	},
	">": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		if left > right {
			m.Stack().Push(-1)
		} else {
			m.Stack().Push(0)
		}
		return 0, nil
	},
	"<": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		if left < right {
			m.Stack().Push(-1)
		} else {
			m.Stack().Push(0)
		}
		return 0, nil
	},
	"=": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		if left == right {
			m.Stack().Push(-1)
		} else {
			m.Stack().Push(0)
		}
		return 0, nil
	},
	"<>": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		if left != right {
			m.Stack().Push(-1)
		} else {
			m.Stack().Push(0)
		}
		return 0, nil
	},
	"um*": func(m memory) (int, error) {
		right := int64(m.Stack().Pop())
		left := int64(m.Stack().Pop())
		product := left * right

		var x, y int
		x = int(product >> 32)
		y = int(product & 0xFFFFFFFF)

		m.Stack().Push(y)
		m.Stack().Push(x)
		return 0, nil
	},

	// bitwise
	"&": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left & right)
		return 0, nil
	},
	"|": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left | right)
		return 0, nil
	},
	"^": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left ^ right)
		return 0, nil
	},
	"&^": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left &^ right)
		return 0, nil
	},
	"<<": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left << right)
		return 0, nil
	},
	">>": func(m memory) (int, error) {
		right := m.Stack().Pop()
		left := m.Stack().Pop()
		m.Stack().Push(left >> right)
		return 0, nil
	},
	"negate": func(m memory) (int, error) {
		m.Stack().Push(-m.Stack().Pop())
		return 0, nil
	},
	"abs": func(m memory) (int, error) {
		abs := int(math.Abs(float64(m.Stack().Pop())))
		m.Stack().Push(abs)
		return 0, nil
	},
	"min": func(m memory) (int, error) {
		n1 := float64(m.Stack().Pop())
		n2 := float64(m.Stack().Pop())
		min := int(math.Min(n1, n2))
		m.Stack().Push(min)
		return 0, nil
	},
	"max": func(m memory) (int, error) {
		n1 := float64(m.Stack().Pop())
		n2 := float64(m.Stack().Pop())
		max := int(math.Max(n1, n2))
		m.Stack().Push(max)
		return 0, nil
	},

	// stack
	"depth": func(m memory) (int, error) {
		m.Stack().Push(m.Stack().Index)
		return 0, nil
	},
	"pick": func(m memory) (int, error) {
		index := m.Stack().Pop()
		m.Stack().Pick(index)
		return 0, nil
	},
	"dup": func(m memory) (int, error) {
		m.Stack().Pick(0)
		return 0, nil
	},
	"drop": func(m memory) (int, error) {
		m.Stack().Pop()
		return 0, nil
	},
	"swap": func(m memory) (int, error) {
		m.Stack().Swap()
		return 0, nil
	},
	"over": func(m memory) (int, error) {
		value := m.Stack().PeekN(2)
		m.Stack().Push(value)
		return 0, nil
	},
	"rot": func(m memory) (int, error) {
		m.Stack().Rotate()
		return 0, nil
	},
	">r": func(m memory) (int, error) {
		m.ReturnStack().Push(m.Stack().Pop())
		return 0, nil
	},
	"r>": func(m memory) (int, error) {
		m.Stack().Push(m.ReturnStack().Pop())
		return 0, nil
	},
	"clear": func(m memory) (int, error) {
		m.Stack().Clear()
		return 0, nil
	},
	"fill": func(m memory) (int, error) {
		chr := m.Stack().Pop()
		u := m.Stack().Pop()
		addr := m.Stack().Pop()

		for i := 0; i < u; i++ {
			m.DataSpace().Put(addr+i, chr)
		}

		return 0, nil
	},

	// data space
	"here": func(m memory) (int, error) {
		m.Stack().Push(m.DataSpace().Index)
		// showStack(m)
		// fmt.Printf(" here: t %T, & %p, i %d, %v\n", m.Stack(), m.Stack(), m.Stack().Index, m.Stack().Values[:m.Stack().Index])
		return 0, nil
	},
	",": func(m memory) (int, error) {
		value := m.Stack().Pop()
		index := m.DataSpace().Index
		m.DataSpace().Allot(1)
		m.DataSpace().Put(index, value)
		// fmt.Printf("\n")
		// m.Stack().Print(",", m.NumberBase())
		// fmt.Printf("\n")

		// fmt.Printf("  ds,: t %T, & %p, i %d, %v\n", m.DataSpace(), m.DataSpace(), m.DataSpace().Index, m.DataSpace().Values[:m.DataSpace().Index])
		// fmt.Printf("    ,: t %T, & %p, i %d, %v\n", m.Stack(), m.Stack(), m.Stack().Index, m.Stack().Values[:m.Stack().Index])

		return 0, nil
		// return 1, nil
	},
	"!": func(m memory) (int, error) {
		addr := m.Stack().Pop()
		value := m.Stack().Pop()
		m.DataSpace().Put(addr, value)
		return 0, nil
	},
	"@": func(m memory) (int, error) {
		addr := m.Stack().Pop()
		value := m.DataSpace().Get(addr)
		m.Stack().Push(value)
		return 0, nil
	},
	"cells": func(m memory) (int, error) {
		n := m.Stack().Pop()
		// todo: make this work on a byte level?
		// m.Stack().Push(n * (strconv.IntSize / 8))
		m.Stack().Push(n * 1)
		return 0, nil
	},
	"allot": func(m memory) (int, error) {
		size := m.Stack().Pop()
		m.DataSpace().Allot(size)
		// todo: error check overflow?
		return 0, nil
	},

	// time
	"ms": func(m memory) (int, error) {
		ms := m.Stack().Pop()
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return 0, nil
	},
	"utime": func(m memory) (int, error) {
		ms := time.Now().UnixNano()
		us := int(ms / 1000)

		m.Stack().Push(us)
		m.Stack().Push(0) // todo: why does gforth do this?
		return 0, nil
	},

	// output
	"(clearscreen)": func(m memory) (int, error) {
		fmt.Print("\033[2J")
		return 0, nil
	},
	"(at-xy)": func(m memory) (int, error) {
		y := m.Stack().Pop()
		x := m.Stack().Pop()

		// ESC[{line};{column}H
		fmt.Printf("\033[%d;%dH", y, x)
		return 0, nil
	},
	"(emit)": func(m memory) (int, error) {
		i := m.Stack().Pop()
		n := strconv.FormatInt(int64(i), m.NumberBase())
		// fmt.Printf("base: %d, ", m.NumberBase())
		fmt.Printf("%v ", n)
		return 0, nil
	},
	"(showstack)": func(m memory) (int, error) {
		return showStack(m)
	},
	"(showreturn)": func(m memory) (int, error) {
		return showReturn(m)
	},
	"i.r": func(m memory) (int, error) {
		pad := m.Stack().Pop()
		n := m.Stack().Pop()
		// '     0'
		s := fmt.Sprintf("%d", n)
		for i := len(s); i < pad; i++ {
			s = " " + s
		}
		fmt.Printf("%s", s)
		return 0, nil
	},
}

func showStack(m memory) (int, error) {
	m.Stack().Print("", m.NumberBase())
	return 0, nil
}

func showReturn(m memory) (int, error) {
	m.ReturnStack().Print("", m.NumberBase())
	return 0, nil
}

func notYetImplemented(word string) error {
	return fmt.Errorf("error: not yet implemented: %v", word)
}
