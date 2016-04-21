package main

import (
	"fmt"
	"io"
	"reflect"

	"github.com/gophercloud/cli/lib"
)

type Result struct {
	// err will store any error encountered while processing the command.
	err   error
	value interface{}
	rt    ResultTyper
}

func (r Result) GetValue() interface{} {
	return r.value
}

func (r Result) SetValue(v interface{}) {
	r.value = v
}

func (r Result) Types() []lib.ResultTyper {
	return []lib.ResultTyper{
		//reflect.TypeOf(map[string]interface{}{}),
		//reflect.TypeOf([]map[string]interface{}{}),
		//reflect.TypeOf(io.Reader),
		//reflect.TypeOf(""),
		MapStringInterface,
	}
}

func (r *Result) SetError(err error) {
	r.err = err
}

func (r Result) GetError() error {
	return r.err
}

func (r *Result) HandleEmpty() error {
	for _, t := range r.Types() {
		if reflect.TypeOf(r.value).AssignableTo(reflect.TypeOf(t)) {
			r.SetValue(t.HandleEmpty())
		}
	}
}

type MapStringInterface map[string]interface{}

func (rt MapStringInterface) HandleEmpty() (interface{}, error) {
	return fmt.Sprintf("No result found.\n"), nil
}

func (rt MapStringInterface) Print() {

}

type SliceOfMapStringInterface []map[string]interface{}

func (rt SliceOfMapStringInterface) HandleEmpty() (interface{}, error) {
	return fmt.Sprintf("No results found.\n"), nil
}

func (rt SliceOfMapStringInterface) Print() {

}

type String string

func (rt String) HandleEmpty() (interface{}, error) {
	return fmt.Sprintf("No result found.\n"), nil
}

func (rt String) Print() {

}

type IOReader io.Reader

func (rt IOReader) HandleEmpty() (interface{}, error) {
	return fmt.Sprintf("No result found.\n"), nil
}

func (rt IOReader) Print() {

}
