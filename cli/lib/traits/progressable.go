package traits

import (
	"fmt"
	"math"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/vbauerster/mpb"

	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	quiet   bool
	startch chan interfaces.ProgressItemer
	stats   ProgressStatsBar
	info    *mpb.Progress
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

func (c *Progressable) InitProgress() {
	c.startch = make(chan interfaces.ProgressItemer)
	c.info = mpb.New(nil)
}

func (c *Progressable) AddSummaryBar() {
	//c.stats = new(ProgressStatsBar)
	c.stats.totals.RWMutex = new(sync.RWMutex)
	c.stats.Bar = c.info.AddBar(2).PrependFunc(func(b *mpb.Statistics) string {
		c.stats.totals.Lock()
		defer c.stats.totals.Unlock()
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", c.stats.totals.active, c.stats.totals.completed, c.stats.totals.errored)
	})
	c.stats.id = "summary"
	c.stats.setBarToText()
}

func (c *Progressable) ProgStartCh() chan interfaces.ProgressItemer {
	return c.startch
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
}

func (p *BytesProgressable) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	b := new(ProgressBarBytes)
	if p.ShouldProgress() {
		b.Bar = p.info.AddBar(pi.Size()).PrependElapsed(6).PrependETA(6).AppendPercentage().AppendFunc(func(s *mpb.Statistics) string {
			return pi.ID()
		})
	}
	return b
}

type PercentageProgressable struct {
	Progressable
}

func (p *PercentageProgressable) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	b := new(ProgressBarBytes)
	if p.ShouldProgress() {
		b.Bar = p.info.AddBar(pi.Size()).PrependElapsed(6).AppendPercentage().AppendFunc(func(s *mpb.Statistics) string {
			return pi.ID()
		})
	}
	return b
}

type TextProgressable struct {
	BytesProgressable
	donesmg, runmsg string
}

func (p *TextProgressable) RunningMsg() string {
	return p.runmsg
}

func (p *TextProgressable) DoneMsg() string {
	return p.donesmg
}

func (p *TextProgressable) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	b := new(ProgressBarText)
	if p.ShouldProgress() {
		b.Bar = p.info.AddBar(pi.Size()).PrependElapsed(6).AppendPercentage().AppendFunc(func(s *mpb.Statistics) string {
			return pi.ID()
		}).AppendFunc(func(s *mpb.Statistics) string {
			if s.Current == s.Total {
				return p.DoneMsg()
			}
			return p.RunningMsg()
		})
	}
	return b
}

type BytesStreamProgressable struct {
	TextProgressable
}

func (p *BytesStreamProgressable) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	b := new(ProgressBarBytes)
	if p.ShouldProgress() {
		var size int64
		if pi.Size() != 0 {
			size = pi.Size()
		} else {
			size = math.MaxInt64
		}
		b.Bar = p.info.AddBar(size).SetWidth(50).PrependElapsed(0).AppendFunc(func(s *mpb.Statistics) string {
			return fmt.Sprintf("%s\t%s", pi.ID(), mpb.Format(s.Current).To(mpb.UnitBytes))
		})
	}
	return b
}
