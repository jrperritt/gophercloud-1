package traits

import "gopkg.in/urfave/cli.v1"

type Waitable struct {
	wait      bool
	donechout chan<- interface{}
}

func (c *Waitable) WaitFor(raw interface{}) {
	//openstack.GC.DoneChan <- raw
	c.donechout <- raw
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
