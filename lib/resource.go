package lib

type Resourcer interface {
	//Result() Resulter
	//SetResult(Resulter)
	StdinValue() interface{}
	SetStdinValue(interface{})
	StdinField() string
	SetStdinField(string)
}
