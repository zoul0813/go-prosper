package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/eiannone/keyboard"
	"github.com/zoul0813/go-prosper/pkg/vm"
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
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("Go Forth %v, Copyright (C) 2024 David Higgins <zoul0813@me.com>\nCompiled with %v\n", VERSION, runtime.Version())

	cwd, _ := os.Getwd()
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	var corePath string
	corePaths := []string{
		filepath.Join(exeDir, "runtime", "core.f"),
		filepath.Join(exeDir, "..", "share", "prosper", "runtime", "core.f"),
		filepath.Join(cwd, "runtime", "core.f"),
	}

	for _, p := range corePaths {
		if _, err := os.Stat(p); err == nil {
			corePath = p
			break
		}
	}

	if corePath == "" {
		fmt.Printf("Failed to find core.f, Tried:\n")
		for _, p := range corePaths {
			fmt.Printf(". - %s\n", p)
		}
		return fmt.Errorf("core.f not found")
	}

	corePathRel, _ := filepath.Rel(cwd, corePath)
	fmt.Printf("Core: %v\n", corePathRel)
	machine = vm.NewVM()
	history = make([]string, 0, HISTORY_SIZE)
	historyIndex = 0
	if err := loadFile(machine, corePath); err != nil {
		return err
	}
	for _, path := range os.Args[1:] {
		if err := loadFile(machine, path); err != nil {
			return err
		}
	}
	machine.SetMode(vm.Standard)

	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		return err
	}
	defer func() {
		_ = keyboard.Close()
	}()

	machine.ReadKey = func() (rune, error) { return readKey(keysEvents) }

	fmt.Printf("Type 'bye' to exit\n\n")
	fmt.Print("> ")
	// for scanner.Scan() {
	line := ""
	exit := false
	for {
		event, ok := <-keysEvents
		if !ok {
			return io.EOF
		}
		if event.Err != nil {
			return event.Err
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
	return nil
}

func backspace(line string) string {
	if len(line) < 1 {
		return ""
	}
	fmt.Print("\b")
	fmt.Printf(" ")
	fmt.Print("\b")
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
	fmt.Print(" ")
	if err := machine.Execute(input); err != nil {
		if errors.Is(err, vm.ErrBye) {
			return true
		}
		fmt.Printf("\nprocess error: %v\n> ", err)
		return false
	}
	fmt.Printf(" ok\n> ")
	return false
}

func loadFile(machine *vm.VM, path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := machine.Load(string(source)); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func readKey(events <-chan keyboard.KeyEvent) (rune, error) {
	for event := range events {
		if event.Err != nil {
			return 0, event.Err
		}
		switch event.Key {
		case keyboard.KeyCtrlC:
			return 0, vm.ErrBye
		case keyboard.KeySpace:
			return ' ', nil
		case keyboard.KeyEnter:
			return '\n', nil
		}
		if event.Rune != 0 {
			return event.Rune, nil
		}
	}
	return 0, io.EOF
}
