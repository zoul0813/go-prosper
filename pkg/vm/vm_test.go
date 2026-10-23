package vm

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func processTest(t *testing.T, m *VM, input string) {
	t.Helper()
	if _, err := m.Process(GetWords(input), "", 0); err != nil {
		t.Fatalf("process %q: %v", input, err)
	}
}

func assertStack(t *testing.T, m *VM, want ...int) {
	t.Helper()
	s := m.Memory.Stack()
	got := s.Values[:s.Index]
	if len(got) != len(want) || (len(got) > 0 && !reflect.DeepEqual(got, want)) {
		t.Fatalf("stack = %v, want %v", got, want)
	}
}

func loadCore(t *testing.T) *VM {
	t.Helper()
	core, err := os.ReadFile("../../runtime/core.f")
	if err != nil {
		t.Fatal(err)
	}
	m := NewVM()
	for _, line := range strings.Split(string(core), "\n") {
		if !strings.HasPrefix(line, "\\") {
			processTest(t, m, line)
		}
	}
	assertStack(t, m)
	return m
}

func TestRotationAndDependentWords(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  []int
	}{
		{"1 2 3 rot", []int{2, 3, 1}},
		{"1 2 3 -rot", []int{3, 1, 2}},
		{"1 2 3 4 2swap", []int{3, 4, 1, 2}},
		{"1 2 3 4 2over", []int{1, 2, 3, 4, 1, 2}},
		{": indices 4 1 do i loop ; indices", []int{1, 2, 3}},
		{": nested 2 0 do 2 0 do i loop i loop ; nested", []int{0, 1, 0, 0, 1, 1}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			m := loadCore(t)
			processTest(t, m, tt.input)
			assertStack(t, m, tt.want...)
			if m.Memory.ReturnStack().Index != 0 {
				t.Fatal("return stack not empty")
			}
		})
	}
}

func TestUntil(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  []int
	}{
		{": count 0 begin 1 + dup 3 = until ; count", []int{3}},
		{": once begin 42 -1 until ; once", []int{42}},
		{": nested 0 begin 0 begin 1 + dup 2 = until + dup 4 = until ; nested", []int{4}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			m := NewVM()
			processTest(t, m, tt.input)
			assertStack(t, m, tt.want...)
		})
	}
}
