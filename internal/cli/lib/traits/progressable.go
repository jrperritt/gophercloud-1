package traits

import (
	"fmt"
	"math"
	"sync"

	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/vbauerster/mpb"

	"gopkg.in/urfave/cli.v1"
)

type Progressable struct {
	quiet   bool
	startch chan interfaces.ProgressItemer
	stats   ProgressStatsBar
	info    *mpb.Progress
	*sync.Mutex
}

func (p *Progressable) ProgressFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "[optional] If provided, only final results are printed.",
		},
	}
}

func (p *Progressable) SetProgress(b bool) {
	p.quiet = b
}

func (p *Progressable) ShouldProgress() bool {
	return !p.quiet
}

func (p *Progressable) InitProgress() {
	p.startch = make(chan interfaces.ProgressItemer)
	p.info = mpb.New(nil)
	p.Mutex = new(sync.Mutex)
}

func (p *Progressable) AddSummaryBar() {
	p.stats.Bar = p.info.AddBar(2).PrependFunc(func(b *mpb.Statistics) string {
		p.Lock()
		defer p.Unlock()
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", p.stats.totals.active, p.stats.totals.completed, p.stats.totals.errored)
	})
	p.stats.id = "summary"
	p.stats.setBarToText()
}

func (p *Progressable) ProgStartCh() chan interfaces.ProgressItemer {
	return p.startch
}

func (p *Progressable) StartBar() {
	p.Lock()
	p.stats.totals.active++
	p.Unlock()
}

func (p *Progressable) CompleteBar() {
	p.Lock()
	p.stats.totals.active--
	p.stats.totals.completed++
	p.Unlock()
}

func (p *Progressable) ErrorBar() {
	p.Lock()
	p.stats.totals.active--
	p.stats.totals.errored++
	p.Unlock()
}

type BytesProgressable struct {
	Progressable
}

func (p *BytesProgressable) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	b := new(ProgressBarBytes)
	if p.ShouldProgress() {
		b.Bar = p.info.AddBar(pi.Size()).PrependElapsed(6).PrependFunc(func(_ *mpb.Statistics) string {
			return fmt.Sprintf("\t")
		}).PrependETA(4).AppendPercentage().AppendFunc(func(s *mpb.Statistics) string {
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
