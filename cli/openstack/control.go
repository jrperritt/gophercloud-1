package openstack

import (
	"sync"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"gopkg.in/urfave/cli.v1"
)

type globalContext struct {
	CommandContext                         *cli.Context
	ServiceClient                          *gophercloud.ServiceClient
	GlobalOptions                          *GlobalOptions
	ExecuteResults, ResultsRunCommand      chan (interface{})
	Command                                interfaces.Commander
	Logger                                 *logger
	DoneChan, ProgressDoneChan, UpdateChan chan (interface{})
	doneChan                               chan (bool)
	wgExecute, wgProgress                  *sync.WaitGroup
}

// GC represents the global context
var GC *globalContext

// Action is the common method all commands run
func Action(ctx *cli.Context, commander interfaces.Commander) error {
	GC = &globalContext{
		ExecuteResults:    make(chan interface{}),
		ResultsRunCommand: make(chan interface{}),
		wgExecute:         new(sync.WaitGroup),
		wgProgress:        new(sync.WaitGroup),
		Command:           commander,
		CommandContext:    ctx,
		DoneChan:          make(chan interface{}),
		ProgressDoneChan:  make(chan interface{}),
		UpdateChan:        make(chan interface{}),
	}

	err := SetGlobalOptions()
	if err != nil {
		return ErrExit1{err}
	}

	GC.Logger = GC.GlobalOptions.logger

	err = Authenticate()
	if err != nil {
		return ErrExit1{err}
	}

	if !GC.GlobalOptions.noCache {
		defer func() {
			StoreCredentials()
		}()
	}

	GC.Command.SetServiceClient(GC.ServiceClient)
	GC.Command.SetContext(GC.CommandContext)

	GC.GlobalOptions.logger.Debugln("Running HandleInterfaceFlags...")
	err = GC.Command.HandleInterfaceFlags()
	if err != nil {
		return ErrExit1{err}
	}

	GC.GlobalOptions.logger.Debugln("Running HandleFlags...")
	err = GC.Command.HandleFlags()
	if err != nil {
		return ErrExit1{err}
	}

	GC.GlobalOptions.logger.Debugln("Running RunCommand...")
	go RunCommand()

	GC.GlobalOptions.logger.Debugln("Running OutputResults...")
	err = OutputResults()
	if err != nil {
		return ErrExit1{err}
	}

	return nil
}
