package openstack

import (
	"fmt"
	"io"
	"reflect"

	"github.com/gophercloud/cli/lib"
)

// Result satisfies the Resulter interface
type Result struct {
	// err will store any error encountered while processing the command.
	err   error
	value interface{}
	rt    lib.ResultTyper
}

// GetValue satisfies the Resulter.GetValue method
func (r *Result) GetValue() interface{} {
	return r.value
}

// SetValue satisfies the Resulter.SetValue method
func (r *Result) SetValue(v interface{}) {
	r.value = v
}

// GetTypes satisfies the Resulter.GetTypes method
func (r *Result) GetTypes() []lib.ResultTyper {
	return []lib.ResultTyper{
		MapStringInterface{},
		SliceOfMapStringInterface{},
		IOReader{},
		String(""),
	}
}

// SetType satisfies the Resulter.SetType method
func (r *Result) SetType() {
	if err := r.GetError(); err != nil {
		r.SetValue(err)
	}

	for _, t := range r.GetTypes() {
		if reflect.TypeOf(r.GetValue()).AssignableTo(reflect.TypeOf(t)) {
			r.rt = t
		}
	}
}

// GetType satisfies the Resulter.GetType method
func (r *Result) GetType() lib.ResultTyper {
	return r.rt
}

// SetError satisfies the Resulter.SetError method
func (r *Result) SetError(err error) {
	r.err = err
}

// GetError satisfies the Resulter.GetError method
func (r *Result) GetError() error {
	return r.err
}

// GetEmptyValue satisfies the Resulter.GetEmptyValue method
func (r *Result) GetEmptyValue() interface{} {
	if r.GetType() == nil {
		r.SetType()
	}
	return r.GetType().GetEmptyValue()
}

// Print satisfies the Resulter.Print method
func (r *Result) Print() {

	r.GetType().Print()

}

// MapStringInterface satisfies the ResultTyper interface
type MapStringInterface map[string]interface{}

// GetEmptyValue satisfies the ResultTyper.GetEmptyValue method
func (rt MapStringInterface) GetEmptyValue() interface{} {
	return fmt.Sprintf("No result found.\n")
}

func (rt MapStringInterface) Print() {

}

// SliceOfMapStringInterface satisfies the ResultTyper interface
type SliceOfMapStringInterface []map[string]interface{}

// GetEmptyValue satisfies the ResultTyper.GetEmptyValue method
func (rt SliceOfMapStringInterface) GetEmptyValue() interface{} {
	return fmt.Sprintf("No results found.\n")
}

func (rt SliceOfMapStringInterface) Print() {

}

// String satisfies the ResultTyper interface
type String string

// GetEmptyValue satisfies the ResultTyper.GetEmptyValue method
func (rt String) GetEmptyValue() interface{} {
	return fmt.Sprintf("No result found.\n")
}

func (rt String) Print() {

}

// IOReader satisfies the ResultTyper interface
type IOReader struct {
	io.Reader
}

// GetEmptyValue satisfies the ResultTyper.GetEmptyValue method
func (rt IOReader) GetEmptyValue() interface{} {
	return fmt.Sprintf("No result found.\n")
}

func (rt IOReader) Print() {

}
