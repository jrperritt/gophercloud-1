package lib

import "github.com/codegangsta/cli"

// Provider should be implemented by one object per cloud provider
type Provider interface {
	// the name of the cli tool
	Name() string

	NewGlobalOptionser(*cli.Context) GlobalOptionser

	NewAuthenticater(GlobalOptionser, string) Authenticater

	ResultsChannel() chan Resulter

	NewResultOutputter(GlobalOptionser) Outputter

	ErrExit1(error)
}

// Cloud is a global variable representing the CLI's global context. It should
// be set (read: overridden) by each cloud provider
var Cloud Provider

// Run executes all the methods of a Provider for each command
func Run(cliContext *cli.Context, commander Commander) {

	if Cloud == nil {
		panic("You must set the Cloud variable")
	}

	globalOptions := Cloud.NewGlobalOptionser(cliContext)
	err := globalOptions.ParseGlobalOptions()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	authenticater := Cloud.NewAuthenticater(globalOptions, commander.ServiceClientType())
	serviceClient, err := authenticater.Authenticate()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	commander.SetServiceClient(serviceClient)

	err = commander.HandleFlags()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	resultsChannel := Cloud.ResultsChannel()

	err = commander.RunCommand(resultsChannel)
	if err != nil {
		Cloud.ErrExit1(err)
	}

	outputter := Cloud.NewResultOutputter(globalOptions)

	for result := range resultsChannel {
		err = outputter.OutputResult(result)
		if err != nil {
			Cloud.ErrExit1(err)
		}
	}
}
