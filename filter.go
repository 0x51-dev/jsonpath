package jsonpath

import (
	"fmt"
	"github.com/0x51-dev/jsonpath/cmp"
	"github.com/0x51-dev/jsonpath/internal/ir"
	"sort"
)

func (ctx *context) applyFilterSelector(selector *ir.FilterSelector, node any, recursive bool) NodeList {
	var nodeList NodeList
	switch node := node.(type) {
	case []any:
		for _, value := range node {
			if err := ctx.checkLogicalExpr(selector.LogicalExpr, value); err == nil {
				nodeList = append(nodeList, value)
			}
		}

		if recursive {
			for _, value := range node {
				nodeList = append(
					nodeList,
					ctx.applyFilterSelector(
						selector,
						value,
						recursive,
					)...,
				)
			}
		}
	case map[string]any:
		var keys []string
		for key := range node {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, k := range keys {
			value := node[k]
			if err := ctx.checkLogicalExpr(selector.LogicalExpr, value); err == nil {
				nodeList = append(nodeList, value)
			}

			if recursive {
				nodeList = append(
					nodeList,
					ctx.applyFilterSelector(
						selector,
						value,
						recursive,
					)...,
				)
			}
		}
	}
	return nodeList
}

func (ctx *context) checkBasicExpr(expr ir.BasicExpr, node any) error {
	switch expr := expr.(type) {
	case *ir.ComparisonExpr:
		left, err := ctx.value(expr.Left, node)
		if err != nil {
			return err
		}
		right, err := ctx.value(expr.Right, node)
		if err != nil {
			return err
		}
		return cmp.Compare(left, right, expr.Op)
	case *ir.ParenExpr:
		return ctx.checkLogicalExpr(expr.LogicalExpr, node)
	case *ir.TestExpr:
		if expr.Negation {
			panic("not implemented: Negation")
		}
		switch expr := expr.TestExpr.(type) {
		case *ir.RelQuery:
			v := newContext(node).applyPath(
				&ir.JSONPathQuery{Segments: expr.Segments},
			)
			if v == nil {
				return fmt.Errorf("no matching expression")
			}
			return nil
		case *ir.JSONPathQuery:
			panic("not implemented: JSONPathQuery")
		case *ir.FunctionExpr:
			return ctx.checkFunctionExpr(expr, node)
		default:
			return fmt.Errorf("unsupported test expression type: %T", expr)
		}
	default:
		return fmt.Errorf("unsupported basic expression type: %T", expr)
	}
}

func (ctx *context) checkLogicalAndExpr(expr *ir.LogicalAndExpr, node any) error {
	for _, e := range expr.Expressions {
		if err := ctx.checkBasicExpr(e, node); err != nil {
			return err
		}
	}
	return nil // All expressions are true.
}

func (ctx *context) checkLogicalExpr(expr *ir.LogicalExpr, node any) error {
	var hit bool
	for _, e := range expr.Expressions {
		if err := ctx.checkLogicalAndExpr(e, node); err != nil {
			continue
		}
		hit = true
	}
	if !hit {
		return fmt.Errorf("no matching expression")
	}
	return nil
}

func (ctx *context) value(comp ir.Comparable, node any) (any, error) {
	switch comp := comp.(type) {
	case *ir.AbsSingularQuery:
		return comp.Value(ctx.root)
	case *ir.RelSingularQuery:
		return comp.Value(node)
	case *ir.FunctionExpr:
		switch comp.Name {
		case "value":
			switch arg := comp.Arguments[0].(type) {
			case *ir.JSONPathQuery:
				nodeList := ctx.applyPath(arg)
				if len(nodeList) != 1 {
					return nil, nil
				}
				return nodeList[0], nil
			case *ir.RelQuery:
				nodeList := newContext(node).applyPath(&ir.JSONPathQuery{
					Segments: arg.Segments,
				})
				if len(nodeList) != 1 {
					return nil, nil
				}
				return nodeList[0], nil
			default:
				panic(fmt.Sprintf("unsupported value arg %T", arg))
			}
		default:
			panic(fmt.Sprintf("unsupported function: %s", comp.Name))
		}
	default:
		return comp.Value(nil)
	}
}
