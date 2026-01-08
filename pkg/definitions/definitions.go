package definitions

import (
	"fmt"
	"strings"

	"github.com/zoul0813/go-prosper/pkg/stack"
)

type memory interface {
	Stack() *stack.Stack
	ReturnStack() *stack.Stack
	DataSpace() *stack.Stack
	NumberBase() int
}
type Processor interface {
	Process(string) (bool, error)
}

type DefinitionFn func(m memory) (int, error)

type Definition struct {
	Name      string
	Definiton string
	Compiled  []string
	Fn        DefinitionFn
	Internal  bool
}

func (d *Definition) String() string {
	s := fmt.Sprintf(": %v\n  %v ;\n\ncompiled:\n  %v", d.Name, d.Definiton, strings.Join(d.Compiled, " "))
	if d.Internal {
		s += " (internal)"
	}
	return s
}
