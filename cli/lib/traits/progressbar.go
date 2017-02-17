package traits

import (
	"fmt"
	"sync"
	"time"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/vbauerster/mpb"
)

type BarType uint8

const (
	BarPercentage BarType = iota
	BarBytes
	BarText
)

type ProgressBar struct {
	*mpb.Bar
	id string
}

func (b *ProgressBar) ID() string {
	return b.id
}

func (b *ProgressBar) Complete(_ interfaces.ProgressCompleteStatuser) {
	b.Completed()
	for b.InProgress() {
		time.Sleep(100 * time.Millisecond)
	}
	/*
		for b.Bar.GetStatistics().Current != b.Bar.GetStatistics().Total {
			time.Sleep(100 * time.Millisecond)
		}
	*/
}

func (b *ProgressBar) Error(s interfaces.ProgressErrorStatuser) string {
	return fmt.Sprintf("[ERROR: %s] %s", s.Err(), b.id)
}

type ProgressBarPercentage struct {
	ProgressBar
}

func (p *ProgressBarPercentage) Start(s interfaces.ProgressStartStatuser) interfaces.ProgressBarrer {
	return p
}

func (b *ProgressBarPercentage) Update(s interfaces.ProgressUpdateStatuser) {
	b.Incr(s.Change().(int))
}

type ProgressBarBytes struct {
	ProgressBar
	size int64
}

func (b *ProgressBarBytes) Start(s interfaces.ProgressStartStatuser) interfaces.ProgressBarrer {
	return b
}

func (b *ProgressBarBytes) Update(s interfaces.ProgressUpdateStatuser) {
	var incr int
	switch inc := s.Change().(type) {
	case int:
		incr = inc
	case float64:
		incr = int(inc)
	}
	b.Incr(incr)
}

type ProgressBarText struct {
	ProgressBar
	donemsg, runmsg string
}

func (p *ProgressBarText) DoneMsg() string {
	return p.donemsg
}

func (p *ProgressBarText) RunMsg() string {
	return p.runmsg
}

func (p *ProgressBarText) Start(s interfaces.ProgressStartStatuser) interfaces.ProgressBarrer {
	return p
}

func (b *ProgressBarText) Update(_ interfaces.ProgressUpdateStatuser) {
}

func (b *ProgressBarText) setBarToText() {
	b.SetLeftEnd(' ')
	b.SetRightEnd(' ')
	b.SetTip(' ')
	b.SetFill(' ')
	b.SetEmpty(' ')
}

type ProgressStatsBar struct {
	ProgressBarText
	totals struct {
		*sync.RWMutex
		active    int
		completed int
		errored   int
	}
}
