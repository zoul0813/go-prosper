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
	LOOP_SIZE       = 64
)

type Mode int

const (
	Initializing Mode = iota
	Standard
)

var (
	loops     = make([]int, LOOP_SIZE)
	loopIndex = 0
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
	vm.Memory.Dict().DefineSystem("branch")
	vm.Memory.Dict().DefineSystem("branch?")
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

func (vm *VM) Process(words string, next string) (int, error) {
	splitWords := strings.Fields(words)
	comment := false
	skipBy := 0

	for wordIndex := 0; wordIndex < len(splitWords); wordIndex++ {
		word := splitWords[wordIndex]
		// fmt.Printf("word: '%v', index: %v, words: \"%v\"\n", word, wordIndex, words)
		if word == "" {
			continue
		}
		if comment {
			// fmt.Printf("skipping: %v (%v, %v)\n", word, skipNext, comment)
			if word == ")" {
				skipBy = 0
				comment = false
			}
			continue
		}
		if vm.compiling {
			if word == ";" {
				// fmt.Printf("Done. loops: %v:%v, %v %v\n", len(loops), cap(loops), loopIndex, loops)
				vm.EndCompile()
				err := vm.Memory.Dict().Define(strings.TrimSpace(vm.defining), strings.TrimSpace(vm.definition))
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
					switch word {
					default:
						vm.definition += " " + word
					}
				}
			}
		} else {
			switch strings.ToLower(word) {
			// data space
			case "constant":
				name, ok := getNext(splitWords, wordIndex+1, next)
				skipBy = 1
				// fmt.Printf("constant: \"%v\" (%v) %d, \"%v\"\n", name, strings.Join(splitWords, " "), wordIndex+1, next)
				if !ok {
					fmt.Printf("attempt to use zero-length string as a name")
					continue
				}
				value := vm.Memory.Stack().Pop()
				vm.Memory.Dict().Define(name, fmt.Sprintf("%d", value))

			// branching
			case "branch":
				// unconditional branch
				offset, ok := getNext(splitWords, wordIndex+1, next)
				skipBy = 1
				// if !vm.compiling {
				// 	fmt.Printf("compile-only word")
				// }
				if !ok {
					fmt.Printf("attempt to use nil branch offset")
					continue
				}
				skip, err := strconv.Atoi(offset)
				if err != nil {
					fmt.Printf("attempt to use non-int as branch offset")
					continue
				}
				// TODO: implement "debug" flag for outputting these branch locations?
				// nw, _ := getNext(splitWords, wordIndex+2, next)
				// fmt.Printf("branch: %v:%v, %v\n", offset, skip, nw)
				skipBy += skip
				// dst, _ := getNext(splitWords, wordIndex+skipBy, next)
				// fmt.Printf("  next: %v\n", dst)
			case "branch?":
				// condition branch
				condition := vm.Memory.Stack().Pop()
				offset, ok := getNext(splitWords, wordIndex+1, next)
				skipBy = 1
				// if !vm.compiling {
				// 	fmt.Printf("compile-only word")
				// }
				if !ok {
					fmt.Printf("attempt to use nil branch offset")
					continue
				}
				if condition == 0 {
					skip, err := strconv.Atoi(offset)
					if err != nil {
						fmt.Printf("attempt to use non-int as branch offset")
						continue
					}
					// TODO: implement "debug" flag for outputting these branch locations?
					// nw, _ := getNext(splitWords, wordIndex+2, next)
					// fmt.Printf("branch?: %v:%v, %v\n", offset, skip, nw)
					skipBy += skip
					// dst, _ := getNext(splitWords, wordIndex+skipBy, next)
					// fmt.Printf("  next: %v\n", dst)
				}
			// compiler
			case ":":
				vm.Compile()
			// system
			case "\\":
				comment = true
			case "(":
				comment = true
			case "see":
				lookup, ok := getNext(splitWords, wordIndex+1, next)
				skipBy = 1
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
				return 0, fmt.Errorf("bye")
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
							return 1, fmt.Errorf("undefined word\n>>>%v<<<", word)
						}
					} else if i, err := strconv.Atoi(word); err == nil {
						vm.Memory.Stack().Push(i)
					} else {
						return 1, fmt.Errorf("undefined word\n>>>%v<<<", word)
					}
				} else {
					if def.Fn != nil {
						def.Fn(vm.Memory)
					} else {
						nextWord, _ := getNext(splitWords, wordIndex+1, next)
						skip, err := vm.Process(def.Definiton, nextWord)
						if err != nil {
							return 1, err
						}
						// fmt.Printf("skipNext: %v (%v) \"%v\"\n", skip, wordIndex, nextWord)
						skipBy = skip
					}
				}
			}
		}
		wordIndex += skipBy
		skipBy = 0
	}
	return skipBy, nil
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
