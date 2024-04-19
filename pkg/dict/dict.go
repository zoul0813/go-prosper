package dict

import (
	"fmt"
	"strings"

	"github.com/zoul0813/go-forth/pkg/definitions"
)

type Dict struct {
	Words map[string]*definitions.Definition
}

func NewDict() *Dict {
	dict :=
		&Dict{
			Words: make(map[string]*definitions.Definition),
		}

	for word, fn := range definitions.Core {
		dict.DefineCore(word, "", fn)
	}

	return dict
}

func (d *Dict) Define(word, definition string, compiled []string) error {
	name := strings.ToLower(word)

	existing, ok := d.Words[name]
	if (ok && !existing.Internal) || !ok {
		d.Words[name] = &definitions.Definition{
			Name:      word,
			Definiton: definition,
			Compiled:  compiled,
			Fn:        nil,
			Internal:  false,
		}
		return nil
	}
	if ok && existing.Internal {
		return fmt.Errorf("can not override internal definition")
	}
	return fmt.Errorf("failed to define see")
}

func (d *Dict) DefineSystem(word string) bool {
	d.Words[strings.ToLower(word)] = &definitions.Definition{
		Name:      word,
		Definiton: "<system>",
		Fn:        nil,
		Internal:  true,
	}
	return true
}

func (d *Dict) DefineCore(word, definition string, fn definitions.DefinitionFn) bool {
	def := &definitions.Definition{
		Name:      word,
		Definiton: "<core>",
		Fn:        fn,
		Internal:  true,
	}

	d.Words[strings.ToLower(word)] = def

	return true
}

func (d *Dict) Lookup(word string) (*definitions.Definition, bool) {
	w, ok := d.Words[strings.ToLower(word)]
	return w, ok
}
