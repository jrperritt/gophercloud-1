package lib

// GlobalOptionser is an interface that types implement when they contain
// global options
type GlobalOptionser interface {
	Sources() []string
	ParseGlobalOptions() error
	//GlobalOptions() []GlobalOptioner
	Defaults() []GlobalOptioner
	MethodsMap() map[string]func() error
	Validate() error
	Set() error
}

// GlobalOptioner is an interface that a global option implements
type GlobalOptioner interface {
	Name() string
	Value() interface{}
	From() string
}
