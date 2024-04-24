package jsonpath

import (
	"fmt"
	"github.com/0x51-dev/jsonpath/internal/grammar"
	"github.com/0x51-dev/jsonpath/internal/ir"
	"github.com/0x51-dev/upeg/parser/op"
)

// NodeList is a list of nodes.
type NodeList []any

// Path represents a JSONPath query.
type Path struct {
	query *ir.JSONPathQuery
}

// New creates a new JSONPath query from the given string.
func New(query string) (*Path, error) {
	p, err := grammar.NewParser([]rune(query))
	if err != nil {
		return nil, err
	}
	n, err := p.Parse(op.And{grammar.JsonpathQuery, op.EOF{}})
	if err != nil {
		return nil, err
	}
	q, err := ir.ParseJSONPathQuery(n)
	if err != nil {
		return nil, err
	}
	return &Path{query: q}, nil
}

// Apply applies the JSONPath query to the given argument.
func (p Path) Apply(queryArgument any) NodeList {
	return newContext(queryArgument).applyPath(p.query)
}

// Query returns the query string.
func (p Path) Query() string {
	return p.query.String()
}

type context struct {
	root any
}

func newContext(root any) *context {
	return &context{root: root}
}

// applyBracketedSelection returns a list of nodes from the given current node.
func (ctx *context) applyBracketedSelection(segment *ir.BracketedSelection, node any, recursive bool) NodeList {
	var nodeList NodeList
	for _, selector := range segment.Selectors {
		nodeList = append(
			nodeList,
			ctx.applySelector(selector, node, recursive)...,
		)
	}
	return nodeList
}

// applyPath returns a list of nodes from the given input, by applying the path segments.
func (ctx *context) applyPath(p *ir.JSONPathQuery) NodeList {
	nodeList := NodeList{ctx.root}
	for _, segment := range p.Segments {
		if len(nodeList) == 0 {
			return nil
		}

		switch segment := segment.(type) {
		case ir.ChildSegment:
			nodeList = ctx.applySegment(segment, nodeList, false)
		case *ir.DescendantSegment:
			nodeList = ctx.applySegment(segment.Segment, nodeList, true)
		default:
			panic(fmt.Sprintf("unsupported segment type: %T", segment))
		}
	}
	return nodeList
}

// applySegment returns a list of nodes from the given input, by applying the segment.
func (ctx *context) applySegment(segment ir.Segment, input NodeList, recursive bool) NodeList {
	var nodeList NodeList
	switch segment := segment.(type) {
	case *ir.BracketedSelection:
		for _, node := range input {
			nodeList = append(
				nodeList,
				ctx.applyBracketedSelection(
					segment,
					node,
					recursive,
				)...,
			)
		}
	case *ir.WildcardSelector:
		for _, node := range input {
			nodeList = append(
				nodeList,
				applyWildcardSelector(
					node,
					recursive,
				)...,
			)
		}
	case *ir.MemberNameShorthand:
		for _, node := range input {
			nodeList = append(
				nodeList,
				applyNameSelector(
					&ir.NameSelector{Name: segment.Name},
					node,
					recursive,
				)...,
			)
		}
	default:
		panic(fmt.Sprintf("unsupported segment type: %T", segment))
	}
	return nodeList
}
