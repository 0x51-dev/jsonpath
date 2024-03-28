package cmp

import "fmt"

// EqualError is returned when two values are equal, but should not be.
type EqualError struct {
	Left  any
	Right any
}

// NewEqualError creates a new EqualError.
func NewEqualError(left, right any) *EqualError {
	return &EqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *EqualError) Error() string {
	return fmt.Sprintf("equal: %v == %v", e.Left, e.Right)
}

// NotEqualError is returned when two values are not equal, but should be.
// For example, when comparing 1 to 2.
type NotEqualError struct {
	Left  any
	Right any
}

// NewNotEqualError creates a new NotEqualError.
func NewNotEqualError(left, right any) *NotEqualError {
	return &NotEqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *NotEqualError) Error() string {
	return fmt.Sprintf("not equal: %v != %v", e.Left, e.Right)
}

// NotGreaterThanError is returned when a value is not greater than another.
type NotGreaterThanError struct {
	Left  any
	Right any
}

// NewNotGreaterThanError creates a new NotGreaterThanError.
func NewNotGreaterThanError(left, right any) *NotGreaterThanError {
	return &NotGreaterThanError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *NotGreaterThanError) Error() string {
	return fmt.Sprintf("not greater than: %v >= %v", e.Left, e.Right)
}

// NotGreaterThanOrEqualError is returned when a value is not greater than or equal to another.
type NotGreaterThanOrEqualError struct {
	Left  any
	Right any
}

// NewNotGreaterThanOrEqualError creates a new NotGreaterThanOrEqualError.
func NewNotGreaterThanOrEqualError(left, right any) *NotGreaterThanOrEqualError {
	return &NotGreaterThanOrEqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *NotGreaterThanOrEqualError) Error() string {
	return fmt.Sprintf("not greater than or equal: %v < %v", e.Left, e.Right)
}

// NotLesserThanOrEqualError is returned when a value is not lesser than or equal to another.
type NotLesserThanOrEqualError struct {
	Left  any
	Right any
}

// NewNotLesserThanOrEqualError creates a new NotLesserThanOrEqualError.
func NewNotLesserThanOrEqualError(left, right any) *NotLesserThanOrEqualError {
	return &NotLesserThanOrEqualError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *NotLesserThanOrEqualError) Error() string {
	return fmt.Sprintf("not lesser than or equal: %v > %v", e.Left, e.Right)
}

// OperatorNotSupportedError is returned when an operator is not supported.
type OperatorNotSupportedError struct {
	Operator string
}

// NewOperatorNotSupportedError creates a new OperatorNotSupportedError.
func NewOperatorNotSupportedError(op string) *OperatorNotSupportedError {
	return &OperatorNotSupportedError{
		Operator: op,
	}
}

// Error returns the error message.
func (e *OperatorNotSupportedError) Error() string {
	return fmt.Sprintf("operator not supported: %s", e.Operator)
}

// TypeMismatchError is returned when a type mismatch is detected.
// For example, when comparing a string to an integer.
type TypeMismatchError struct {
	Expected any
	Actual   any
}

// NewTypeMismatchError creates a new TypeMismatchError.
func NewTypeMismatchError(expected, actual any) *TypeMismatchError {
	return &TypeMismatchError{
		Expected: expected,
		Actual:   actual,
	}
}

// Error returns the error message.
func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch: expected %T, got %T", e.Expected, e.Actual)
}

// TypeNotSupportedError is returned when a type is not supported.
type TypeNotSupportedError struct {
	Type any
}

// NewTypeNotSupportedError creates a new TypeNotSupportedError.
func NewTypeNotSupportedError(t any) *TypeNotSupportedError {
	return &TypeNotSupportedError{
		Type: t,
	}
}

// Error returns the error message.
func (e *TypeNotSupportedError) Error() string {
	return fmt.Sprintf("type not supported: %T", e.Type)
}

// notLesserThanError is returned when a value is not lesser than another.
type notLesserThanError struct {
	Left  any
	Right any
}

// NewNotLesserThanError creates a new notLesserThanError.
func NewNotLesserThanError(left, right any) *notLesserThanError {
	return &notLesserThanError{
		Left:  left,
		Right: right,
	}
}

// Error returns the error message.
func (e *notLesserThanError) Error() string {
	return fmt.Sprintf("not lesser than: %v <= %v", e.Left, e.Right)
}
