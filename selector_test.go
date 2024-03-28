package jsonpath_test

import "testing"

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-5
func TestPath_Apply_arraySliceSelector(t *testing.T) {
	example := []any{"a", "b", "c", "d", "e", "f", "g"}
	testCases{
		{
			comment: "Slice with default step",
			query:   "$[1:3]",
			result:  []any{"b", "c"},
		},
		{
			comment: "Slice with no end index",
			query:   "$[5:]",
			result:  []any{"f", "g"},
		},
		{
			comment: "Slice with step 2",
			query:   "$[1:5:2]",
			result:  []any{"b", "d"},
		},
		{
			comment: "Slice with negative step",
			query:   "$[5:1:-2]",
			result:  []any{"f", "d"},
		},
		{
			comment: "Slice in reverse order",
			query:   "$[::-1]",
			result:  []any{"g", "f", "e", "d", "c", "b", "a"},
		},
	}.Run(t, example)
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-4
func TestPath_Apply_indexSelector(t *testing.T) {
	example := []any{"a", "b"}
	testCases{
		{
			comment: "Element of array",
			query:   "$[1]",
			result:  []any{"b"},
		},
		{
			comment: "Element of array, from the end",
			query:   "$[-2]",
			result:  []any{"a"},
		},
	}.Run(t, example)
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-2
func TestPath_Apply_nameSelector(t *testing.T) {
	example := map[string]any{
		"o": map[string]any{"j j": map[string]any{"k.k": 3}},
		"'": map[string]any{"@": 2},
	}
	testCases{
		{
			comment: "Named value in a nested object",
			query:   "$.o['j j']",
			result:  []any{map[string]any{"k.k": 3}},
		},
		{
			comment: "Nesting further down",
			query:   "$.o['j j']['k.k']",
			result:  []any{3},
		},
		{
			comment: "Different delimiter in the query, unchanged Normalized Path",
			query:   "$.o[\"j j\"][\"k.k\"]",
			result:  []any{3},
		},
		{
			comment: "Unusual member names",
			query:   "$[\"'\"][\"@\"]",
			result:  []any{2},
		},
	}.Run(t, example)
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-3
func TestPath_Apply_wildcardSelector(t *testing.T) {
	example := map[string]any{
		"o": map[string]any{"j": 1, "k": 2},
		"a": []any{5, 3},
	}
	testCases{
		{
			comment: "Object values",
			query:   "$[*]",
			result:  []any{[]any{5, 3}, map[string]any{"j": 1, "k": 2}},
		},
		{
			comment: "Object values",
			query:   "$.o[*]",
			result:  []any{1, 2},
		},
		{
			comment: "Non-deterministic ordering",
			query:   "$.o[*, *]",
			result:  []any{1, 2, 1, 2},
		},
		{
			comment: "Array members",
			query:   "$.a[*]",
			result:  []any{5, 3},
		},
	}.Run(t, example)
}
