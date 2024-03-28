package cmp

import "errors"

func Compare(a, b any, op string) error {
	switch op {
	case "==":
		return eq(a, b)
	case "!=":
		if err := eq(a, b); err == nil {
			return NewEqualError(a, b)
		}
		return nil
	case "<=":
		err := gt(a, b)
		if err == nil {
			return NewNotLesserThanOrEqualError(a, b)
		}
		var typeNotSupported *TypeNotSupportedError
		if errors.As(err, &typeNotSupported) {
			return eq(a, b)
		}
		var notGreaterThan *NotGreaterThanError
		if !errors.As(err, &notGreaterThan) {
			return err
		}
		return nil
	case ">=":
		err := lt(a, b)
		if err == nil {
			return NewNotGreaterThanOrEqualError(a, b)
		}
		var typeNotSupported *TypeNotSupportedError
		if errors.As(err, &typeNotSupported) {
			return eq(a, b)
		}
		var notLesserThan *notLesserThanError
		if !errors.As(err, &notLesserThan) {
			return err
		}
		return nil
	case "<":
		return lt(a, b)
	case ">":
		return gt(a, b)
	default:
		return NewOperatorNotSupportedError(op)
	}
}

func eq(a, b any) error {
	if a == nil && b == nil {
		return nil
	}
	if a == nil || b == nil {
		return NewNotEqualError(a, b)
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
				return NewTypeMismatchError(a, b)
			}
			if float64(a) != f {
				return NewNotEqualError(a, b)
			}
			return nil
		}
		if a != i {
			return NewNotEqualError(a, b)
		}
		return nil
	case float64:
		f, ok := b.(float64)
		if !ok {
			i, ok := b.(int64)
			if !ok {
				return NewTypeMismatchError(a, b)
			}
			f = float64(i)
		}
		if a != f {
			return NewNotEqualError(a, b)
		}
		return nil
	case string:
		s, ok := b.(string)
		if !ok {
			return NewTypeMismatchError(a, b)
		}
		if a != s {
			return NewNotEqualError(a, b)
		}
		return nil
	case bool:
		b, ok := b.(bool)
		if !ok {
			return NewTypeMismatchError(a, b)
		}
		if a != b {
			return NewNotEqualError(a, b)
		}
		return nil
	case []any:
		b, ok := b.([]any)
		if !ok {
			return NewTypeMismatchError(a, b)
		}
		if len(a) != len(b) {
			return NewNotEqualError(a, b)
		}
		for i := range a {
			if err := eq(a[i], b[i]); err != nil {
				return err
			}
		}
		return nil
	case map[string]any:
		b, ok := b.(map[string]any)
		if !ok {
			return NewTypeMismatchError(a, b)
		}
		if len(a) != len(b) {
			return NewNotEqualError(a, b)
		}
		for key, value := range a {
			if err := eq(value, b[key]); err != nil {
				return err
			}
		}
		return nil
	default:
		return NewTypeNotSupportedError(a)
	}
}

func gt(a, b any) error {
	if a == nil && b == nil {
		return NewNotGreaterThanError(a, b)
	}
	if a == nil || b == nil {
		return NewTypeMismatchError(a, b)
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
				return NewTypeMismatchError(a, b)
			}
			return nil
		}
		if a <= i {
			return NewNotGreaterThanError(a, b)
		}
		return nil
	case float64:
		f, ok := b.(float64)
		if !ok {
			i, ok := b.(int64)
			if !ok {
				return NewTypeMismatchError(a, b)
			}
			f = float64(i)
		}
		if a <= f {
			return NewNotGreaterThanError(a, b)
		}
		return nil
	case string:
		s, ok := b.(string)
		if !ok {
			return NewTypeMismatchError(a, b)
		}
		if a <= s {
			return NewNotGreaterThanError(a, b)
		}
		return nil
	default:
		return NewTypeNotSupportedError(a)
	}
}

func lt(a, b any) error {
	if a == nil && b == nil {
		return NewNotLesserThanError(a, b)
	}
	if a == nil || b == nil {
		return NewTypeMismatchError(a, b)
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
				return NewTypeMismatchError(a, b)
			}
			return nil
		}
		if a >= i {
			return NewNotLesserThanError(a, b)
		}
		return nil
	case float64:
		f, ok := b.(float64)
		if !ok {
			i, ok := b.(int64)
			if !ok {
				return NewTypeMismatchError(a, b)
			}
			f = float64(i)
		}
		if a >= f {
			return NewNotLesserThanError(a, b)
		}
		return nil
	case string:
		s, ok := b.(string)
		if !ok {
			return NewTypeMismatchError(a, b)
		}
		if a >= s {
			return NewNotLesserThanError(a, b)
		}
		return nil
	default:
		return NewTypeNotSupportedError(a)
	}
}
