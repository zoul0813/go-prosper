package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/zoul0813/go-forth/pkg/memory"
)

const (
	STACK_SIZE      = 2048
	DATA_SPACE_SIZE = 1024
)

type Mode int

const (
	Initializing Mode = iota
	Standard
)

type VM struct {
	// ip          int // instruction pointer
	// Stack       *stack.Stack
	// ReturnStack *stack.Stack
	// Dict        *dict.Dict
	Memory memory.Memory

	compiling  bool
	defining   string
	definition string
	mode       Mode
}

func NewVM() *VM {
	vm := &VM{
		// ip:          0,
		Memory:     *memory.NewMemory(STACK_SIZE, DATA_SPACE_SIZE),
		compiling:  false,
		defining:   "",
		definition: "",
		mode:       Initializing,
	}

	// internal words, can not be overwritten
	vm.Memory.Dict().DefineSystem("constant")
	vm.Memory.Dict().DefineSystem("words")
	vm.Memory.Dict().DefineSystem("(debug)")
	vm.Memory.Dict().DefineSystem("bye")

	return vm
}

func getNext(words []string, index int, next string) (string, bool) {
	len := len(words)
	// fmt.Printf(
	// 	"getNext: [%v] %d < %d = %v \"%v\"\n",
	// 	strings.Join(words, ", "),
	// 	len,
	// 	index,
	// 	len < index,
	// 	next,
	// )
	var word string = ""
	if index < len {
		word = words[index]
		// fmt.Printf("index++ \"%v\"\n", word)
	} else if next != "" {
		word = next
		// fmt.Printf("next \"%v\"\n", word)
	}
	// fmt.Printf("getNext = (\"%v\", %v)\n", word, word != "")
	return word, word != ""
}

func (vm *VM) Process(words string, next string) (bool, error) {
	splitWords := strings.Split(words, " ")
	skipNext := false
	comment := false
	for wordIndex, word := range splitWords {
		if word == "" {
			continue
		}
		if comment || skipNext {
			// fmt.Printf("skipping: %v (%v, %v)\n", word, skipNext, comment)
			skipNext = false
			if word == ")" {
				comment = false
			}
			continue
		}
		if vm.compiling {
			if word == ";" {
				// fmt.Printf("Done.\n")
				vm.EndCompile()
				err := vm.Memory.Dict().Define(vm.defining, vm.definition)
				if err != nil {
					fmt.Printf("error: %v (%v)", err, vm.defining)
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
			// data space
			case "constant":
				skipNext = true
				name, ok := getNext(splitWords, wordIndex+1, next)
				// fmt.Printf("constant: \"%v\" (%v) %d, \"%v\"\n", name, strings.Join(splitWords, ", "), wordIndex+1, next)
				if !ok {
					fmt.Printf("attempt to use zero-length string as a name")
					continue
				}
				value := vm.Memory.Stack().Pop()
				vm.Memory.Dict().Define(name, fmt.Sprintf("%d", value))
			// compiler
			case ":":
				vm.Compile()
			// system
			case "\\":
				comment = true
			case "(":
				comment = true
			case "see":
				skipNext = true
				lookup, ok := getNext(splitWords, wordIndex+1, next)
				// fmt.Printf("see: \"%v\" %v\n", lookup, ok)
				if !ok {
					fmt.Printf("undefined word, missing")
					continue
				}
				def, ok := vm.Memory.Dict().Lookup(lookup)
				if !ok {
					fmt.Printf("undefined word: \"%v\" %v", lookup, def)
					continue
				}
				fmt.Printf("%v", def)
			case "words":
				for word, _ := range vm.Memory.Dict().Words {
					fmt.Printf("%v ", word)
				}
			case "(debug)":
				fmt.Printf("\n\nDEBUG\n")
				fmt.Printf("  Stack: <%d> %v\n", vm.Memory.Stack().Index, vm.Memory.Stack().Values[:vm.Memory.Stack().Index])
				fmt.Printf("  Return: <%d> %v\n", vm.Memory.ReturnStack().Index, vm.Memory.ReturnStack().Values[:vm.Memory.ReturnStack().Index])
				fmt.Printf("  Mode: %v\n", vm.compiling)
			case "bye":
				return false, fmt.Errorf("bye")
			default:
				def, ok := vm.Memory.Dict().Lookup(word)
				// fmt.Printf(" > \"%v\" (%v, %v, %v) => \"%v\"\n", word, ok, skipNext, comment, def)
				if !ok {
					// TODO: refactor Atoi/ParseUint into "ParseNumber"
					if strings.HasPrefix(word, "$") {
						_, i := utf8.DecodeRuneInString(word)
						hex := word[i:]
						// fmt.Printf("hex: %v (%v, %v)", hex, i, word)
						if i, err := strconv.ParseUint(hex, 16, strconv.IntSize); err == nil {
							vm.Memory.Stack().Push(int(i))
						} else {
							return false, fmt.Errorf("undefined word\n>>>%v<<<", word)
						}
					} else if i, err := strconv.Atoi(word); err == nil {
						vm.Memory.Stack().Push(i)
					} else {
						return false, fmt.Errorf("undefined word\n>>>%v<<<", word)
					}
				} else {
					if def.Fn != nil {
						def.Fn(vm.Memory)
					} else {
						nextWord, _ := getNext(splitWords, wordIndex+1, next)
						skip, err := vm.Process(def.Definiton, nextWord)
						if err != nil {
							return false, err
						}
						// fmt.Printf("skipNext: %v (%v) \"%v\"\n", skip, skipNext, nextWord)
						skipNext = skip
					}
				}
			}
		}
	}
	return skipNext, nil
}

func (vm *VM) Compile() {
	vm.compiling = true
	vm.defining = ""
	vm.definition = ""
}

func (vm *VM) EndCompile() {
	vm.compiling = false
}

func (vm *VM) Mode() Mode {
	return vm.mode
}

func (vm *VM) SetMode(mode Mode) Mode {
	vm.mode = mode
	return vm.mode
}
