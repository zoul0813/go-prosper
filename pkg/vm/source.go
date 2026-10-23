package vm

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrBye = errors.New("bye")

// GetWords tokenizes source for callers using Process directly. Execute also
// reports lexical errors and is the preferred entry point for source text.
func GetWords(input string) []string {
	words, _ := scanWords(input)
	return words
}

func scanWords(input string) ([]string, error) {
	runes := []rune(input)
	var words []string
	for i := 0; i < len(runes); {
		if unicode.IsSpace(runes[i]) {
			i++
			continue
		}
		// Existing core files also use repeated backslashes for comment banners.
		if runes[i] == '\\' {
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}
		start := i
		for i < len(runes) && !unicode.IsSpace(runes[i]) {
			i++
		}
		word := string(runes[start:i])
		switch strings.ToLower(word) {
		case "(":
			for i < len(runes) && runes[i] != ')' {
				i++
			}
			if i == len(runes) {
				return nil, fmt.Errorf("unclosed comment")
			}
			i++
		case ".\"", "abort\"":
			words = append(words, word)
			// The whitespace separating the parsing word from its text is consumed.
			if i < len(runes) {
				i++
			}
			start = i
			for i < len(runes) {
				if runes[i] == '\\' && i+1 < len(runes) {
					i += 2
					continue
				}
				if runes[i] == '"' {
					break
				}
				i++
			}
			if i == len(runes) {
				return nil, fmt.Errorf("%s: unclosed string", word)
			}
			i++
			words = append(words, string(runes[start:i]))
		default:
			words = append(words, word)
		}
	}
	return words, nil
}

// Execute processes source, allowing a definition to continue in the next call.
// Errors abort the current input and leave the interpreter ready for another.
func (vm *VM) Execute(source string) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
		}
		if err != nil {
			vm.compiling = false
			vm.definition = nil
			vm.Memory.Recover()
		}
	}()
	words, err := scanWords(source)
	if err != nil {
		return err
	}
	_, err = vm.Process(words, "", 0)
	return err
}

// Load processes a complete source file and rejects unfinished definitions.
func (vm *VM) Load(source string) error {
	if err := vm.Execute(source); err != nil {
		return err
	}
	if vm.compiling {
		vm.compiling = false
		vm.definition = nil
		vm.Memory.Recover()
		return fmt.Errorf("unfinished definition %q", vm.defining)
	}
	return nil
}
