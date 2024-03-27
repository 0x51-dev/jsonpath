package jsonpath

import (
	"errors"
	"fmt"
)

func (ctx *context) compare(a, b any, op string) error {
	switch op {
	case "==":
		return ctx.eq(a, b)
	case "!=":
		if err := ctx.eq(a, b); err == nil {
			return newEqualError(a, b)
		}
		return nil
	case "<=":
		err := ctx.gt(a, b)
		if err == nil {
			return newNotLesserThanOrEqualError(a, b)
		}
		var typeNotSupported *typeNotSupportedError
		if errors.As(err, &typeNotSupported) {
			return ctx.eq(a, b)
		}
		var notGreaterThan *notGreaterThanError
		if !errors.As(err, &notGreaterThan) {
			return err
		}
		return nil
	case ">=":
		err := ctx.lt(a, b)
		if err == nil {
			return newNotGreaterThanOrEqualError(a, b)
		}
		var typeNotSupported *typeNotSupportedError
		if errors.As(err, &typeNotSupported) {
			return ctx.eq(a, b)
		}
		var notLesserThan *notLesserThanError
		if !errors.As(err, &notLesserThan) {
			return err
		}
		return nil
	case "<":
		return ctx.lt(a, b)
	case ">":
		return ctx.gt(a, b)
	default:
		return newOperatorNotSupportedError(op)
	}
}

func (ctx *context) eq(a, b any) error {
	if a == nil && b == nil {
		return nil
	}
	if a == nil || b == nil {
		return newNotEqualError(a, b)
	}

	// Support for int to int64 conversion.
	if i, ok := a.(int); ok {
		a = int64(i)
	}
	if i, ok := b.(int); ok {
		b = int64(i)
	}

	switch a := a.(type) {
	case int64:
		i, ok := b.(int64)
		if !ok {
			f, ok := b.(float64)
			if !ok {
				return newTypeMismatchError(a, b)
			}
			if float64(a) != f {
				return newNotEqualError(a, b)
			}
			return nil
		}
		if a != i {
			return newNotEqualError(a, b)
		}
		return nil
	case float64:
		f, ok := b.(float64)
		if !ok {
			i, ok := b.(int64)
			if !ok {
				return newTypeMismatchError(a, b)
			}
			f = float64(i)
		}
		if a != f {
			return newNotEqualError(a, b)
		}
		return nil
	case string:
		s, ok := b.(string)
		if !ok {
			return newTypeMismatchError(a, b)
		}
		if a != s {
			return newNotEqualError(a, b)
		}
		return nil
	case bool:
		b, ok := b.(bool)
		if !ok {
			return newTypeMismatchError(a, b)
		}
		if a != b {
			return newNotEqualError(a, b)
		}
		return nil
	case []any:
		b, ok := b.([]any)
		if !ok {
			return newTypeMismatchError(a, b)
		}
		if len(a) != len(b) {
			return newNotEqualError(a, b)
		}
		for i := range a {
			if err := ctx.eq(a[i], b[i]); err != nil {
				return err
			}
		}
		return nil
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
			return newTypeMismatchError(a, b)
		}
		if len(a) != len(b) {
			return newNotEqualError(a, b)
		}
		for key, value := range a {
			if err := ctx.eq(value, b[key]); err != nil {
				return err
			}
		}
		return nil
	default:
		return newTypeNotSupportedError(a)
	}
}

func (ctx *context) gt(a, b any) error {
	if a == nil && b == nil {
		return newNotGreaterThanError(a, b)
	}
	if a == nil || b == nil {
		return newTypeMismatchError(a, b)
	}

	// Support for int to int64 conversion.
	if i, ok := a.(int); ok {
		a = int64(i)
	}
	if i, ok := b.(int); ok {
		b = int64(i)
	}

	switch a := a.(type) {
	case int64:
		i, ok := b.(int64)
		if !ok {
			f, ok := b.(float64)
			if !ok || float64(a) <= f {
				return newTypeMismatchError(a, b)
			}
			return nil
		}
		if a <= i {
			return newNotGreaterThanError(a, b)
		}
		return nil
	case float64:
		f, ok := b.(float64)
		if !ok {
			i, ok := b.(int64)
			if !ok {
				return newTypeMismatchError(a, b)
			}
			f = float64(i)
		}
		if a <= f {
			return newNotGreaterThanError(a, b)
		}
		return nil
	case string:
		s, ok := b.(string)
		if !ok {
			return newTypeMismatchError(a, b)
		}
		if a <= s {
			return newNotGreaterThanError(a, b)
		}
		return nil
	default:
		return newTypeNotSupportedError(a)
	}
}

func (ctx *context) lt(a, b any) error {
	if a == nil && b == nil {
		return newNotLesserThanError(a, b)
	}
	if a == nil || b == nil {
		return newTypeMismatchError(a, b)
	}

	// Support for int to int64 conversion.
	if i, ok := a.(int); ok {
		a = int64(i)
	}
	if i, ok := b.(int); ok {
		b = int64(i)
	}

	switch a := a.(type) {
	case int64:
		i, ok := b.(int64)
		if !ok {
			f, ok := b.(float64)
			if !ok || float64(a) >= f {
				return newTypeMismatchError(a, b)
			}
			return nil
		}
		if a >= i {
			return newNotLesserThanError(a, b)
		}
		return nil
	case float64:
		f, ok := b.(float64)
		if !ok {
			i, ok := b.(int64)
			if !ok {
				return newTypeMismatchError(a, b)
			}
			f = float64(i)
		}
		if a >= f {
			return newNotLesserThanError(a, b)
		}
		return nil
	case string:
		s, ok := b.(string)
		if !ok {
			return newTypeMismatchError(a, b)
		}
		if a >= s {
			return newNotLesserThanError(a, b)
		}
		return nil
	default:
		return newTypeNotSupportedError(a)
	}
}

// equalError is returned when two values are equal, but should not be.
type equalError struct {
	Left  any
	Right any
}

// newEqualError creates a new equalError.
func newEqualError(left, right any) *equalError {
	return &equalError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *equalError) Error() string {
	return fmt.Sprintf("equal: %v == %v", e.Left, e.Right)
}

// NotEqualError is returned when two values are not equal, but should be.
// For example, when comparing 1 to 2.
type notEqualError struct {
	Left  any
	Right any
}

// newNotEqualError creates a new notEqualError.
func newNotEqualError(left, right any) *notEqualError {
	return &notEqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *notEqualError) Error() string {
	return fmt.Sprintf("not equal: %v != %v", e.Left, e.Right)
}

// notGreaterThanError is returned when a value is not greater than another.
type notGreaterThanError struct {
	Left  any
	Right any
}

// newNotGreaterThanError creates a new notGreaterThanError.
func newNotGreaterThanError(left, right any) *notGreaterThanError {
	return &notGreaterThanError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *notGreaterThanError) Error() string {
	return fmt.Sprintf("not greater than: %v >= %v", e.Left, e.Right)
}

// notGreaterThanOrEqualError is returned when a value is not greater than or equal to another.
type notGreaterThanOrEqualError struct {
	Left  any
	Right any
}

// newNotGreaterThanOrEqualError creates a new notGreaterThanOrEqualError.
func newNotGreaterThanOrEqualError(left, right any) *notGreaterThanOrEqualError {
	return &notGreaterThanOrEqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *notGreaterThanOrEqualError) Error() string {
	return fmt.Sprintf("not greater than or equal: %v < %v", e.Left, e.Right)
}

// notLesserThanError is returned when a value is not lesser than another.
type notLesserThanError struct {
	Left  any
	Right any
}

// newNotLesserThanError creates a new notLesserThanError.
func newNotLesserThanError(left, right any) *notLesserThanError {
	return &notLesserThanError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *notLesserThanError) Error() string {
	return fmt.Sprintf("not lesser than: %v <= %v", e.Left, e.Right)
}

// notLesserThanOrEqualError is returned when a value is not lesser than or equal to another.
type notLesserThanOrEqualError struct {
	Left  any
	Right any
}

// newNotLesserThanOrEqualError creates a new notLesserThanOrEqualError.
func newNotLesserThanOrEqualError(left, right any) *notLesserThanOrEqualError {
	return &notLesserThanOrEqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *notLesserThanOrEqualError) Error() string {
	return fmt.Sprintf("not lesser than or equal: %v > %v", e.Left, e.Right)
}

// operatorNotSupportedError is returned when an operator is not supported.
type operatorNotSupportedError struct {
	Operator string
}

// newOperatorNotSupportedError creates a new operatorNotSupportedError.
func newOperatorNotSupportedError(op string) *operatorNotSupportedError {
	return &operatorNotSupportedError{
		Operator: op,
	}
}

// Error returns the error message.
func (e *operatorNotSupportedError) Error() string {
	return fmt.Sprintf("operator not supported: %s", e.Operator)
}

// TypeMismatchError is returned when a type mismatch is detected.
// For example, when comparing a string to an integer.
type typeMismatchError struct {
	Expected any
	Actual   any
}

// newTypeMismatchError creates a new typeMismatchError.
func newTypeMismatchError(expected, actual any) *typeMismatchError {
	return &typeMismatchError{
		Expected: expected,
		Actual:   actual,
	}
}

// Error returns the error message.
func (e *typeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch: expected %T, got %T", e.Expected, e.Actual)
}

// typeNotSupportedError is returned when a type is not supported.
type typeNotSupportedError struct {
	Type any
}

// newTypeNotSupportedError creates a new typeNotSupportedError.
func newTypeNotSupportedError(t any) *typeNotSupportedError {
	return &typeNotSupportedError{
		Type: t,
	}
}

// Error returns the error message.
func (e *typeNotSupportedError) Error() string {
	return fmt.Sprintf("type not supported: %T", e.Type)
}
