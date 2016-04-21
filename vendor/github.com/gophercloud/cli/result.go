package main

import (
	"io"
	"reflect"
)

type Result struct{}

func (r Result) Types() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(map[string]interface{}{}),
		reflect.TypeOf([]map[string]interface{}{}),
		reflect.TypeOf(io.Reader),
		reflect.TypeOf(""),
	}
}
