package dict

import (
	"strings"

	"github.com/zoul0813/go-forth/pkg/definitions"
)

type Dict struct {
	Words map[string]*definitions.Definition
}

func NewDict() *Dict {
	return &Dict{
		Words: make(map[string]*definitions.Definition),
	}
}

func (d *Dict) Define(word, definition string) bool {
	d.Words[strings.ToLower(word)] = &definitions.Definition{
		Name:      word,
		Definiton: definition,
		Fn:        nil,
	}
	return true
}

func (d *Dict) DefineCore(word, definition string, fn func()) bool {
	d.Words[strings.ToLower(word)] = &definitions.Definition{
		Name:      word,
		Definiton: "<core>",
		Fn:        fn,
	}
	return true
}

func (d *Dict) Lookup(word string) (*definitions.Definition, bool) {
	w, ok := d.Words[strings.ToLower(word)]
	return w, ok
}
