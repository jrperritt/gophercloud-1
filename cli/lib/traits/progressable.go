package traits

import (
	"sync"

	"github.com/gophercloud/gophercloud/cli/openstack"
	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	quiet               bool
	donechin, donechout chan interface{}
	updatechin          chan interface{}
	wg                  *sync.WaitGroup
	*openstack.ProgressInfo
}

func (c *Progressable) BarID(raw interface{}) string {
	return raw.(string)
}

func (c *Progressable) InitProgress(donechout chan interface{}) {
	c.updatechin = make(chan interface{})
	c.donechin = make(chan interface{})
	c.donechout = donechout
	go c.Listen()
	c.ProgressInfo.Start()
}

func (c *Progressable) ProgDoneChIn() chan interface{} {
	return c.donechin
}

func (c *Progressable) ProgDoneChOut() chan interface{} {
	return c.donechout
}

func (c *Progressable) ProgUpdateChIn() chan interface{} {
	return c.updatechin
}

func (c *Progressable) WG() *sync.WaitGroup {
	return c.wg
}

func (c *Progressable) ProgressFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] If provided, only final results are printed.",
		},
	}
}

func (c *Progressable) SetProgress(b bool) {
	c.quiet = b
}

func (c *Progressable) ShouldProgress() bool {
	return !c.quiet
}

type PercentageProgressable struct {
	Progressable
}

func (c *PercentageProgressable) InitProgress(donech chan interface{}) {
	c.ProgressInfo = openstack.NewProgressInfo(0)
	c.Progressable.InitProgress(donech)
}

func (c *PercentageProgressable) ShowBar(id string) {
	s := new(openstack.ProgressStatusStart)
	s.Name = id
	c.StartChan <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(openstack.ProgressStatusComplete)
			s.Name = id
			c.ProgressInfo.CompleteChan <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(openstack.ProgressStatusUpdate)
			s.Name = id
			s.Increment = int(r.(float64))
			c.ProgressInfo.UpdateChan <- s
		}
	}
}

type BytesProgressable struct {
	Progressable
}

func (c *BytesProgressable) InitProgress(donech chan interface{}) {
	c.ProgressInfo = openstack.NewProgressInfo(1)
	c.Progressable.InitProgress(donech)
}

func (c *BytesProgressable) ShowBar(id string) {
	s := new(openstack.ProgressStatusStart)
	s.Name = id
	c.StartChan <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(openstack.ProgressStatusComplete)
			s.Name = id
			c.ProgressInfo.CompleteChan <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(openstack.ProgressStatusUpdate)
			s.Name = id
			s.Increment = int(r.(float64))
			c.ProgressInfo.UpdateChan <- s
		}
	}
}

type TextProgressable struct {
	Progressable
	//RunningMsg, DoneMsg string
}

func (c *TextProgressable) InitProgress(donech chan interface{}) {
	c.ProgressInfo = openstack.NewProgressInfo(2)
	c.Progressable.InitProgress(donech)
}

func (c *TextProgressable) ShowBar(id string) {
	s := new(openstack.ProgressStatusStart)
	s.Name = id
	c.StartChan <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(openstack.ProgressStatusComplete)
			s.Name = id
			c.ProgressInfo.CompleteChan <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(openstack.ProgressStatusUpdate)
			s.Name = id
			s.Msg = r.(string)
			c.ProgressInfo.UpdateChan <- s
		}
	}
}
