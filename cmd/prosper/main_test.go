package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eiannone/keyboard"
	"github.com/zoul0813/go-prosper/pkg/vm"
)

func TestReadKey(t *testing.T) {
	events := make(chan keyboard.KeyEvent, 5)
	events <- keyboard.KeyEvent{Key: keyboard.KeyArrowUp}
	events <- keyboard.KeyEvent{Rune: 'w'}
	events <- keyboard.KeyEvent{Key: keyboard.KeySpace}
	events <- keyboard.KeyEvent{Key: keyboard.KeyEnter}
	events <- keyboard.KeyEvent{Key: keyboard.KeyCtrlC}
	close(events)
	for _, want := range []rune{'w', ' ', '\n'} {
		got, err := readKey(events)
		if err != nil || got != want {
			t.Fatalf("key = %q, %v; want %q", got, err, want)
		}
	}
	if _, err := readKey(events); !errors.Is(err, vm.ErrBye) {
		t.Fatalf("Ctrl-C error = %v", err)
	}
	if _, err := readKey(events); err != io.EOF {
		t.Fatalf("closed input error = %v", err)
	}
	failure := errors.New("terminal disconnected")
	events = make(chan keyboard.KeyEvent, 1)
	events <- keyboard.KeyEvent{Err: failure}
	if _, err := readKey(events); !errors.Is(err, failure) {
		t.Fatalf("terminal error = %v", err)
	}
}

func TestLoadFile(t *testing.T) {
	m := vm.NewVM()
	if err := loadFile(m, "../../runtime/core.f"); err != nil {
		t.Fatal(err)
	}
	if err := loadFile(m, "../../2048.fs"); err != nil {
		t.Fatal(err)
	}
	if _, ok := m.Memory.Dict().Lookup("play"); !ok {
		t.Fatal("game not loaded")
	}
	path := filepath.Join(t.TempDir(), "bad.fs")
	if err := os.WriteFile(path, []byte(": unfinished 1"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := loadFile(m, path); err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("missing source context: %v", err)
	}
	if err := m.Execute("1 2 +"); err != nil {
		t.Fatal(err)
	}
	if m.Memory.Stack().Pop() != 3 {
		t.Fatal("failed file left compiler active")
	}
	if err := loadFile(m, filepath.Join(t.TempDir(), "missing.fs")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing file error = %v", err)
	}
}
