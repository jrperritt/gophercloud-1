package traits

import (
	"sync"

	"gopkg.in/urfave/cli.v1"
)

type Waitable struct {
	wait   bool
	donech chan interface{}
	wg     *sync.WaitGroup
}

func (c *Waitable) WaitDoneCh() chan interface{} {
	return c.donech
}

func (c *Waitable) WG() *sync.WaitGroup {
	return c.wg
}

func (c *Waitable) WaitFor(raw interface{}) {
	c.donech <- raw
}

func (c *Waitable) WaitFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "wait",
			Usage: "[optional] If provided, wait to return until the operation is complete.",
		},
	}
}

func (c *Waitable) SetWait(b bool) {
	c.wait = b
}

func (c *Waitable) ShouldWait() bool {
	return c.wait
}
