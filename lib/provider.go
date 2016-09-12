package lib

// Provider should be implemented by one object per cloud provider
type Provider interface {
	NewGlobalOptionser(Contexter) GlobalOptionser

	NewAuthenticater(GlobalOptionser, string) Authenticater

	InputChannel() chan interface{}

	FillInputChannel(Commander, chan interface{})

	ResultsChannel() chan interface{}

	NewResultOutputter(GlobalOptionser, Commander) Outputter

	ErrExit1(error)
}
