package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/zoul0813/go-forth/pkg/dict"
	"github.com/zoul0813/go-forth/pkg/stack"
)

func main() {
	fmt.Println("Go Forth")

	stack := stack.NewStack(1024)
	dict := dict.NewDict()
	compiling := false
	var defining string
	var definition string

	scanner := bufio.NewScanner(os.Stdin)

	var process func(string) bool
	process = func(words string) bool {
		for _, word := range strings.Split(words, " ") {
			switch word {
			// arithmetic
			case "+":
				// add
				right := stack.Pop()
				left := stack.Pop()
				stack.Push(left + right)
			case "-":
				// sub
				right := stack.Pop()
				left := stack.Pop()
				stack.Push(left - right)
			case "*":
				// mul
				right := stack.Pop()
				left := stack.Pop()
				stack.Push(left * right)
			case "/":
				// div
				right := stack.Pop()
				left := stack.Pop()
				stack.Push(left / right)
			// stack
			case "dup":
				stack.Push(stack.Peek())
			case "drop":
				stack.Pop()
			case "nip":
				outer := stack.Pop()
				stack.Pop()
				stack.Push(outer)
			case "rot":
				outer := stack.Pop()
				inner := stack.Pop()
				value := stack.Pop()
				stack.Push(outer)
				stack.Push(inner)
				stack.Push(value)
			// output
			case ".":
				// print
				i := stack.Pop()
				fmt.Printf("%d", i)
			case ".s":
				// inspect
				fmt.Printf("<%d> %v", stack.Index, stack.Values[:stack.Index])
			// compiler
			case ":":
				defining = ""
				compiling = true
			// system
			case "bye":
				fmt.Printf("bye!\n")
				return false
			default:
				def, ok := dict.Words[word]
				if !ok {
					if i, err := strconv.Atoi(word); err == nil {
						stack.Push(i)
					}
				} else {
					process(def)
				}
			}
		}
		return true
	}

	fmt.Print("> ")
	for scanner.Scan() {

		running := true
		input := scanner.Text()

		words := strings.Split(strings.ToLower(input), " ")

		for _, word := range words {
			if compiling {
				if word == ";" {
					ok := dict.Define(defining, definition)
					if !ok {
						fmt.Printf("error: failed to define %v ", defining)
					} else {
						fmt.Printf("defined %v", defining)
					}
					compiling = false
				} else {
					if defining == "" {
						defining = word
					} else {
						definition += " " + word
					}
				}
			} else {
				running = process(word)

				if !running {
					break
				}
			}
		}
		if !running {
			break
		}
		fmt.Printf("  ok\n> ")
	}

}
