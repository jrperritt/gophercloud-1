package openstack

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/gophercloud"

	"gopkg.in/urfave/cli.v1"
)

type globalContext struct {
	CommandContext                         *cli.Context
	ServiceClient                          *gophercloud.ServiceClient
	GlobalOptions                          *GlobalOptions
	ExecuteResults, ResultsRunCommand      chan (interface{})
	Command                                Commander
	Logger                                 *logrus.Logger
	DoneChan, ProgressDoneChan, UpdateChan chan (interface{})
	doneChan                               chan (bool)
	wgExecute, wgProgress                  *sync.WaitGroup
}

var GC *globalContext

func Action(ctx *cli.Context, commander Commander) error {
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

	GC.GlobalOptions.logger.Debug("Running HandleFlags...")
	err = GC.Command.HandleFlags()
	if err != nil {
		return ErrExit1{err}
	}

	GC.GlobalOptions.logger.Debug("Running RunCommand...")
	go RunCommand()

	GC.GlobalOptions.logger.Debug("Running OutputResults...")
	err = OutputResults()
	if err != nil {
		return ErrExit1{err}
	}

	return nil
}
