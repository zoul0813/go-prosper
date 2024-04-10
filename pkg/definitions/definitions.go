package definitions

import (
	"fmt"

	"github.com/zoul0813/go-forth/pkg/stack"
)

type memory interface {
	Stack() *stack.Stack
	ReturnStack() *stack.Stack
	DataSpace() *stack.Stack
}
type Processor interface {
	Process(string) (bool, error)
}

type DefinitionFn func(m memory) (bool, error)

type Definition struct {
	Name      string
	Definiton string
	Fn        DefinitionFn
	Internal  bool
}

func (d *Definition) String() string {
	s := fmt.Sprintf(": %v\n  %v ;", d.Name, d.Definiton)
	if d.Internal {
		s += " (internal)"
	}
	return s
}
