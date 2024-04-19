package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/zoul0813/go-forth/pkg/vm"
)

const (
	VERSION      = "0.0.0-alpha"
	HISTORY_SIZE = 32
)

var (
	machine      *vm.VM
	history      []string
	historyIndex int
)

func main() {
	fmt.Printf("Go Forth %v, Copyright (C) 2024 David Higgins <zoul0813@me.com>\nCompiled with %v\n", VERSION, runtime.Version())

	cwd, _ := os.Getwd()
	corePath := path.Join(cwd, "cmd/prosper/core.f")
	fmt.Printf("Core: %v\n", corePath)
	core, err := os.ReadFile(corePath)
	if err != nil {
		fmt.Println("Failed to load core.f: ", err)
		return
	}

	machine = vm.NewVM()
	history = make([]string, 0, HISTORY_SIZE)
	historyIndex = 0
	// machine.SetMode(vm.Standard)
	// scanner := bufio.NewScanner(os.Stdin)

	for _, line := range strings.Split(string(core), "\n") {
		if !strings.HasPrefix(line, "\\") {
			words := vm.GetWords(line)
			machine.Process(words, "", 0)
		}
		// fmt.Printf("DEFINE: %v\n", line)
	}

	machine.SetMode(vm.Standard)

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Printf("Type 'bye' to exit\n\n")
	fmt.Print("> ")
	// for scanner.Scan() {
	line := ""
	exit := false
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
		}

		if event.Key == keyboard.KeyCtrlC {
			break
		}

		switch event.Key {
		case keyboard.KeyEnter:
			addToHistory(line)
			// fmt.Printf("\nline: '%v'\nhist: [%v]\n", line, strings.Join(history, ", "))
			words := line + "\n"
			line = ""
			exit = process(words)
			// fmt.Print("exit: %v", exit)

		case keyboard.KeyArrowUp:
			line = gotoHistory(line, -1)
		case keyboard.KeyArrowDown:
			line = gotoHistory(line, 1)
		case keyboard.KeyDelete:
			line = backspace(line)
		case keyboard.KeyBackspace2:
			line = backspace(line)
		case keyboard.KeyBackspace:
			line = backspace(line)
		case keyboard.KeySpace:
			fmt.Print(" ")
			line += " "
		default:
			fmt.Printf("%c", event.Rune)
			line += string(event.Rune)
		}

		// fmt.Print("you pressed: rune: %q, key: %X\n", event.Rune, event.Key)
		if exit {
			break
		}
	}
}

func backspace(line string) string {
	if len(line) < 1 {
		return ""
	}
	fmt.Printf(string(keyboard.KeyBackspace))
	fmt.Printf(" ")
	fmt.Printf(string(keyboard.KeyBackspace))
	line = line[:len(line)-1]
	return line
}

func gotoHistory(line string, i int) string {
	if len(history) < 1 {
		return ""
	}

	historyIndex += i
	if historyIndex < 0 {
		historyIndex = len(history) - 1
	}

	if historyIndex > len(history) {
		historyIndex = 0
	}

	for _, _ = range line {
		line = backspace(line)
	}

	if historyIndex == len(history) {
		line = ""
	} else {
		line = history[historyIndex]
	}

	fmt.Printf("%v", line)
	return line
}

func addToHistory(line string) {
	if len(line) < 1 {
		historyIndex++
		return
	}

	if len(history) > 0 {
		last := history[len(history)-1]

		if last == line {
			// don't duplicate the last used word
			historyIndex++
			return
		}
	}

	if len(history) >= HISTORY_SIZE {
		history = history[:HISTORY_SIZE-1]
	}
	history = append(history, line)
	historyIndex = len(history)
	// fmt.Printf("history: [%v]", strings.Join(history, ", "))
}

func process(input string) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("\n\nERROR\n")
			fmt.Printf(" Indexes: S: %d, R: %v, D: %v\n", machine.Memory.Stack().Index, machine.Memory.ReturnStack().Index, machine.Memory.DataSpace().Index)
			fmt.Printf("   Stack: <%-4d> %v\n", machine.Memory.Stack().Index, machine.Memory.Stack().Values[:32])
			fmt.Printf("  Return: <%-4d> %v\n", machine.Memory.ReturnStack().Index, machine.Memory.ReturnStack().Values[:32])
			fmt.Printf("    Data: <%-4d> %v\n", machine.Memory.DataSpace().Index, machine.Memory.DataSpace().Values[:32])
			machine.Memory.Recover()
			fmt.Printf("\n\npanic occurred: %v\n\n", err)
		}
	}()

	// input := scanner.Text()
	// words := input
	// fmt.Printf("words: %v\n", words)
	fmt.Printf(" ")
	words := vm.GetWords(input)
	_, err := machine.Process(words, "", 0)
	if err != nil {
		if err.Error() == "bye" {
			return true
		}
		fmt.Printf("\nprocess error: %v\n\n", err)
	}
	fmt.Printf(" ok\n> ")
	return false
}
