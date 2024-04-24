package jsonpath_test

import (
	"github.com/0x51-dev/jsonpath"
	"testing"
)

// https://www.rfc-editor.org/rfc/rfc9535.html#name-examples-7
func TestFunctions(t *testing.T) {
	for _, test := range []struct {
		query     string
		wellTyped bool
	}{
		{
			query:     "$[?length(@) < 3]",
			wellTyped: true,
		},
		{
			query: "$[?length(@.*) < 3]",
		},
		{
			query:     "$[?count(@.*) == 1]",
			wellTyped: true,
		},
		{
			query: "$[?count(1) == 1]",
		},
		{
			query:     "$[?match(@.timezone, 'Europe/.*')]",
			wellTyped: true,
		},
		{
			query: "$[?match(@.timezone, 'Europe/.*') == true]",
		},
		{
			query:     "$[?value(@..color) == \"red\"]",
			wellTyped: true,
		},
		{
			query: "$[?value(@..color)]",
		},
	} {
		t.Run(test.query, func(t *testing.T) {
			if _, err := jsonpath.New(test.query); (err == nil) != test.wellTyped {
				t.Fatalf("%q, %t", err, test.wellTyped)
			}
		})
	}
}
