package vm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zoul0813/go-forth/pkg/dict"
	"github.com/zoul0813/go-forth/pkg/stack"
)

const (
	STACK_SIZE = 2048
)

type Mode int

const (
	Initializing Mode = iota
	Standard
)

type VM struct {
	ip          int // instruction pointer
	Stack       *stack.Stack
	ReturnStack *stack.Stack
	Dict        *dict.Dict

	compiling  bool
	defining   string
	definition string
	mode       Mode
}

func NewVM() *VM {
	return &VM{
		ip:          0,
		Stack:       stack.NewStack(STACK_SIZE),
		ReturnStack: stack.NewStack(STACK_SIZE),
		Dict:        dict.NewDict(),
		compiling:   false,
		defining:    "",
		definition:  "",
		mode:        Initializing,
	}
}

func (vm *VM) Process(words string) (bool, error) {
	splitWords := strings.Split(words, " ")
	skipNext := false
	for wordIndex, word := range splitWords {
		if skipNext {
			skipNext = false
			continue
		}
		if vm.compiling {
			if word == ";" {
				// fmt.Printf("Done.\n")
				vm.EndCompile()
				ok := vm.Dict.Define(vm.defining, vm.definition)
				if !ok {
					fmt.Printf("error: failed to define %v ", vm.defining)
				} else {
					if vm.mode == Standard {
						fmt.Printf("defined %v\n", vm.defining)
					}
				}
			} else {
				if vm.defining == "" {
					// fmt.Printf("Defining %v\n", word)
					vm.defining = word
				} else {
					// fmt.Printf("Compiling... %v\n", word)
					vm.definition += " " + word
				}
			}
		} else {
			switch strings.ToLower(word) {
			// arithmetic
			case "+":
				// add
				right := vm.Stack.Pop()
				left := vm.Stack.Pop()
				vm.Stack.Push(left + right)
			case "-":
				// sub
				right := vm.Stack.Pop()
				left := vm.Stack.Pop()
				vm.Stack.Push(left - right)
			case "*":
				// mul
				right := vm.Stack.Pop()
				left := vm.Stack.Pop()
				vm.Stack.Push(left * right)
			case "/":
				// div
				right := vm.Stack.Pop()
				left := vm.Stack.Pop()
				vm.Stack.Push(left / right)
			// stack
			case "depth":
				vm.Stack.Push(vm.Stack.Index)
			case "dup":
				vm.Stack.Dup()
			case "drop":
				vm.Stack.Pop()
			case "swap":
				vm.Stack.Swap()
			case "over":
				value := vm.Stack.PeekN(2)
				vm.Stack.Push(value)
			case "rot":
				vm.Stack.Rotate()
			case ">r":
				vm.ReturnStack.Push(vm.Stack.Pop())
			case "r>":
				vm.Stack.Push(vm.ReturnStack.Pop())
			case "clear":
				vm.Stack.Clear()
				// output
			case "(emit)": // emit ???
				// print
				i := vm.Stack.Pop()
				fmt.Printf("%d ", i)
			case "(showstack)":
				// inspect
				fmt.Printf("<%d> %v ", vm.Stack.Index, vm.Stack.Values[:vm.Stack.Index])
			// branching
			case "branch":
				// unconditional branch
				notYetImplemented(word)
			case "branch?":
				// condition branch
				notYetImplemented(word)
			// compiler
			case ":":
				vm.Compile()
			// system
			case "see":
				skipNext = true
				lookup := splitWords[wordIndex+1]
				def, ok := vm.Dict.Lookup(lookup)
				if !ok {
					fmt.Printf("undefined word: %v %v", word, def)
				} else {
					fmt.Printf(": %v\n  %v ;", lookup, def)
				}
			case "words":
				for word, _ := range vm.Dict.Words {
					fmt.Printf("%v ", word)
				}
			case "debug":
				fmt.Printf("\n\nDEBUG\n")
				fmt.Printf("  Stack: <%d> %v\n", vm.Stack.Index, vm.Stack.Values[:vm.Stack.Index])
				fmt.Printf("  Return: <%d> %v\n", vm.ReturnStack.Index, vm.ReturnStack.Values[:vm.ReturnStack.Index])
				fmt.Printf("  Mode: %v\n", vm.compiling)
			case "bye":
				return false, fmt.Errorf("bye")
			default:
				def, ok := vm.Dict.Lookup(word)
				if !ok {
					if i, err := strconv.Atoi(word); err == nil {
						vm.Stack.Push(i)
					}
				} else {
					// fmt.Printf("Found: %v = %v", word, def)
					vm.Process(def.Definiton)
				}
			}
		}
	}
	return true, nil
}

func (vm *VM) Compile() {
	vm.compiling = true
	vm.defining = ""
	vm.definition = ""
}

func (vm *VM) EndCompile() {
	vm.compiling = false
}

func notYetImplemented(word string) {
	fmt.Printf("error: not yet implemented: %v", word)
}

func (vm *VM) Mode() Mode {
	return vm.mode
}

func (vm *VM) SetMode(mode Mode) Mode {
	vm.mode = mode
	return vm.mode
}
