package definitions_test

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/zoul0813/go-prosper/pkg/definitions"
	"github.com/zoul0813/go-prosper/pkg/memory"
)

func TestDivMod(t *testing.T) {
	for _, tt := range []struct {
		left, right, remainder, quotient int
	}{
		{7, 3, 1, 2}, {-7, 3, -1, -2}, {7, -3, 1, -2}, {6, 3, 0, 2},
	} {
		m := memory.NewMemory(16, 16)
		m.Stack().Push(tt.left)
		m.Stack().Push(tt.right)
		if _, err := definitions.Core["/mod"](m); err != nil {
			t.Fatal(err)
		}
		want := []int{tt.remainder, tt.quotient}
		if got := m.Stack().Values[:m.Stack().Index]; !reflect.DeepEqual(got, want) {
			t.Fatalf("%d %d /mod = %v, want %v", tt.left, tt.right, got, want)
		}
	}
}

func TestIntegerPrecision(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, tt := range []struct {
		word  string
		input []int
		want  int
	}{
		{"abs", []int{maxInt}, maxInt},
		{"abs", []int{-maxInt}, maxInt},
		{"abs", []int{0}, 0},
		{"min", []int{maxInt, maxInt - 1}, maxInt - 1},
		{"min", []int{-maxInt, -maxInt + 1}, -maxInt},
		{"max", []int{maxInt - 1, maxInt}, maxInt},
		{"max", []int{-maxInt + 1, -maxInt}, -maxInt + 1},
	} {
		t.Run(tt.word+"/"+strconv.Itoa(tt.want), func(t *testing.T) {
			m := memory.NewMemory(16, 16)
			for _, n := range tt.input {
				m.Stack().Push(n)
			}
			if _, err := definitions.Core[tt.word](m); err != nil {
				t.Fatal(err)
			}
			if m.Stack().Index != 1 || m.Stack().Peek() != tt.want {
				t.Fatalf("stack = %v, want [%d]", m.Stack().Values[:m.Stack().Index], tt.want)
			}
		})
	}
}
