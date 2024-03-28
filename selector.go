package jsonpath

import (
	"fmt"
	"github.com/0x51-dev/jsonpath/internal/ir"
	"sort"
)

// applyIndexSelector returns a list of nodes from the given current node.
// If the current node is a list, it returns the element at the index.
// Otherwise, it returns nil.
func applyIndexSelector(selector *ir.IndexSelector, node any, recursive bool) NodeList {
	var nodeList NodeList
	if node, ok := node.([]any); ok {
		idx := selector.Index
		if idx < 0 {
			// A negative index-selector counts from the array end backwards, obtaining an equivalent non-negative
			// index-selector by adding the length of the array to the negative index.
			idx += len(node)
		}
		if len(node) <= idx {
			// Nothing is selected, and it is not an error, if the index lies outside the range of the array.
			return nil
		}
		nodeList = append(nodeList, node[idx])

		if recursive {
			for _, value := range node {
				if value := applyIndexSelector(selector, value, recursive); value != nil {
					nodeList = append(nodeList, value...)
				}
			}
		}
	}
	if node, ok := node.(map[string]any); ok && recursive {
		var keys []string
		for key := range node {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			if value := applyIndexSelector(selector, node[key], recursive); value != nil {
				nodeList = append(nodeList, value...)
			}
		}
	}
	return nodeList
}

// applyNameSelector returns a value from the given current node.
// If the current node is a map, it returns the value associated with the name.
// Otherwise, it returns nil.
func applyNameSelector(selector *ir.NameSelector, node any, recursive bool) NodeList {
	var nodeList NodeList
	if node, ok := node.(map[string]any); ok {
		// Applying the name-selector to an object node selects a member value whose name equals the member name `M` or
		// selects nothing if there is no such member value.
		if value, ok := node[selector.Name]; ok {
			nodeList = append(nodeList, value)
		}
		if recursive {
			var keys []string
			for key := range node {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				if value := applyNameSelector(selector, node[key], recursive); value != nil {
					nodeList = append(nodeList, value...)
				}
			}
		}
	}
	if node, ok := node.([]any); ok && recursive {
		for _, value := range node {
			if value := applyNameSelector(selector, value, recursive); value != nil {
				nodeList = append(nodeList, value...)
			}
		}
	}
	return nodeList
}

// applySliceSelector returns a list of nodes from the given current node.
// If the current node is a list, it returns a slice of elements.
// Otherwise, it returns nil.
func applySliceSelector(selector *ir.SliceSelector, node any, recursive bool) NodeList {
	var nodeList NodeList
	if node, ok := node.([]any); ok {
		if selector.Step == 0 {
			// When step is 0, no elements are selected.
			return nil
		}

		if 0 < selector.Step {
			// When step is negative, elements are selected in reverse order. Thus, for example, 5:1:-2 selects elements
			// with indices 5 and 3 (in that order), and ::-1 selects all the elements of an array in reverse order.
			if selector.End == -1 {
				selector.End = len(node)
			}
			for i := max(selector.Start, 0); i < min(selector.End, len(node)); i += selector.Step {
				nodeList = append(nodeList, node[i])
			}
		} else {
			// The array slice expression start:end:step selects elements at indices starting at start, incrementing by
			// step, and ending with end (which is itself excluded).
			if selector.Start == 0 && selector.End == -1 {
				selector.Start = len(node) - 1
			}
			for i := min(selector.Start, len(node)); max(selector.End, -1) < i; i += selector.Step {
				nodeList = append(nodeList, node[i])
			}
		}

		if recursive {
			for _, value := range node {
				nodeList = append(
					nodeList,
					applySliceSelector(
						selector,
						value,
						recursive,
					)...,
				)
			}
		}
	}
	if node, ok := node.(map[string]any); ok && recursive {
		var keys []string
		for key := range node {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			if value := applySliceSelector(selector, node[key], recursive); value != nil {
				nodeList = append(nodeList, value...)
			}
		}
	}
	return nodeList
}

// applyWildcardSelector returns a list of nodes from the given current node.
// If the current node is a map, it returns a list of values sorted by keys.
// If the current node is a list, it returns the list itself.
// Otherwise, it returns nil.
func applyWildcardSelector(node any, recursive bool) NodeList {
	var nodeList NodeList
	switch node := node.(type) {
	case map[string]any:
		var keys []string
		for key := range node {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := node[key]
			nodeList = append(nodeList, value)

			if recursive {
				nodeList = append(
					nodeList,
					applyWildcardSelector(
						value,
						recursive,
					)...,
				)
			}
		}
		return nodeList
	case []any:
		nodeList = append(nodeList, node...)
		if recursive {
			for _, value := range node {
				nodeList = append(
					nodeList,
					applyWildcardSelector(
						value,
						recursive,
					)...,
				)
			}
		}
	}
	return nodeList
}

// applySelector returns a list of nodes from the given current node.
// A selector produces a node list consisting of zero or more children of the input value.
func (ctx *context) applySelector(selector ir.Selector, node any, recursive bool) NodeList {
	switch selector := selector.(type) {
	case *ir.NameSelector:
		if v := applyNameSelector(selector, node, recursive); v != nil {
			return v
		}
		return nil
	case *ir.WildcardSelector:
		return applyWildcardSelector(node, recursive)
	case *ir.SliceSelector:
		return applySliceSelector(selector, node, recursive)
	case *ir.IndexSelector:
		if v := applyIndexSelector(selector, node, recursive); v != nil {
			return v
		}
		return nil
	case *ir.FilterSelector:
		return ctx.applyFilterSelector(selector, node, recursive)
	default:
		panic(fmt.Sprintf("unsupported selector type: %T", selector))
	}
}
