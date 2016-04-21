package lib

type GlobalOptionser interface {
	InitGlobalOptions() error
	Sources() []string
	GetGlobalOptions() []GlobalOptioner
	Defaults() []GlobalOptioner
	MethodsMap() map[string]func() error
	Validate() error
	Set() error
}

type GlobalOptioner interface {
}
