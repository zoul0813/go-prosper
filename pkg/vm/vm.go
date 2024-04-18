package vm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zoul0813/go-forth/pkg/memory"
	"github.com/zoul0813/go-forth/pkg/numbers"
	"github.com/zoul0813/go-forth/pkg/stack"
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

type VM struct {
	// ip          int // instruction pointer
	// Stack       *stack.Stack
	// ReturnStack *stack.Stack
	// Dict        *dict.Dict
	Memory *memory.Memory

	compiling  bool
	defining   string
	definition string
	mode       Mode
}

func NewVM() *VM {
	vm := &VM{
		// ip:          0,
		Memory:     memory.NewMemory(STACK_SIZE, DATA_SPACE_SIZE),
		compiling:  false,
		defining:   "",
		definition: "",
		mode:       Initializing,
	}

	// temporary
	vm.Memory.Dict().DefineSystem("cr")

	// internal words, can not be overwritten

	internal := []string{
		// control flow
		"branch",
		"?branch",
		"[char]",
		"do",
		"begin",
		"until",
		"loop",
		"i",
		// variables
		"constant",
		".\"",
		// system
		"words",
		"(debug)",
		"bye",
	}

	for _, word := range internal {
		vm.Memory.Dict().DefineSystem(word)
	}

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

func getString(word, words string, stringIndex, wordIndex, skipBy int) (string, int, int) {
	splitWords := strings.Fields(words)
	index := len(word) + 1 + stringIndex + strings.Index(words[stringIndex:], word)
	str := make([]rune, len(words)-len(word)-index)
	skip := 0
	stringIndex += len(word) + 1

	if index < 0 {
		fmt.Printf("string not found: \"%v\" in \"%v\"", word, words)
		return "!INVALID", stringIndex, skip
	}

	for i, c := range words[index:] {
		if string(c) != "\"" {
			str = append(str, c)
		} else {
			stringIndex += i
			break
		}
	}
	sw := strings.Fields(string(str))
	skip += len(sw)
	nw, _ := getNext(splitWords, wordIndex+(skipBy+skip)+1, "")
	if nw == "\"" {
		skip++
	}
	// fmt.Printf("\nstring: '%s'\n", string(str))
	return string(str), stringIndex, skip
	// fmt.Printf(" (next: %v, %d)", nw, skipBy)
}

func (vm *VM) Process(words string, next string, skipBy int) (int, error) {
	splitWords := strings.Fields(words)
	comment := false
	stringIndex := 0

	// fmt.Printf("Stack: t %T, & %p, i %d, %v\n", vm.Memory.Stack(), vm.Memory.Stack(), vm.Memory.Stack().Index, vm.Memory.Stack().Values[:vm.Memory.Stack().Index])

	for wordIndex := skipBy; wordIndex < len(splitWords); wordIndex++ {
		word := splitWords[wordIndex]
		// fmt.Printf(
		// 	"word: '%v', index: %v  <%v> %v, <%v> %v\n",
		// 	word,
		// 	wordIndex,
		// 	vm.Memory.Stack().Index,
		// 	vm.Memory.Stack().Values[:vm.Memory.Stack().Index],
		// 	vm.Memory.ReturnStack().Index,
		// 	vm.Memory.ReturnStack().Values[:vm.Memory.ReturnStack().Index],
		// )
		// fmt.Printf("word: '%v', index: %v, words: \"%v\"\n", word, wordIndex, words)
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
				err := vm.EndCompile()

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
					case ".\"":
						str, index, skip := getString(word, words, stringIndex, wordIndex, skipBy)
						stringIndex = index
						// fmt.Printf("%s", str)
						vm.definition += " " + word + " " + str + "\""
						// fmt.Printf("defining: str, %v, %v\n", word, str)
						skipBy += skip
					default:
						vm.definition += " " + word
					}
				}
			}
		} else {
			switch strings.ToLower(word) {
			// number base
			case "(decimal)":
				vm.Memory.SetNumberBase(10)
			case "(octal)":
				vm.Memory.SetNumberBase(8)
			case "(hex)":
				vm.Memory.SetNumberBase(16)
			case "(binary)":
				vm.Memory.SetNumberBase(2)
			// data space
			case "cr":
				{
					fmt.Printf("\n")
				}
			case "constant":
				name, ok := getNext(splitWords, wordIndex+1, next)
				skipBy += 1
				// vm.Memory.Stack().Print("C", vm.Memory.NumberBase())
				value := vm.Memory.Stack().Pop()
				// vm.Memory.Stack().Print("C", vm.Memory.NumberBase())
				// fmt.Printf("constant: \"%v\" (%v) %d %d, \"%v\", %v !\n", name, strings.Join(splitWords, " "), skipBy, wordIndex+1, next, value)
				if !ok {
					fmt.Printf("attempt to use zero-length string as a name")
					continue
				}
				vword := fmt.Sprintf("%d", value)
				vm.Memory.Dict().Define(name, vword, vword)
				// vm.Memory.Stack().Print("C", vm.Memory.NumberBase())
			case ".\"":
				// fmt.Printf("\nwor: '%v'\nsid: %v\nwid: %v\n", word, stringIndex, wordIndex)
				str, index, skip := getString(word, words, stringIndex, wordIndex, skipBy)
				stringIndex = index
				fmt.Printf("%s ", string(str))
				skipBy += skip

				// fmt.Printf("\nstr: '%v'\nidx: %v\nski: %v\n", str, index, skip)

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
			case "?branch":
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
					// fmt.Printf("?branch: %v:%v, %v\n", offset, skip, nw)
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
				fmt.Printf("   Stack: <%d> %v\n", vm.Memory.Stack().Index, vm.Memory.Stack().Values[:vm.Memory.Stack().Index])
				fmt.Printf("  Return: <%d> %v\n", vm.Memory.ReturnStack().Index, vm.Memory.ReturnStack().Values[:vm.Memory.ReturnStack().Index])
				fmt.Printf("    Data: <%d> %v\n", vm.Memory.DataSpace().Index, vm.Memory.DataSpace().Values[:vm.Memory.DataSpace().Index])
				fmt.Printf("    Mode: %v\n", vm.compiling)
			case "bye":
				return 0, fmt.Errorf("bye")
			default:
				def, ok := vm.Memory.Dict().Lookup(word)
				// fmt.Printf(" > \"%v\" (%v, %v, %v) => \"%v\"\n", word, ok, skipBy, comment, def)
				// fmt.Printf(" word: %v : %d\n", word, skipBy)
				if ok {
					if def.Fn != nil {
						skip, err := def.Fn(vm.Memory)
						if err != nil {
							return skip, err
						}
					} else {
						nextWord, _ := getNext(splitWords, wordIndex+1, next)
						skip, err := vm.Process(def.Compiled, nextWord, 0) // we just assume we aren't skipping here ... yes?
						if err != nil {
							return 1, err
						}
						// fmt.Printf("skipNext: %v (%v) \"%v\"\n", skip, wordIndex, nextWord)
						skipBy += skip
					}
				} else {
					// TODO: refactor Atoi/ParseUint into "ParseNumber"
					n, ok := numbers.Parse(word, vm.Memory.NumberBase())
					if ok {
						vm.Memory.Stack().Push(n)
					} else {
						return 0, fmt.Errorf("undefined word\n>>>%v<<<", word)
					}
				}
			}
		}
		wordIndex += skipBy
		if wordIndex >= len(splitWords) {
			return skipBy, nil
		}
		skipBy = 0
	}
	return skipBy, nil
}

func (vm *VM) Compile() {
	vm.compiling = true
	vm.defining = ""
	vm.definition = ""
}

func (vm *VM) EndCompile() error {
	definition := strings.TrimSpace(vm.definition)
	compiled, err := vm.parse(definition)

	if err == nil {
		err = vm.Memory.Dict().Define(strings.TrimSpace(vm.defining), definition, compiled)
	}

	vm.compiling = false

	return err
}

func (vm *VM) parse(d string) (string, error) {
	words := strings.Fields(d)

	compiled := make([]string, 0, 1024)
	control := stack.NewStack(1024)
	// for _, word := range words {
	skipNext := false
	for wordIndex := 0; wordIndex < len(words); wordIndex++ {
		word := words[wordIndex]
		if skipNext {
			if word == ")" {
				skipNext = false
			}
			continue
		}
		switch strings.ToLower(word) {
		case "(":
			skipNext = true

		case "[char]":
			if len(words) >= wordIndex+1 {
				char := words[wordIndex+1]
				n := int(char[0])
				wordIndex++
				compiled = append(compiled, strconv.Itoa(n))
			} else {
				return "", fmt.Errorf("[char] ")
			}

		case "if":
			compiled = append(compiled, "?branch")
			compiled = append(compiled, "?") // location of 'else'
			control.Push(len(compiled))

		case "else":
			compiled = append(compiled, "branch")
			compiled = append(compiled, "?") // location of 'then'
			jump := control.Pop()            // location of if ?branch index
			to := len(compiled)
			compiled[jump-1] = fmt.Sprintf("%d", to-2)
			control.Push(len(compiled))

		case "then":
			jump := control.Pop() // location of else ?branch index
			to := len(compiled) - jump
			compiled[jump-1] = fmt.Sprintf("%d", to)

		case "do":
			compiled = append(compiled, "2>R")
			control.Push(len(compiled))

		case "begin":
			control.Push(len(compiled))

		case "until":
			// jump out of loop if 0
			compiled = append(compiled, "?branch")
			compiled = append(compiled, "2")

			offset := control.Pop()
			jump := (offset - (len(compiled) + 2))
			// fmt.Printf("UNTIL : %v, %v, %v\n", offset, jump, len(compiled))
			compiled = append(compiled, "branch")
			compiled = append(compiled, fmt.Sprintf("%d", jump))

		case "loop":
			compiled = append(compiled, "2R>")
			compiled = append(compiled, "1")
			compiled = append(compiled, "+")
			compiled = append(compiled, "2DUP")
			compiled = append(compiled, "=")
			compiled = append(compiled, "ROT")
			compiled = append(compiled, "SWAP")
			compiled = append(compiled, "2>R")

			offset := control.Pop()
			jump := (offset - (len(compiled) + 2))
			// fmt.Printf("LOOP : %v, %v, %v\n", offset, jump, len(compiled))
			compiled = append(compiled, "?branch")
			compiled = append(compiled, fmt.Sprintf("%d", jump))

			compiled = append(compiled, "2R>")
			compiled = append(compiled, "2DROP")

		// inner loop index
		case "i":
			compiled = append(compiled, "2R@")
			compiled = append(compiled, "SWAP")
			compiled = append(compiled, "DROP")
		// outer loop index
		// case "j":
		// 	compiled = append(compiled, "2R>")
		// 	compiled = append(compiled, "2R>")
		// 	compiled = append(compiled, "2SWAP")
		// 	compiled = append(compiled, "2DROP")

		default:
			// fmt.Printf("PASSTHRU: '%v'\n", word)
			compiled = append(compiled, word)
		}
	}

	c := strings.Join(compiled, " ")

	// fmt.Printf(" PARSED: %v\n COMPILE: %v\n", d, c)
	return c, nil
}

func (vm *VM) Mode() Mode {
	return vm.mode
}

func (vm *VM) SetMode(mode Mode) Mode {
	vm.mode = mode
	return vm.mode
}
