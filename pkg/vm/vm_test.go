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

func TestCompiledValuesIgnoreExecutionBase(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  []int
	}{
		{"10 constant ten hex ten", []int{10}},
		{"-10 constant negative binary negative", []int{-10}},
		{": ten 10 ; hex ten", []int{10}},
		{"hex : sixteen 10 ; decimal sixteen", []int{16}},
		{": char [char] A ; binary char", []int{65}},
		{": indices 3 0 do i loop ; binary indices", []int{0, 1, 2}},
		{": named 10 constant ff ; hex named ff", []int{10}},
		{": conditional if 10 else 20 then ; binary -1 conditional 0 conditional", []int{10, 20}},
	} {
		t.Run(tt.input, func(t *testing.T) {
			m := loadCore(t)
			processTest(t, m, tt.input)
			assertStack(t, m, tt.want...)
		})
	}
}

func TestInvalidDefinitionsRecover(t *testing.T) {
	for _, body := range []string{
		"then", "else", "loop", "until", "if 42", "begin 42", "do 42",
		"begin then", "if loop", "do until", "if else else then", "( unclosed", "[char]",
	} {
		t.Run(body, func(t *testing.T) {
			m := NewVM()
			processTest(t, m, ": saved 7 ;")
			if _, err := m.Process(GetWords(": saved "+body+" ;"), "", 0); err == nil {
				t.Fatal("invalid definition accepted")
			}
			if m.compiling {
				t.Fatal("compilation mode not reset")
			}
			processTest(t, m, ": good 1 2 + ; saved good")
			assertStack(t, m, 7, 3)
		})
	}
	m := NewVM()
	if _, err := m.Process(GetWords(": ;"), "", 0); err == nil {
		t.Fatal("empty name accepted")
	}
}

func TestNestedControlFlow(t *testing.T) {
	m := loadCore(t)
	processTest(t, m, ": choose if 1 if 10 else 20 then else 30 then ; -1 choose 0 choose")
	assertStack(t, m, 10, 30)
	processTest(t, m, ": mixed 3 0 do i 1 = if 99 else i then loop ; mixed")
	assertStack(t, m, 10, 30, 0, 99, 2)
}
