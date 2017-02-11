package traits

import (
	"sync"
	"time"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gosuri/uiprogress"

	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	quiet      bool
	donechin   chan interface{}
	donechout  chan interface{}
	updatechin chan interface{}
	listench   chan interfaces.ProgressStatuser
	bt         BarType
	stats      *ProgressStatsBar
	info       *uiprogress.Progress
	barsByName map[string]interfaces.ProgressBarrer
	namesByBar map[interfaces.ProgressBarrer]string
}

func (c *Progressable) BarID(raw interface{}) string {
	return raw.(string)
}

func NewProgress(donechout chan interface{}, bt BarType) *Progressable {
	c := new(Progressable)
	c.donechout = donechout
	c.donechin = make(chan interface{})
	c.updatechin = make(chan interface{})
	c.listench = make(chan interfaces.ProgressStatuser)

	c.info = uiprogress.New()
	c.info.RefreshInterval = time.Second * 1

	c.bt = bt

	c.stats = NewProgressStatsBar()
	c.info.AddBar(2)
	c.info.Bars[0] = c.stats.Bar

	c.barsByName = make(map[string]interfaces.ProgressBarrer, 0)
	c.namesByBar = make(map[interfaces.ProgressBarrer]string, 0)

	c.info.Start()
	return c
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

func (p *Progressable) Update() {
	p.info.Bars[0].Incr()
	p.info.Bars[0].Set(0)
}

func (p *Progressable) StartBar(s interfaces.ProgressStartStatuser) {
	if statusBarInfo, ok := p.barsByName[s.BarID()]; ok {
		p.stats.totals.Lock()
		p.stats.totals.active++
		p.stats.totals.errored--
		p.stats.totals.Unlock()
		statusBarInfo.Reset()
		return
	}

	var bar interfaces.ProgressBarrer
	switch p.bt {
	case BarPercentage:
		//bar = new(ProgressBarPercentage).Start(s)
	case BarBytes:
		bb := new(ProgressBarBytes)
		uib := p.info.AddBar(s.BarSize()).PrependElapsed().AppendCompleted().AppendFunc(func(b *uiprogress.Bar) string {
			return s.BarID()
		})
		bb.ProgressBar = new(ProgressBar)
		bb.ProgressBar.Bar = uib
		bar = bb
	case BarText:
	}

	p.barsByName[s.BarID()] = bar
	p.namesByBar[bar] = s.BarID()
	p.stats.totals.Lock()
	p.stats.totals.active++
	p.stats.totals.Unlock()
}

func (p *Progressable) UpdateBar(s interfaces.ProgressUpdateStatuser) {
	if statusBarInfo := p.barsByName[s.BarID()]; statusBarInfo != nil {
		statusBarInfo.Update(s)
	}
}

func (p *Progressable) CompleteBar(s interfaces.ProgressCompleteStatuser) {
	if statusBarInfo := p.barsByName[s.BarID()]; statusBarInfo != nil {
		statusBarInfo.Complete(s)
		p.stats.totals.Lock()
		p.stats.totals.active--
		p.stats.totals.completed++
		p.stats.totals.Unlock()
	}
}

func (p *Progressable) ErrorBar(s interfaces.ProgressErrorStatuser) {
	if statusBarInfo := p.barsByName[s.BarID()]; statusBarInfo != nil {
		p.namesByBar[statusBarInfo] = statusBarInfo.Error(s)
		p.stats.totals.Lock()
		p.stats.totals.active--
		p.stats.totals.errored++
		p.stats.totals.Unlock()
	}
}

type PercentageProgressable struct {
	Progressable
}

func (c *PercentageProgressable) InitProgress(donechout chan interface{}) {
	c.Progressable = *NewProgress(donechout, BarPercentage)
}

func (c *PercentageProgressable) ShowBar(id string) {
	s := new(ProgressStatusStart)
	c.ProgListenCh() <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(ProgressStatusComplete)
			s.id = id
			c.ProgListenCh() <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(ProgressStatusUpdate)
			s.id = id
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

func (c *BytesProgressable) InitProgress(donechout chan interface{}) {
	c.Progressable = *NewProgress(donechout, BarBytes)
	c.Sizes = new(ByteSizesMap)
	c.Sizes.m = make(map[string]int)
	c.Sizes.Mutex = new(sync.Mutex)
	go c.Listen(c.ProgListenCh())
}

func (c *BytesProgressable) Listen(statusch chan interfaces.ProgressStatuser) {
	listen(c, statusch)
}

func (c *BytesProgressable) ShowBar(id string) {
	s := new(ProgressStatusStart)
	s.id = id
	s.size = c.Sizes.Get(id)
	c.ProgListenCh() <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(ProgressStatusComplete)
			s.id = id
			c.ProgListenCh() <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(ProgressStatusUpdate)
			s.id = id
			s.Increment = r.(int)
			c.ProgListenCh() <- s
		}
	}
}

type TextProgressable struct {
	Progressable
	RunningMsg, DoneMsg string
}

func (c *TextProgressable) InitProgress(donechout chan interface{}) {
	c.Progressable = *NewProgress(donechout, BarText)
}

func (c *TextProgressable) ShowBar(id string) {
	s := new(ProgressStatusStart)
	s.id = id
	c.ProgListenCh() <- s

	for {
		select {
		case r := <-c.ProgDoneChIn():
			s := new(ProgressStatusComplete)
			s.id = id
			c.ProgListenCh() <- s
			c.ProgDoneChOut() <- r
			return
		case r := <-c.ProgUpdateChIn():
			s := new(ProgressStatusUpdate)
			s.id = id
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

func listen(p interfaces.Progresser, statusch chan interfaces.ProgressStatuser) {
	defer close(statusch)
	for status := range statusch {
		switch s := status.(type) {
		case *ProgressStatusStart:
			p.StartBar(s)
		case *ProgressStatusUpdate:
			p.UpdateBar(s)
		case *ProgressStatusComplete:
			p.CompleteBar(s)
		case *ProgressStatusError:
			p.ErrorBar(s)
		}
		p.Update()
	}
}
