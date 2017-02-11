package traits

import (
	"fmt"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gosuri/uiprogress"
)

type BarType uint8

const (
	BarPercentage BarType = iota
	BarBytes
	BarText
)

type ProgressBar struct {
	*uiprogress.Bar
	id string
}

func (b *ProgressBar) ID() string {
	return b.id
}

func (b *ProgressBar) Complete(_ interfaces.ProgressCompleteStatuser) {
	b.Set(b.Total)
}

func (b *ProgressBar) Error(s interfaces.ProgressErrorStatuser) string {
	return fmt.Sprintf("[ERROR: %s] %s", s.Err(), b.id)
}

func (b *ProgressBar) Reset() {
	b.Set(0)
}

type ProgressBarPercentage struct {
	*ProgressBar
}

func (p *ProgressBarPercentage) BarSize() int {
	return 100
}

func (p *ProgressBarPercentage) Start(s interfaces.ProgressStartStatuser) interfaces.ProgressBarrer {
	statusBar := uiprogress.NewBar(s.BarSize()).AppendCompleted().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%s  %s", s.BarID(), b.TimeElapsedString())
	})
	p.ProgressBar = new(ProgressBar)
	p.ProgressBar.Bar = statusBar
	return p
}

func (b *ProgressBarPercentage) Update(s interfaces.ProgressUpdateStatuser) {
	incr := s.(*ProgressStatusUpdate).Increment
	b.Incr()
	b.Set(incr - 1)
}

type ProgressBarBytes struct {
	*ProgressBar
	size int
}

func (p *ProgressBarBytes) BarSize() int {
	return p.size
}

func (b *ProgressBarBytes) Start(s interfaces.ProgressStartStatuser) interfaces.ProgressBarrer {
	statusBar := uiprogress.NewBar(s.BarSize()).AppendCompleted().PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("%s  %s", s.BarID(), b.TimeElapsedString())
	})
	b.ProgressBar = new(ProgressBar)
	b.ProgressBar.Bar = statusBar
	return b
}

func (b *ProgressBarBytes) Update(s interfaces.ProgressUpdateStatuser) {
	incr := s.(*ProgressStatusUpdate).Increment
	b.Set(b.Current() + incr)
	//lib.Log.Debugf("b: %+v\n", b)
}

type ProgressBarText struct {
	*ProgressBar
	donemsg, runmsg string
}

func (p *ProgressBarText) BarSize() int {
	return 2
}

func (p *ProgressBarText) DoneMsg() string {
	return p.donemsg
}

func (p *ProgressBarText) RunMsg() string {
	return p.runmsg
}

func (p *ProgressBarText) Start(s interfaces.ProgressStartStatuser) interfaces.ProgressBarrer {
	statusBar := uiprogress.NewBar(2).PrependElapsed().AppendFunc(func(b *uiprogress.Bar) string {
		var msg string
		switch b.Current() {
		case b.Total:
			msg = p.DoneMsg()
		default:
			msg = p.RunMsg()
		}
		return fmt.Sprintf("%s\t%s", s.BarID(), msg)
	})
	p.ProgressBar = new(ProgressBar)
	p.ProgressBar.Bar = statusBar
	p.setBarToText()
	return p
}

func (b *ProgressBarText) Update(_ interfaces.ProgressUpdateStatuser) {
	b.Incr()
	b.Set(0)
}

func (b *ProgressBarText) setBarToText() {
	b.LeftEnd = ' '
	b.RightEnd = ' '
	b.Head = ' '
	b.Fill = ' '
	b.Empty = ' '
}

type ProgressStatsBar struct {
	*ProgressBarText
	totals struct {
		*sync.RWMutex
		active    int
		completed int
		errored   int
	}
}

func NewProgressStatsBar() *ProgressStatsBar {
	c := new(ProgressStatsBar)
	c.totals.RWMutex = new(sync.RWMutex)
	c.ProgressBarText = new(ProgressBarText)
	c.ProgressBarText.ProgressBar = new(ProgressBar)
	c.ProgressBarText.Bar = uiprogress.NewBar(2).PrependFunc(func(b *uiprogress.Bar) string {
		c.totals.Lock()
		defer c.totals.Unlock()
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", c.totals.active, c.totals.completed, c.totals.errored)
	}).PrependElapsed()

	//c.ProgressBarText.ProgressBar.Bar = p.Bars[0]
	c.id = "summary"
	c.setBarToText()
	return c
}
