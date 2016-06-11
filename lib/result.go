package lib

type Resulter interface {
	GetValue() interface{}
	SetValue(interface{})
	GetError() error
	SetError(error)
	GetTypes() []ResultTyper
	SetType()
	GetType() ResultTyper
	GetEmptyValue() interface{}
	Print()
}

type ResultTyper interface {
	GetEmptyValue() interface{}
	Print()
}
