package lib

import (
	"fmt"

	"github.com/codegangsta/cli"
)

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

// Context is a global variable representing the CLI's global context. It should
// be set (read: overridden) by each cloud provider
var Context Provider

// Run executes all the methods of a Provider for each command
func Run(cliContext *cli.Context, commander Commander) {

	debugChannel := make(chan string)
	commander.SetDebugChannel(debugChannel)
	go func() {
		for msg := range debugChannel {
			fmt.Printf("[DEBUGCHANNEL] %s\n", msg)
		}
	}()

	if Context == nil {
		panic("You must set the Cloud variable")
	}

	globalOptions := Context.NewGlobalOptionser(cliContext)
	err := globalOptions.ParseGlobalOptions()
	if err != nil {
		Context.ErrExit1(err)
	}

	authenticater := Context.NewAuthenticater(globalOptions, commander.ServiceClientType())
	serviceClient, err := authenticater.Authenticate()
	if err != nil {
		Context.ErrExit1(err)
	}

	commander.SetServiceClient(serviceClient)

	err = commander.HandleFlags()
	if err != nil {
		Context.ErrExit1(err)
	}

	resultsChannel := Context.ResultsChannel()

	fmt.Println("running command...")
	err = commander.RunCommand(resultsChannel)
	if err != nil {
		Context.ErrExit1(err)
	}

	fmt.Println("creating result outputter...")
	outputter := Context.NewResultOutputter(globalOptions)

	fmt.Println("fetching results...")
	for result := range resultsChannel {
		fmt.Println("outputting a result: ", result)
		err = outputter.OutputResult(result)
		if err != nil {
			Context.ErrExit1(err)
		}
	}
}
