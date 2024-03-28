package jsonpath_test

import "testing"

// https://www.rfc-editor.org/rfc/rfc9535.html#name-filter-selector
func TestPath_Apply_filterSelector(t *testing.T) {
	example := map[string]any{
		"a": []any{3, 5, 1, 2, 4, 6,
			map[string]any{"b": "j"},
			map[string]any{"b": "k"},
			map[string]any{"b": map[string]any{}},
			map[string]any{"b": "kilo"},
		},
		"o": map[string]any{"p": 1, "q": 2, "r": 3, "s": 5, "t": map[string]any{"u": 6}},
		"e": "f",
	}
	testCases{
		{
			comment: "Member value comparison",
			query:   "$.a[?@.b == 'kilo']",
			result:  []any{map[string]any{"b": "kilo"}},
		},
		{
			comment: "Equivalent query with enclosing parentheses",
			query:   "$.a[?(@.b == 'kilo')]",
			result:  []any{map[string]any{"b": "kilo"}},
		},
		{
			comment: "Array value comparison",
			query:   "$.a[?@>3.5]",
			result:  []any{5, 4, 6},
		},
		{
			comment: "Array value existence",
			query:   "$.a[?@.b]",
			result: []any{
				map[string]any{"b": "j"},
				map[string]any{"b": "k"},
				map[string]any{"b": map[string]any{}},
				map[string]any{"b": "kilo"},
			},
		},
		{
			comment: "Existence of non-singular queries",
			query:   "$[?@.*]",
			result: []any{
				[]any{3, 5, 1, 2, 4, 6,
					map[string]any{"b": "j"},
					map[string]any{"b": "k"},
					map[string]any{"b": map[string]any{}},
					map[string]any{"b": "kilo"},
				},
				map[string]any{
					"p": 1, "q": 2, "r": 3, "s": 5,
					"t": map[string]any{"u": 6},
				},
			},
		},
		{
			comment: "Nested filters",
			query:   "$[?@[?@.b]]",
			result: []any{
				[]any{3, 5, 1, 2, 4, 6,
					map[string]any{"b": "j"},
					map[string]any{"b": "k"},
					map[string]any{"b": map[string]any{}},
					map[string]any{"b": "kilo"},
				},
			},
		},
		{
			comment: "Non-deterministic ordering",
			query:   "$.o[?@<3, ?@<3]",
			result:  []any{1, 2, 1, 2},
		},
		{
			comment: "Array value logical OR",
			query:   `$.a[?@<2 || @.b == "k"]`,
			result:  []any{1, map[string]any{"b": "k"}},
		},
		{
			comment: "Array value regular expression match",
			query:   `$.a[?match(@.b, "[jk]")]`,
			result: []any{
				map[string]any{"b": "j"},
				map[string]any{"b": "k"},
			},
		},
		{
			comment: "Array value regular expression search",
			query:   `$.a[?search(@.b, "[jk]")]`,
			result: []any{
				map[string]any{"b": "j"},
				map[string]any{"b": "k"},
				map[string]any{"b": "kilo"},
			},
		},
		{
			comment: "Object value logical AND",
			query:   "$.o[?@>1 && @<4]",
			result:  []any{2, 3},
		},
		{
			comment: "Object value logical OR",
			query:   "$.o[?@.u || @.x]",
			result:  []any{map[string]any{"u": 6}},
		},
		{
			comment: "Comparison of queries with no values",
			query:   "$.a[?@.b == $.x]",
			result:  []any{3, 5, 1, 2, 4, 6},
		},
		{
			comment: "Comparisons of primitive and of structured values",
			query:   "$.a[?@ == @]",
			result: []any{3, 5, 1, 2, 4, 6,
				map[string]any{"b": "j"},
				map[string]any{"b": "k"},
				map[string]any{"b": map[string]any{}},
				map[string]any{"b": "kilo"},
			},
		},
	}.Run(t, example)
}
