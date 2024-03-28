package cmp_test

import (
	"fmt"
	"github.com/0x51-dev/jsonpath"
	"github.com/0x51-dev/jsonpath/cmp"
	"strings"
	"testing"
)

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-6
func TestComparisons(t *testing.T) {
	var source = map[string]any{
		"obj": map[string]any{"x": "y"},
		"arr": []any{2, 3},
	}
	for _, test := range []struct {
		valueA, valueB any
		pathA, pathB   string
		op             string
		result         bool
	}{
		{
			pathA:  "$.absent1",
			pathB:  "$.absent2",
			op:     "==",
			result: true,
		},
		{
			pathA:  "$.absent1",
			pathB:  "$.absent2",
			op:     "<=",
			result: true,
		},
		{
			pathA:  "$.absent1",
			valueB: "g",
			op:     "==",
		},
		{
			pathA: "$.absent1",
			pathB: "$.absent2",
			op:    "!=",
		},
		{
			pathA:  "$.absent1",
			valueB: "g",
			op:     "!=",
			result: true,
		},
		{
			valueA: 1,
			valueB: 2,
			op:     "<=",
			result: true,
		},
		{
			valueA: 1,
			valueB: 2,
			op:     ">",
		},
		{
			valueA: 13,
			valueB: "13",
			op:     "==",
		},
		{
			valueA: "a",
			valueB: "b",
			op:     "<=",
			result: true,
		},
		{
			valueA: "a",
			valueB: "b",
			op:     ">",
		},
		{
			pathA: "$.obj",
			pathB: "$.arr",
			op:    "==",
		},
		{
			pathA:  "$.obj",
			pathB:  "$.arr",
			op:     "!=",
			result: true,
		},
		{
			pathA:  "$.obj",
			pathB:  "$.obj",
			op:     "==",
			result: true,
		},
		{
			pathA: "$.obj",
			pathB: "$.obj",
			op:    "!=",
		},
		{
			pathA:  "$.arr",
			pathB:  "$.arr",
			op:     "==",
			result: true,
		},
		{
			pathA: "$.arr",
			pathB: "$.arr",
			op:    "!=",
		},
		{
			pathA:  "$.obj",
			valueB: 17,
			op:     "==",
		},
		{
			pathA:  "$.obj",
			valueB: 17,
			op:     "!=",
			result: true,
		},
		{
			pathA: "$.obj",
			pathB: "$.arr",
			op:    "<=",
		},
		{
			pathA: "$.obj",
			pathB: "$.arr",
			op:    "<",
		},
		{
			pathA:  "$.obj",
			pathB:  "$.obj",
			op:     "<=",
			result: true,
		},
		{
			pathA:  "$.arr",
			pathB:  "$.arr",
			op:     "<=",
			result: true,
		},
		{
			valueA: 1,
			pathB:  "$.arr",
			op:     "<=",
		},
		{
			valueA: 1,
			pathB:  "$.arr",
			op:     ">=",
		},
		{
			valueA: 1,
			pathB:  "$.arr",
			op:     ">",
		},
		{
			valueA: 1,
			pathB:  "$.arr",
			op:     "<",
		},
		{
			valueA: true,
			valueB: true,
			op:     "<=",
			result: true,
		},
		{
			valueA: true,
			valueB: true,
			op:     ">",
		},
	} {
		var name strings.Builder
		if test.pathA != "" {
			name.WriteString(test.pathA)
			q, err := jsonpath.New(test.pathA)
			if err != nil {
				t.Fatalf("error parsing %q: %v", test.pathA, err)
			}
			valueA := []any(q.Apply(source))
			if err != nil {
				t.Fatalf("error applying %q: %v", test.pathA, err)
			}
			test.valueA = valueA
		} else {
			name.WriteString(fmt.Sprintf("%v", test.valueA))
		}
		name.WriteString(fmt.Sprintf(" %v ", test.op))
		if test.pathB != "" {
			name.WriteString(test.pathB)
			q, err := jsonpath.New(test.pathB)
			if err != nil {
				t.Fatalf("error parsing %q: %v", test.pathB, err)
			}
			valueB := []any(q.Apply(source))
			if err != nil {
				t.Fatalf("error applying %q: %v", test.pathB, err)
			}
			test.valueB = valueB
		} else {
			name.WriteString(fmt.Sprintf("%v", test.valueB))
		}
		t.Run(name.String(), func(t *testing.T) {
			if err := cmp.Compare(test.valueA, test.valueB, test.op); (err == nil) != test.result {
				t.Errorf("compare(%v, %v, %q) = %v; want %v", test.valueA, test.valueB, test.op, err, test.result)
			}
		})
	}
}
