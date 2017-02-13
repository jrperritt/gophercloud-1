package traits

import (
	"fmt"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/vbauerster/mpb"

	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	quiet      bool
	updatechin chan interface{}
	startch    chan interfaces.ProgressItemer
	bt         BarType
	stats      *ProgressStatsBar
	info       *mpb.Progress
}

func NewProgressable(bt BarType) *Progressable {
	c := new(Progressable)
	c.updatechin = make(chan interface{})
	c.startch = make(chan interfaces.ProgressItemer)
	c.info = mpb.New(nil)
	c.bt = bt
	c.stats = new(ProgressStatsBar)
	c.stats.totals.RWMutex = new(sync.RWMutex)
	return c
}

func (c *Progressable) AddSummaryBar() {
	c.stats.ProgressBarText = new(ProgressBarText)
	c.stats.ProgressBarText.ProgressBar = new(ProgressBar)
	c.stats.ProgressBarText.Bar = c.info.AddBar(2).PrependFunc(func(b *mpb.Statistics) string {
		c.stats.totals.Lock()
		defer c.stats.totals.Unlock()
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", c.stats.totals.active, c.stats.totals.completed, c.stats.totals.errored)
	}).PrependElapsed(2)
	c.stats.id = "summary"
	c.stats.setBarToText()
}

func (c *Progressable) ProgStartCh() chan interfaces.ProgressItemer {
	return c.startch
}

func (c *Progressable) ProgUpdateChIn() chan interface{} {
	return c.updatechin
}

func (p *Progressable) StartBar() {
	p.stats.totals.Lock()
	p.stats.totals.active++
	p.stats.totals.Unlock()
}

func (p *Progressable) CompleteBar() {
	p.stats.totals.Lock()
	p.stats.totals.active--
	p.stats.totals.completed++
	p.stats.totals.Unlock()
}

func (p *Progressable) ErrorBar() {
	p.stats.totals.Lock()
	p.stats.totals.active--
	p.stats.totals.errored++
	p.stats.totals.Unlock()
}

type BytesProgressable struct {
	Progressable
	updatechin chan interface{}
}

func (c *BytesProgressable) ProgUpdateCh() chan interface{} {
	return c.updatechin
}

func (p *BytesProgressable) InitProgress() {
	//p.Progressable = *NewProgressable(BarBytes)
	p.updatechin = make(chan interface{})
	p.startch = make(chan interfaces.ProgressItemer)
	p.info = mpb.New(nil)
	p.bt = BarBytes
	p.stats = new(ProgressStatsBar)
	p.stats.totals.RWMutex = new(sync.RWMutex)
}

func (p *BytesProgressable) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	b := new(ProgressBarBytes)
	b.ProgressBar = new(ProgressBar)
	if p.ShouldProgress() {
		b.ProgressBar.Bar = p.info.AddBar(pi.Size()).PrependElapsed(2).AppendPercentage().AppendFunc(func(s *mpb.Statistics) string {
			return pi.ID()
		})
	}
	return b
}

type PercentageProgressable struct {
	BytesProgressable
}

type TextProgressable struct {
	BytesProgressable
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
