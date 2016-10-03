package commands

import (
	"github.com/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	Quiet bool
	*openstack.ProgressInfo
}

func (c *Progressable) BarID(raw interface{}) string {
	return raw.(string)
}

func (c *Progressable) InitProgress() {
	go c.Listen()
	c.ProgressInfo.Start()
}

func (c *Progressable) ProgressFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] If provided, only final results are printed.",
		},
	}
}

func (c *Progressable) ShouldProgress() bool {
	return !c.Quiet
}

type TextProgressable struct {
	Progressable
	//RunningMsg, DoneMsg string
}

func (c *TextProgressable) ShowBar(id string) {
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

type PercentageProgressable struct {
	Progressable
}

type BytesProgressable struct {
	Progressable
}
