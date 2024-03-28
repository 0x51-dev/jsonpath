package jsonpath_test

import (
	"github.com/0x51-dev/jsonpath"
	"reflect"
	"testing"
)

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-8
func TestPath_Apply_childSegment(t *testing.T) {
	example := []any{"a", "b", "c", "d", "e", "f", "g"}
	testCases{
		{
			comment: "Indices",
			query:   "$[0, 3]",
			result:  []any{"a", "d"},
		},
		{
			comment: "Slice and index",
			query:   "$[0:2, 5]",
			result:  []any{"a", "b", "f"},
		},
		{
			comment: "Duplicate entries",
			query:   "$[0, 0]",
			result:  []any{"a", "a"},
		},
	}.Run(t, example)
}

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-9
func TestPath_Apply_descendantSegment(t *testing.T) {
	example := map[string]any{
		"o": map[string]any{"j": 1, "k": 2},
		"a": []any{5, 3, []any{map[string]any{"j": 4}, map[string]any{"k": 6}}},
	}
	testCases{
		{
			comment: "Object values",
			query:   "$..j",
			result:  []any{4, 1},
		},
		{
			comment: "Array values",
			query:   "$..[0]",
			result:  []any{5, map[string]any{"j": 4}},
		},
		{
			comment: "All values",
			query:   "$..[*]",
			result: []any{
				[]any{5, 3, []any{map[string]any{"j": 4}, map[string]any{"k": 6}}},
				5,
				3,
				[]any{map[string]any{"j": 4}, map[string]any{"k": 6}},
				map[string]any{"j": 4},
				map[string]any{"k": 6},
				4,
				6,
				map[string]any{"j": 1, "k": 2},
				1,
				2,
			},
		},
		{
			comment: "Input value is visited",
			query:   "$..o",
			result: []any{
				map[string]any{"j": 1, "k": 2},
			},
		},
		{
			comment: "Non-deterministic ordering",
			query:   "$.o..[*, *]",
			result: []any{
				1, 2, 1, 2,
			},
		},
		{
			comment: "Multiple segments",
			query:   "$.a..[0, 1]",
			result: []any{
				5,
				map[string]any{"j": 4},
				3,
				map[string]any{"k": 6},
			},
		},
	}.Run(t, example)
}

type testCase struct {
	comment string
	query   string
	result  jsonpath.NodeList
}

func (c testCase) Run(t *testing.T, v any) {
	t.Run(c.comment, func(t *testing.T) {
		q, err := jsonpath.New(c.query)
		if err != nil {
			t.Fatal(err)
		}
		r := q.Apply(v)
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
