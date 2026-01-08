package vm

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/zoul0813/go-prosper/pkg/memory"
	"github.com/zoul0813/go-prosper/pkg/numbers"
	"github.com/zoul0813/go-prosper/pkg/stack"
)

const (
	STACK_SIZE      = 2048
	DATA_SPACE_SIZE = 1024
	// LOOP_SIZE       = 64
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
	definition []string
	mode       Mode
}

func NewVM() *VM {
	vm := &VM{
		// ip:          0,
		Memory:    memory.NewMemory(STACK_SIZE, DATA_SPACE_SIZE),
		compiling: false,
		// defining:   "",
		// definition: make([]string, 0, STACK_SIZE),
		mode: Initializing,
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

func getWord(words []string, index int, next string) (string, bool) {
	len := len(words)
	var word string = ""
	if index < len {
		word = words[index]
	} else if next != "" {
		word = next
	}
	return word, word != ""
}

func GetWords(input string) []string {
	split := make([]string, 0, strings.Count(input, " "))
	word := make([]rune, 0, 256)

	l := len(input)

	next := func(i int) (rune, bool) {
		if i+1 < l {
			return rune(input[i+1]), true
		}
		return 0, false
	}

	isString := false

	for i := 0; i < len(input); i++ {
		c := rune(input[i])

		switch c {
		case '.':
			if n, ok := next(i); ok {
				if n == '"' {
					// fmt.Printf("found string\n")
					isString = true
					i += 2
					split = append(split, ".\"")
					word = make([]rune, 0, 256)
					continue
				}
			}
		case '\\':
			if isString {
				if n, ok := next(i); ok {
					if n == '"' {
						word = append(word, '\\')
						word = append(word, '"')
						i++
						continue
					}
				}
			}
		case '"':
			// fmt.Printf("string complete\n")
			isString = false
			word = append(word, '"')
			str := string(word)
			split = append(split, str)
			word = make([]rune, 0, 256)
			continue
		}

		if unicode.IsSpace(c) {
			if !isString {
				if len(word) > 0 {
					str := string(word)
					split = append(split, str)
					word = make([]rune, 0, 256)
				}
				continue
			}
		}
		word = append(word, c)
	}

	if len(word) > 0 {
		split = append(split, string(word))
	}
	return split
}

func (vm *VM) Process(words []string, next string, skipBy int) (int, error) {
	// fmt.Printf("split: %v\n", strings.Join(words, " | "))
	comment := false
	// stringIndex := 0

	// fmt.Printf("Stack: t %T, & %p, i %d, %v\n", vm.Memory.Stack(), vm.Memory.Stack(), vm.Memory.Stack().Index, vm.Memory.Stack().Values[:vm.Memory.Stack().Index])

	for index := skipBy; index < len(words); index++ {
		word := words[index]
		// fmt.Printf(
		// 	"word: '%v', index: %v  <%v> %v, <%v> %v\n",
		// 	word,
		// 	wordIndex,
		// 	vm.Memory.Stack().Index,
		// 	vm.Memory.Stack().Values[:vm.Memory.Stack().Index],
		// 	vm.Memory.ReturnStack().Index,
		// 	vm.Memory.ReturnStack().Values[:vm.Memory.ReturnStack().Index],
		// )
		// fmt.Printf("word: '%v', index: %v, input: \"%v\"\n", word, wordIndex, input)
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
					vm.definition = append(vm.definition, word)
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
				name, ok := getWord(words, index+1, next)
				skipBy += 1
				value := vm.Memory.Stack().Pop()
				// fmt.Printf("\n")
				// vm.Memory.Stack().Print("C", vm.Memory.NumberBase())
				// fmt.Printf("\n")
				// fmt.Printf("constant: \"%v\" (%v) sk: %d, idx: %d, val: \"%v\", %v !\n", name, strings.Join(words, " "), skipBy, index, next, value)
				if !ok {
					fmt.Printf("attempt to use zero-length string as a name")
					continue
				}
				vword := fmt.Sprintf("%d", value)
				vm.Memory.Dict().Define(name, vword, []string{vword})
				// vm.Memory.Stack().Print("C", vm.Memory.NumberBase())
			case ".\"":
				w, ok := getWord(words, index+1, next)
				if !ok {
					fmt.Printf("no string found")
					continue
				}
				skipBy += 1
				if len(w) > 0 {
					w = string(w[:len(w)-1])
					str, err := strconv.Unquote("\"" + w + "\"")
					if err != nil {
						fmt.Printf("failed to unquote '%v' -> %v", w, err)
						continue
					}
					fmt.Printf("%s", str)
				}

			// branching
			case "branch":
				// unconditional branch
				offset, ok := getWord(words, index+1, next)
				skipBy += 1
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
				// nw, _ := getWord(words, wordIndex+2, next)
				// fmt.Printf("branch: %v:%v, %v\n", offset, skip, nw)
				skipBy += skip
				// dst, _ := getWord(words, wordIndex+skipBy, next)
				// fmt.Printf("  next: %v\n", dst)
			case "?branch":
				// condition branch
				condition := vm.Memory.Stack().Pop()
				offset, ok := getWord(words, index+1, next)
				skipBy += 1
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
					// nw, _ := getWord(words, wordIndex+2, next)
					// fmt.Printf("?branch: %v:%v, %v\n", offset, skip, nw)
					skipBy += skip
					// dst, _ := getWord(words, wordIndex+skipBy, next)
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
				lookup, ok := getWord(words, index+1, next)
				skipBy += 1
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
				fmt.Printf("\n%v\n\n ", def)
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
						nextWord, _ := getWord(words, index+1, next)
						// fmt.Printf("\n")
						// vm.Memory.Stack().Print("PB", vm.Memory.NumberBase())
						// fmt.Printf("\nPw: %v", def.Compiled)
						// fmt.Printf("\n")
						skip, err := vm.Process(def.Compiled, nextWord, 0) // we just assume we aren't skipping here ... yes?
						// fmt.Printf("\n")
						// vm.Memory.Stack().Print("PA", vm.Memory.NumberBase())
						// fmt.Printf("\nPW: %v", def.Compiled)
						// fmt.Printf("\n")
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
		index += skipBy
		if index >= len(words) {
			// fmt.Printf("\nSkipping: sk: %v, idx: %v, len: %v %v\n", skipBy, index, index-(len(words)-1), words)
			return index - (len(words) - 1), nil
		}
		skipBy = 0
	}
	return skipBy, nil
}

func (vm *VM) Compile() {
	vm.compiling = true
	vm.defining = ""
	vm.definition = make([]string, 0, STACK_SIZE)
}

func (vm *VM) EndCompile() error {
	compiled, err := vm.parse(vm.definition)

	if err == nil {
		err = vm.Memory.Dict().Define(strings.TrimSpace(vm.defining), strings.Join(vm.definition, " "), compiled)
	}

	vm.compiling = false

	return err
}

func (vm *VM) parse(words []string) ([]string, error) {
	// fmt.Printf("Parsing...\n")
	// words := strings.Fields(d)
	// words := getWords(input)

	compiled := make([]string, 0, STACK_SIZE)
	control := stack.NewStack(STACK_SIZE)
	// stringIndex := 0
	// for _, word := range words {
	skipNext := false
	for index := 0; index < len(words); index++ {
		word := words[index]
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
			if len(words) >= index+1 {
				char, ok := getWord(words, index+1, "")
				if !ok {
					return []string{}, fmt.Errorf("[char] without character")
				}
				n := int(char[0])
				index++
				compiled = append(compiled, strconv.Itoa(n))
			} else {
				return []string{}, fmt.Errorf("[char] ")
			}

		case "if":
			compiled = append(compiled, "?branch")
			compiled = append(compiled, "?") // location of 'else'
			control.Push(len(compiled) - 1)

		case "else":
			compiled = append(compiled, "branch")
			compiled = append(compiled, "?") // location of 'then'
			jump := control.Pop()            // location of if ?branch index
			to := len(compiled) - jump - 1
			// fmt.Printf("\nelse: %v, %v\n", jump, to)
			compiled[jump] = fmt.Sprintf("%d", to)
			control.Push(len(compiled) - 1)

		case "then":
			jump := control.Pop() // location of else ?branch index
			to := len(compiled) - jump - 1
			// fmt.Printf("\nthen: %v, %v\n", jump, to)
			compiled[jump] = fmt.Sprintf("%d", to)

		case "do":
			compiled = append(compiled, "2>R")
			control.Push(len(compiled) - 1)

		case "begin":
			control.Push(len(compiled) - 1)

		case "until":
			// jump out of loop if 0
			compiled = append(compiled, "?branch")
			compiled = append(compiled, "2")

			offset := control.Pop()
			jump := (offset - (len(compiled) - 1 + 2))
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
			jump := (offset - (len(compiled) - 1 + 2))
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

	// c := strings.Join(compiled, " ")

	// fmt.Printf(" PARSED: %v\n COMPILE: %v\n", d, c)
	return compiled, nil
}

func (vm *VM) Mode() Mode {
	return vm.mode
}

func (vm *VM) SetMode(mode Mode) Mode {
	vm.mode = mode
	return vm.mode
}

//// DEPRECATED

// func getString(word, input string, stringIndex, wordIndex, skipBy int) (string, int, int) {
// 	// fmt.Printf("get: '%v', '%v', sid: %v, wid: %v, sk: %v", word, words, stringIndex, wordIndex, skipBy)
// 	// words := strings.Fields(words)
// 	words := getWords(input)
// 	index := len(word) + 1 + stringIndex + strings.Index(words[stringIndex:], word)
// 	str := make([]rune, len(words)-len(word)-index)
// 	skip := 0
// 	stringIndex += len(word) + 1

// 	if index < 0 {
// 		fmt.Printf("string not found: \"%v\" in \"%v\"\n", word, words)
// 		return "!INVALID", stringIndex, skip
// 	}

// 	for i, c := range words[index:] {
// 		if string(c) != "\"" {
// 			str = append(str, c)
// 		} else {
// 			stringIndex += i
// 			break
// 		}
// 	}
// 	// sw := strings.Fields(string(str))
// 	sw := getWords(string(str))
// 	skip += len(sw)
// 	nw, _ := getWord(words, wordIndex+(skipBy+skip)+1, "")
// 	if nw == "\"" {
// 		skip++
// 	}
// 	// fmt.Printf("get: '%s'\n", string(str))
// 	return string(str), stringIndex, skip
// 	// fmt.Printf(" (next: %v, %d)", nw, skipBy)
// }
