package lib

type Contexter interface {
	FlagNames() []string

	NumFlags() int
}

// Provider is a global variable representing the CLI's global context. It should
// be set (read: overridden) by each cloud provider
var CloudProvider Provider

// Run executes all the methods of a Provider for each command
func Run(context Contexter, commander Commander) {

	if CloudProvider == nil {
		panic("You must set the Cloud variable")
	}

	globalOptions := CloudProvider.NewGlobalOptionser(context)
	err := globalOptions.ParseGlobalOptions()
	if err != nil {
		CloudProvider.ErrExit1(err)
	}

	authenticater := CloudProvider.NewAuthenticater(globalOptions, commander.ServiceClientType())
	serviceClient, err := authenticater.Authenticate()
	if err != nil {
		CloudProvider.ErrExit1(err)
	}

	commander.SetServiceClient(serviceClient)

	defer func() {
		if authFromCacher, ok := authenticater.(AuthFromCacher); ok {
			authFromCacher.StoreCredentials()
		}
	}()

	err = commander.HandleFlags()
	if err != nil {
		CloudProvider.ErrExit1(err)
	}

	inChannel := CloudProvider.InputChannel()

	go CloudProvider.FillInputChannel(commander, inChannel)

	outChannel := CloudProvider.ResultsChannel()

	waiter, ok := commander.(Waiter)
	switch ok && waiter.ShouldWait() {
	case true:
		go waiter.ExecuteAndWait(inChannel, outChannel)
	case false:
		go commander.Execute(inChannel, outChannel)
	}

	outputter := CloudProvider.NewResultOutputter(globalOptions, commander)

	for result := range outChannel {
		err = outputter.OutputResult(result)
		if err != nil {
			CloudProvider.ErrExit1(err)
		}
	}
}
