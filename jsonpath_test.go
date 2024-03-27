package jsonpath_test

import (
	"github.com/0x51-dev/jsonpath"
	"reflect"
	"testing"
)

// https://www.rfc-editor.org/rfc/rfc9535.html#name-array-slice-selector
func TestPath_Apply_arraySliceSelector(t *testing.T) {
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
	}.Run(t, []any{"a", "b", "c", "d", "e", "f", "g"})
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-filter-selector
func TestPath_Apply_filterSelector(t *testing.T) {
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
			result:  []any{[]any{1, 2}, []any{1, 2}},
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
	}.Run(t, map[string]any{
		"a": []any{3, 5, 1, 2, 4, 6,
			map[string]any{"b": "j"},
			map[string]any{"b": "k"},
			map[string]any{"b": map[string]any{}},
			map[string]any{"b": "kilo"},
		},
		"o": map[string]any{"p": 1, "q": 2, "r": 3, "s": 5, "t": map[string]any{"u": 6}},
		"e": "f",
	})
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-index-selector
func TestPath_Apply_indexSelector(t *testing.T) {
	testCases{
		{
			comment: "Element of array",
			query:   "$[1]",
			result:  "b",
		},
		{
			comment: "Element of array, from the end",
			query:   "$[-2]",
			result:  "a",
		},
	}.Run(t, []any{"a", "b"})
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-name-selector
func TestPath_Apply_nameSelector(t *testing.T) {
	testCases{
		{
			comment: "Named value in a nested object",
			query:   "$.o['j j']",
			result:  map[string]any{"k.k": 3},
		},
		{
			comment: "Nesting further down",
			query:   "$.o['j j']['k.k']",
			result:  3,
		},
		{
			comment: "Different delimiter in the query, unchanged Normalized Path",
			query:   "$.o[\"j j\"][\"k.k\"]",
			result:  3,
		},
		{
			comment: "Unusual member names",
			query:   "$[\"'\"][\"@\"]",
			result:  2,
		},
	}.Run(t, map[string]any{
		"o": map[string]any{"j j": map[string]any{"k.k": 3}},
		"'": map[string]any{"@": 2},
	})
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-wildcard-selector
func TestPath_Apply_wildcardSelector(t *testing.T) {
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
			result:  []any{[]any{1, 2}, []any{1, 2}},
		},
		{
			comment: "Array members",
			query:   "$.a[*]",
			result:  []any{5, 3},
		},
	}.Run(t, map[string]any{
		"o": map[string]any{"j": 1, "k": 2},
		"a": []any{5, 3},
	})
}

type testCase struct {
	comment string
	query   string
	result  any
}

func (c testCase) Run(t *testing.T, v any) {
	t.Run(c.comment, func(t *testing.T) {
		q, err := jsonpath.New(c.query)
		if err != nil {
			t.Fatal(err)
		}
		r, err := q.Apply(v)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(r, c.result) {
			t.Errorf("expected %v, got %v", c.result, r)
		}
	})
}

type testCases []testCase

func (c testCases) Run(t *testing.T, v any) {
	for _, test := range c {
		test.Run(t, v)
	}
}
