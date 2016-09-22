package commands

import (
	"github.com/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

type WaitCommand struct {
	UpdateChan, CompleteChan chan (interface{})
	ErrorChan                chan (error)
	OutChan                  chan (interface{})
	Wait                     bool
}

func (c *WaitCommand) WaitFor(raw interface{}) {
	openstack.GC.DoneChan <- raw
}

func (c *WaitCommand) WaitFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "wait",
			Usage: "[optional] If provided, wait to return until the operation is complete.",
		},
	}
}

func (c *WaitCommand) ShouldWait() bool {
	return c.Wait
}
