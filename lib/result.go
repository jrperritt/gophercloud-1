package lib

type Resulter interface {
	GetValue() interface{}
	SetValue(interface{})
	GetError() error
	SetError(error)
	Types() []ResultTyper
	HandleEmpty() error
	Print()
}

type ResultTyper interface {
	HandleEmpty() (interface{}, error)
	Print()
}
