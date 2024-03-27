package jsonpath

import (
	"fmt"
	"github.com/0x51-dev/jsonpath/internal/grammar"
	"github.com/0x51-dev/jsonpath/internal/ir"
	"github.com/0x51-dev/upeg/parser/op"
	"regexp"
	"sort"
)

type Path struct {
	query *ir.JSONPathQuery
}

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
func (p Path) Apply(queryArgument any) (any, error) {
	return newContext(queryArgument).applyPath(p.query)
}

func (p Path) Query() string {
	return p.query.String()
}

type context struct {
	root any
}

func newContext(root any) *context {
	return &context{root: root}
}

func (ctx *context) applyChildSegment(segment ir.ChildSegment, current any) (any, error) {
	switch segment := segment.(type) {
	case *ir.BracketedSelection:
		var nodeList []any

		start := current
		for _, selector := range segment.Selectors {
			var err error
			if current, err = ctx.applySelector(selector, start); err != nil {
				return nil, err
			}
			if current != nil {
				// Only append non-empty nodes.
				if n, ok := current.([]any); !ok || len(n) != 0 {
					nodeList = append(nodeList, current)
				}
			}
			start = current
		}

		if len(nodeList) == 0 {
			return nil, fmt.Errorf("no matching expression")
		}

		if len(nodeList) == 1 {
			return nodeList[0], nil
		}
		return nodeList, nil
	case *ir.WildcardSelector:
		switch c := current.(type) {
		case map[string]any:
			return ctx.applyWildcardSelector(current)
		case []any:
			var hit bool
			for _, value := range c {
				if value, err := ctx.applyChildSegment(segment, value); err == nil {
					hit = true
					current = value
				}
			}
			if hit {
				return current, nil
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("unsupported child segment type: %T", segment)
		}
	case *ir.MemberNameShorthand:
		switch current := current.(type) {
		case map[string]any:
			if value, ok := current[segment.Name]; ok {
				return value, nil
			}
			return nil, nil
		default:
			return nil, fmt.Errorf("unsupported child segment type: %T", segment)
		}
	default:
		return nil, fmt.Errorf("unsupported child segment type: %T", segment)
	}
}

func (ctx *context) applyFilterSelector(selector *ir.FilterSelector, current any) (any, error) {
	switch c := current.(type) {
	case []any:

		var nodeList []any
		for _, value := range c {
			if err := ctx.checkLogicalExpr(selector.LogicalExpr, value); err != nil {
				continue
			}
			nodeList = append(nodeList, value)
		}
		return nodeList, nil
	case map[string]any:
		var keys []string
		for key := range c {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		var nodeList []any
		for _, k := range keys {
			current = c[k]
			if err := ctx.checkLogicalExpr(selector.LogicalExpr, current); err != nil {
				continue
			}
			nodeList = append(nodeList, current)
		}
		return nodeList, nil
	default:
		return nil, nil
	}
}

func (ctx *context) applyPath(p *ir.JSONPathQuery) (any, error) {
	current := ctx.root
	for _, segment := range p.Segments {
		if childSegment, ok := segment.(ir.ChildSegment); ok {
			var err error
			if current, err = ctx.applyChildSegment(childSegment, current); err != nil {
				return nil, err
			}
			continue
		}
		return nil, fmt.Errorf("unsupported segment type: %T", segment)
	}
	return current, nil
}

func (ctx *context) applySelector(selector ir.Selector, current any) (any, error) {
	switch selector := selector.(type) {
	case *ir.NameSelector:
		switch c := current.(type) {
		case map[string]any:
			if value, ok := c[selector.Name]; ok {
				return value, nil
			}
		}
		return nil, nil
	case *ir.WildcardSelector:
		return ctx.applyWildcardSelector(current)
	case *ir.SliceSelector:
		switch c := current.(type) {
		case []any:
			var slice []any
			if 0 < selector.Step {
				if selector.End == -1 {
					selector.End = len(c)
				}
				for i := max(selector.Start, 0); i < min(selector.End, len(c)); i += selector.Step {
					slice = append(slice, c[i])
				}
			} else {
				if selector.Start == 0 && selector.End == -1 {
					selector.Start = len(c) - 1
				}
				for i := min(selector.Start, len(c)); max(selector.End, -1) < i; i += selector.Step {
					slice = append(slice, c[i])
				}
			}
			return slice, nil
		}
		return nil, nil
	case *ir.IndexSelector:
		switch c := current.(type) {
		case []any:
			idx := selector.Index
			if idx < 0 {
				idx += len(c)
			}
			return c[idx], nil
		}
		return nil, nil
	case *ir.FilterSelector:
		return ctx.applyFilterSelector(selector, current)
	default:
		return nil, fmt.Errorf("unsupported selector type: %T", selector)
	}
}

func (ctx *context) applyWildcardSelector(current any) (any, error) {
	switch c := current.(type) {
	case map[string]any:
		var keyList []string
		for key := range c {
			keyList = append(keyList, key)
		}
		sort.Strings(keyList)

		var nodeList []any
		for _, key := range keyList {
			nodeList = append(nodeList, c[key])
		}
		return nodeList, nil
	default:
		return current, nil
	}
}

func (ctx *context) checkBasicExpr(expr ir.BasicExpr, current any) error {
	switch expr := expr.(type) {
	case *ir.ComparisonExpr:
		left, err := ctx.value(expr.Left, current)
		if err != nil {
			return err
		}
		right, err := ctx.value(expr.Right, current)
		if err != nil {
			return err
		}
		return ctx.compare(left, right, expr.Op)
	case *ir.ParenExpr:
		return ctx.checkLogicalExpr(expr.LogicalExpr, current)
	case *ir.TestExpr:
		if expr.Negation {
			panic("not implemented")
		}
		switch expr := expr.TestExpr.(type) {
		case *ir.RelQuery:
			v, err := newContext(current).applyPath(&ir.JSONPathQuery{Segments: expr.Segments})
			if err != nil {
				return err
			}
			if v == nil {
				return fmt.Errorf("no matching expression")
			}
			return nil
		case *ir.JSONPathQuery:
			panic("not implemented: JSONPathQuery")
		case *ir.FunctionExpr:
			switch name := expr.Name; name {
			case "match", "search":
				if len(expr.Arguments) != 2 {
					return fmt.Errorf("invalid number of arguments for match: %d", len(expr.Arguments))
				}
				s, ok := expr.Arguments[1].(*ir.StringLiteral)
				if !ok {
					return fmt.Errorf("unsupported argument type for match: %T", expr.Arguments[1])
				}
				v, _ := ctx.value(s, current)
				r, err := regexp.Compile(v.(string))
				if err != nil {
					return err
				}
				switch arg := expr.Arguments[0].(type) { //  TODO: extract argument type
				case *ir.RelQuery:
					v, err := newContext(current).applyPath(&ir.JSONPathQuery{Segments: arg.Segments})
					if err != nil {
						return err
					}
					if v == nil {
						return fmt.Errorf("no matching expression")
					}
					s, ok := v.(string)
					if !ok {
						return fmt.Errorf("unsupported argument type for match: %T", v)
					}
					if name == "match" && r.FindString(s) != s {
						return fmt.Errorf("no matching expression")
					}
					if name == "search" && !r.MatchString(s) {
						return fmt.Errorf("no matching expression")
					}
					return nil
				default:
					panic("not implemented: match argument type")
				}
			default:
				panic("not implemented: function name")
			}
		default:
			return fmt.Errorf("unsupported test expression type: %T", expr)
		}
	default:
		return fmt.Errorf("unsupported basic expression type: %T", expr)
	}
}

func (ctx *context) checkLogicalAndExpr(expr *ir.LogicalAndExpr, current any) error {
	for _, e := range expr.Expressions {
		if err := ctx.checkBasicExpr(e, current); err != nil {
			return err
		}
	}
	return nil // All expressions are true.
}

func (ctx *context) checkLogicalExpr(expr *ir.LogicalExpr, current any) error {
	var hit bool
	for _, e := range expr.Expressions {
		if err := ctx.checkLogicalAndExpr(e, current); err != nil {
			continue
		}
		hit = true
	}
	if !hit {
		return fmt.Errorf("no matching expression")
	}
	return nil
}

func (ctx *context) value(cmp ir.Comparable, current any) (any, error) {
	switch cmp := cmp.(type) {
	case *ir.AbsSingularQuery:
		return cmp.Value(ctx.root)
	case *ir.RelSingularQuery:
		return cmp.Value(current)
	default:
		return cmp.Value(nil)
	}
}
