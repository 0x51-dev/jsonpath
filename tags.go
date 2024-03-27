package jsonpath

import (
	"reflect"
)

type Tag struct {
	path *Path
}

func getTags(field reflect.StructField) (*Tag, error) {
	tag := field.Tag.Get("jsonpath")
	if tag == "" {
		return nil, nil
	}
	q, err := New(tag)
	if err != nil {
		return nil, err
	}
	return &Tag{
		path: q,
	}, nil
}
