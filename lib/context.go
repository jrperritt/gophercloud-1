package lib

type Contexter interface {
	FlagNames() []string
	NumFlags() int
}

// Provider should be implemented by one object per cloud provider
type CloudProvider interface {
	// the name of the cli tool
	Name() string

	NewGlobalOptionser(Contexter) GlobalOptionser

	NewAuthenticater(GlobalOptionser, string) Authenticater

	InputChannel() chan interface{}

	FillInputChannel(Commander, chan interface{})

	ResultsChannel() chan interface{}

	NewResultOutputter(GlobalOptionser, Commander) Outputter

	ErrExit1(error)
}

// Provider is a global variable representing the CLI's global context. It should
// be set (read: overridden) by each cloud provider
var Provider CloudProvider

// Run executes all the methods of a Provider for each command
func Run(context Contexter, commander Commander) {

	if Provider == nil {
		panic("You must set the Cloud variable")
	}

	globalOptions := Provider.NewGlobalOptionser(context)
	err := globalOptions.ParseGlobalOptions()
	if err != nil {
		Provider.ErrExit1(err)
	}

	authenticater := Provider.NewAuthenticater(globalOptions, commander.ServiceClientType())
	serviceClient, err := authenticater.Authenticate()
	if err != nil {
		Provider.ErrExit1(err)
	}

	commander.SetServiceClient(serviceClient)

	defer func() {
		if authFromCacher, ok := authenticater.(AuthFromCacher); ok {
			authFromCacher.StoreCredentials()
		}
	}()

	err = commander.HandleFlags()
	if err != nil {
		Provider.ErrExit1(err)
	}

	inChannel := Provider.InputChannel()

	go Provider.FillInputChannel(commander, inChannel)

	outChannel := Provider.ResultsChannel()

	for item := range inChannel {
		go commander.Execute(item, outChannel)
	}

	outputter := Provider.NewResultOutputter(globalOptions, commander)

	for result := range outChannel {
		err = outputter.OutputResult(result)
		if err != nil {
			Provider.ErrExit1(err)
		}
	}
}
