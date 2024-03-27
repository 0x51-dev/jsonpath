package jsonpath

import (
	"reflect"
	"testing"
)

func TestGetTags(t *testing.T) {
	var obj struct {
		Tagged int `jsonpath:"$.tagged"`
	}

	v := reflect.TypeOf(obj)
	tag, err := getTags(v.Field(0))
	if err != nil {
		t.Fatal(err)
	}
	if tag.path.Query() != "$['tagged']" {
		t.Fatalf("unexpected query: %s", tag.path.Query())
	}
}
