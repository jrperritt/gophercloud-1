package traits

import (
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"

	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	quiet      bool
	donechin   chan interface{}
	donechout  chan interface{}
	updatechin chan interface{}
	listench   chan interfaces.ProgressStatuser
	*ProgressInfo
}

func (c *Progressable) BarID(raw interface{}) string {
	return raw.(string)
}

func (c *Progressable) InitProgress(donechout chan interface{}) {
	c.donechout = donechout
	c.donechin = make(chan interface{})
	c.updatechin = make(chan interface{})
	c.listench = make(chan interfaces.ProgressStatuser)
	go c.Listen(c.listench)
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

func (c *Progressable) ProgListenCh() chan interfaces.ProgressStatuser {
	return c.listench
}

type PercentageProgressable struct {
	Progressable
}

func (c *PercentageProgressable) InitProgress(donech chan interface{}) {
	c.ProgressInfo = NewProgressInfo(BarPercentage)
	c.Progressable.InitProgress(donech)
}

func (c *PercentageProgressable) ShowBar(id string) {
	s := new(ProgressStatusStart)
	s.Name = id
	c.ProgListenCh() <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(ProgressStatusComplete)
			s.Name = id
			c.ProgListenCh() <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(ProgressStatusUpdate)
			s.Name = id
			s.Increment = int(r.(float64))
			c.ProgListenCh() <- s
		}
	}
}

type ByteSizesMap struct {
	*sync.Mutex
	m map[string]int
}

func (b *ByteSizesMap) Set(id string, v int) {
	b.Lock()
	b.m[id] = v
	b.Unlock()
}

func (b *ByteSizesMap) Get(id string) (v int) {
	b.Lock()
	v = b.m[id]
	b.Unlock()
	return
}

type BytesProgressable struct {
	Progressable
	Sizes *ByteSizesMap
}

func (c *BytesProgressable) InitByteSizesMap() {
	c.Sizes = new(ByteSizesMap)
	c.Sizes.m = make(map[string]int)
	c.Sizes.Mutex = new(sync.Mutex)
}

func (c *BytesProgressable) InitProgress(donech chan interface{}) {
	c.ProgressInfo = NewProgressInfo(BarBytes)
	c.Progressable.InitProgress(donech)
}

func (c *BytesProgressable) ShowBar(id string) {
	s := new(ProgressStatusStart)
	s.Name = id
	s.TotalSize = c.Sizes.Get(id)
	c.ProgListenCh() <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(ProgressStatusComplete)
			s.Name = id
			c.ProgListenCh() <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(ProgressStatusUpdate)
			s.Name = id
			s.Increment = r.(int)
			c.ProgListenCh() <- s
		}
	}
}

type TextProgressable struct {
	Progressable
	//RunningMsg, DoneMsg string
}

func (c *TextProgressable) InitProgress(donech chan interface{}) {
	c.ProgressInfo = NewProgressInfo(BarText)
	c.Progressable.InitProgress(donech)
}

func (c *TextProgressable) ShowBar(id string) {
	s := new(ProgressStatusStart)
	s.Name = id
	c.ProgListenCh() <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(ProgressStatusComplete)
			s.Name = id
			c.ProgListenCh() <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(ProgressStatusUpdate)
			s.Name = id
			s.Msg = r.(string)
			c.ProgListenCh() <- s
		}
	}
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
