package ir

import (
	"fmt"
	"github.com/0x51-dev/upeg/parser"
)

func NewInvalidNodeStructureError(name string, n *parser.Node) error {
	return &InvalidNodeStructure{
		Name: name,
		Node: n,
	}
}

type InvalidNodeStructure struct {
	Name string
	Node *parser.Node
}

func (e InvalidNodeStructure) Error() string {
	return fmt.Sprintf("invalid node structure for %q: %v", e.Name, e.Node)
}
