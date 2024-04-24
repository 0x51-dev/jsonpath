package jsonpath

import (
	"fmt"
	"github.com/0x51-dev/jsonpath/internal/ir"
	"regexp"
)

func (ctx *context) checkFunctionExpr(expr *ir.FunctionExpr, node any) error {
	switch name := expr.Name; name {
	case "match", "search":
		if len(expr.Arguments) != 2 {
			return fmt.Errorf("invalid number of arguments for match: %d", len(expr.Arguments))
		}
		s, ok := expr.Arguments[1].(*ir.String)
		if !ok {
			return fmt.Errorf("unsupported argument type for match: %T", expr.Arguments[1])
		}
		v, err := ctx.value(s, node)
		if err != nil {
			return err
		}
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return err
		}
		switch arg := expr.Arguments[0].(type) {
		case *ir.RelQuery:
			v := newContext(node).applyPath(
				&ir.JSONPathQuery{Segments: arg.Segments},
			)
			if v == nil || len(v) != 1 {
				return fmt.Errorf("no matching expression")
			}
			s, ok := v[0].(string)
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
}
