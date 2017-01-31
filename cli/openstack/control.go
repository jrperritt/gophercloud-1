package openstack

import (
	"log"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"gopkg.in/urfave/cli.v1"
)

type globalctx struct {
	*cli.Context
	ExecuteResults, ResultsRunCommand      chan (interface{})
	Logger                                 *logger
	DoneChan, ProgressDoneChan, UpdateChan chan (interface{})
	wgExecute, wgProgress                  *sync.WaitGroup
}

// gctx represents the global context
var gctx *globalctx

// Action is the common function all commands run
func Action(ctx *cli.Context, cmd interfaces.Commander) error {
	cmd.SetContext(ctx)

	gctx = &globalctx{
		ExecuteResults:    make(chan interface{}),
		ResultsRunCommand: make(chan interface{}),
		DoneChan:          make(chan interface{}),
		ProgressDoneChan:  make(chan interface{}),
		UpdateChan:        make(chan interface{}),
		wgExecute:         new(sync.WaitGroup),
		wgProgress:        new(sync.WaitGroup),
	}

	gopts, err := globalopts(ctx)
	if err != nil {
		return ErrExit1{err}
	}

	l := new(logger)
	l.Logger = log.New(ctx.App.Writer, "", log.LstdFlags)
	l.debug = gopts.debug
	gctx.Logger = l

	ao := &authopts{
		cmd:     cmd,
		region:  gopts.region,
		gao:     gopts.authOptions,
		nocache: gopts.noCache,
		urltype: gopts.urlType,
	}
	sc, err := auth(ao)
	if err != nil {
		return ErrExit1{err}
	}

	if !gopts.noCache {
		defer func() {
			cachecreds(ao, sc)
		}()
	}

	cmd.SetServiceClient(sc)

	gctx.Logger.Debugln("Running HandleInterfaceFlags...")
	err = cmd.HandleInterfaceFlags()
	if err != nil {
		return ErrExit1{err}
	}

	gctx.Logger.Debugln("Running HandleFlags...")
	err = cmd.HandleFlags()
	if err != nil {
		return ErrExit1{err}
	}

	gctx.Logger.Debugln("Running RunCommand...")
	go runcmd(cmd)

	gctx.Logger.Debugln("Running OutputResults...")
	err = outres(cmd)
	if err != nil {
		return ErrExit1{err}
	}

	return nil
}
