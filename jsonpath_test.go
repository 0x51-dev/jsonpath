package jsonpath_test

import (
	"encoding/json"
	"fmt"
	"github.com/0x51-dev/jsonpath"
	"reflect"
	"testing"
)

func ExamplePath_Apply() {
	example := map[string]any{
		"store": map[string]any{
			"book": []any{
				map[string]any{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]any{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
				map[string]any{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
				map[string]any{
					"category": "fiction",
					"author":   "J. R. R. Tolkien",
					"title":    "The Lord of the Rings",
					"isbn":     "0-395-19395-8",
					"price":    22.99,
				},
			},
			"bicycle": map[string]any{
				"color": "red",
				"price": 399,
			},
		},
	}

	// This is just to make the response more readable.
	marshal := func(v []any) string {
		var w any = v
		if len(v) == 1 {
			w = v[0]
		}
		raw, _ := json.Marshal(w)
		return string(raw)
	}

	q, _ := jsonpath.New("$.store.book[*].author")
	fmt.Println("The authors of all books in the store:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..author")
	fmt.Println("All authors:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$.store")
	fmt.Println("All things in the store, which are some books and a red bicycle:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$.store..price")
	fmt.Println("Prices of everything in the store:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[2]")
	fmt.Println("The third book:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[2].author")
	fmt.Println("The third book's author:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[2].publisher")
	fmt.Println("The third book's publisher:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[-1]")
	fmt.Println("The last book on order:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[0,1]")
	fmt.Println("The first two books:", marshal(q.Apply(example)))
	q, _ = jsonpath.New("$..book[:2]")
	fmt.Println("The first two books:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[?@.isbn]")
	fmt.Println("All books with an ISBN number:", marshal(q.Apply(example)))

	q, _ = jsonpath.New("$..book[?@.price<10]")
	fmt.Println("All books cheaper than 10:", marshal(q.Apply(example)))

	// Output:
	// The authors of all books in the store: ["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]
	// All authors: ["Nigel Rees","Evelyn Waugh","Herman Melville","J. R. R. Tolkien"]
	// All things in the store, which are some books and a red bicycle: {"bicycle":{"color":"red","price":399},"book":[{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]}
	// Prices of everything in the store: [399,8.95,12.99,8.99,22.99]
	// The third book: {"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}
	// The third book's author: "Herman Melville"
	// The third book's publisher: null
	// The last book on order: {"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}
	// The first two books: [{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"}]
	// The first two books: [{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Evelyn Waugh","category":"fiction","price":12.99,"title":"Sword of Honour"}]
	// All books with an ISBN number: [{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"},{"author":"J. R. R. Tolkien","category":"fiction","isbn":"0-395-19395-8","price":22.99,"title":"The Lord of the Rings"}]
	// All books cheaper than 10: [{"author":"Nigel Rees","category":"reference","price":8.95,"title":"Sayings of the Century"},{"author":"Herman Melville","category":"fiction","isbn":"0-553-21311-3","price":8.99,"title":"Moby Dick"}]
}

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
