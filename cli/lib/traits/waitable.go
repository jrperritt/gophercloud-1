package traits

import (
	"github.com/gophercloud/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

type Waitable struct {
	Wait bool
}

func (c *Waitable) WaitFor(raw interface{}) {
	openstack.GC.DoneChan <- raw
}

func (c *Waitable) WaitFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "wait",
			Usage: "[optional] If provided, wait to return until the operation is complete.",
		},
	}
}

func (c *Waitable) ShouldWait() bool {
	return c.Wait
}
