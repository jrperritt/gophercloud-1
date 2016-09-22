package commands

import (
	"github.com/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

type ProgressCommand struct {
	WaitCommand
	Quiet bool
	*openstack.ProgressInfo
	OutChan chan (interface{})
}

func (c *ProgressCommand) BarID(raw interface{}) string {
	return raw.(string)
}

func (c *ProgressCommand) InitProgress() {
	go c.Listen()
	c.ProgressInfo.Start()
}

func (c *ProgressCommand) ProgressFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] If provided, only final results are printed.",
		},
	}
}

func (c *ProgressCommand) ShouldProgress() bool {
	if c.Quiet {
		return false
	}
	c.Wait = false
	return true
}

type TextProgressCommand struct {
	ProgressCommand
	//RunningMsg, DoneMsg string
}

func (c *TextProgressCommand) ShowBar(id string) {
	s := new(openstack.ProgressStatusStart)
	s.Name = id
	c.StartChan <- s

	for {
		select {
		case r := <-openstack.GC.DoneChan:
			s := new(openstack.ProgressStatusComplete)
			s.Name = id
			c.ProgressInfo.CompleteChan <- s
			openstack.GC.ProgressDoneChan <- r
			return
		case r := <-openstack.GC.UpdateChan:
			s := new(openstack.ProgressStatusUpdate)
			s.Name = id
			s.Msg = r.(string)
			c.ProgressInfo.UpdateChan <- s
		}
	}
}

type PercentageProgressCommand struct {
	ProgressCommand
}

type BytesProgressCommand struct {
	ProgressCommand
}
