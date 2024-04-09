package dict

import (
	"strings"
)

type Definition struct {
	Name      string
	Definiton string
	Fn        func()
}

type Dict struct {
	Words map[string]*Definition
}

func NewDict() *Dict {
	return &Dict{
		Words: make(map[string]*Definition),
	}
}

func (d *Dict) Define(word, definition string) bool {
	d.Words[strings.ToLower(word)] = &Definition{
		Name:      word,
		Definiton: definition,
		Fn:        nil,
	}
	return true
}

// func (d *Dict) DefineCore(word, definition string, fn func()) bool {
// 	d.Words[strings.ToLower(word)] = &Word{
// 		Name:      word,
// 		Definiton: "<core>",
// 		Fn:        fn,
// 	}
// 	return true
// }

func (d *Dict) Lookup(word string) (*Definition, bool) {
	w, ok := d.Words[strings.ToLower(word)]
	return w, ok
}
