package lib

type Resourcer interface {
	GetStdInParams() interface{}
	GetResult() Resulter
}
