package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/zoul0813/go-forth/pkg/vm"
)

const (
	VERSION = "0.0.0-alpha"
)

func main() {
	fmt.Printf("Go Forth %v, Copyright (C) 2024 David Higgins <zoul0813@me.com>\nCompiled with %v\n", VERSION, runtime.Version())

	cwd, _ := os.Getwd()
	corePath := path.Join(cwd, "cmd/core.f")
	fmt.Printf("Core: %v\n", corePath)
	core, err := os.ReadFile(corePath)
	if err != nil {
		fmt.Println("Failed to load core.f: ", err)
		return
	}

	machine := vm.NewVM()

	scanner := bufio.NewScanner(os.Stdin)

	for _, line := range strings.Split(string(core), "\n") {
		if !strings.HasPrefix(line, "\\") {
			machine.Process(line)
		}
		// fmt.Printf("DEFINE: %v\n", line)
	}

	machine.SetMode(vm.Standard)
	fmt.Printf("Type 'bye' to exit\n\n")
	fmt.Print("> ")
	for scanner.Scan() {
		running := true
		input := scanner.Text()
		words := input

		running, err = machine.Process(words)
		if err != nil {
			fmt.Printf("%v", err)
		}
		if !running {
			break
		}

		fmt.Printf("  ok\n> ")
	}

}
