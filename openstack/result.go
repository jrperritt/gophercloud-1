package openstack

import "github.com/gophercloud/cli/lib"

type DebugMsg string

// Result satisfies the Resulter interface
type Result struct {
	value interface{}
	rt    lib.ResultTyper
}

func NewResult() *Result {
	return new(Result)
}

/*
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
*/
