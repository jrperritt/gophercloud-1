package openstack

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
	"gopkg.in/urfave/cli.v1"
)

type BarType uint8

const (
	BarPecentage BarType = iota
	BarBytes
	BarText
)

type Progresser interface {
	Waiter
	InitProgress()
	BarID(item interface{}) string
	ShowBar(id string)
	ShouldProgress() bool
	ProgressFlags() []cli.Flag
}

type ProgressStatus struct {
	Name string
}

type ProgressStatusStart struct {
	ProgressStatus
	TotalSize int
}

type ProgressStatusError struct {
	ProgressStatus
	Err error
}

type ProgressStatusUpdate struct {
	ProgressStatus
	Increment int
	Msg       string
}

type ProgressStatusComplete struct {
	ProgressStatus
	Result interface{}
}

type ProgressItem struct {
	Content ProgressBarrer
}

type ProgressInfo struct {
	Totals struct {
		*sync.RWMutex
		Active    int
		Completed int
		Errored   int
	}
	*uiprogress.Progress
	StartChan           chan *ProgressStatusStart
	UpdateChan          chan *ProgressStatusUpdate
	ErrorChan           chan *ProgressStatusError
	CompleteChan        chan *ProgressStatusComplete
	SummaryBar          *ProgressBarText
	BarType             BarType
	BarsByName          map[string]*ProgressItem
	NamesByBar          map[*ProgressItem]string
	RunningMsg, DoneMsg string
}

func NewProgressInfo(barType BarType) *ProgressInfo {
	p := new(ProgressInfo)

	p.Totals.RWMutex = new(sync.RWMutex)
	p.StartChan = make(chan *ProgressStatusStart)
	p.UpdateChan = make(chan *ProgressStatusUpdate)
	p.CompleteChan = make(chan *ProgressStatusComplete)
	p.ErrorChan = make(chan *ProgressStatusError)

	p.Progress = uiprogress.New()
	p.Progress.RefreshInterval = time.Second * 1
	p.BarType = barType
	p.BarsByName = make(map[string]*ProgressItem, 0)
	p.NamesByBar = make(map[*ProgressItem]string, 0)

	p.AddBar(2).PrependFunc(func(b *uiprogress.Bar) string {
		p.Totals.Lock()
		defer p.Totals.Unlock()
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", p.Totals.Active, p.Totals.Completed, p.Totals.Errored)
	}).PrependElapsed()

	p.SummaryBar = &ProgressBarText{&ProgressBar{p.Bars[0]}}
	p.SummaryBar.setBarToText()

	return p
}

func (p *ProgressInfo) Listen() {
	for {
		select {
		case s := <-p.StartChan:
			p.StartBar(s)
		case s := <-p.UpdateChan:
			p.UpdateBar(s)
		case s := <-p.CompleteChan:
			p.CompleteBar(s)
		case s := <-p.ErrorChan:
			p.ErrorBar(s)
		}
		p.Update()
	}
}

func (p *ProgressInfo) Update() {
	p.Bars[0].Incr()
	p.Bars[0].Set(0)
}

func (p *ProgressInfo) StartBar(status *ProgressStatusStart) {
	statusBarInfo := p.BarsByName[status.Name]
	switch statusBarInfo {
	case nil:
		var bar ProgressBarrer
		switch p.BarType {
		case 0:
			statusBar := p.AddBar(100).PrependElapsed().AppendCompleted().PrependFunc(func(b *uiprogress.Bar) string {
				return status.Name
			})
			bar = &ProgressBarPercentage{&ProgressBar{statusBar}}
		case 1:
			statusBar := p.AddBar(status.TotalSize).PrependElapsed().AppendCompleted().PrependFunc(func(b *uiprogress.Bar) string {
				return status.Name
			})
			bar = &ProgressBarBytes{&ProgressBar{statusBar}}
		default:
			statusBar := p.AddBar(2).PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
				var msg string
				switch b.Current() {
				case b.Total:
					msg = p.DoneMsg
				default:
					msg = p.RunningMsg
				}
				return fmt.Sprintf("%s\t%s", status.Name, msg)
			})
			pt := &ProgressBarText{&ProgressBar{statusBar}}
			pt.setBarToText()
			bar = pt
		}

		pi := &ProgressItem{bar}
		p.BarsByName[status.Name] = pi
		p.NamesByBar[pi] = status.Name
		p.Totals.Lock()
		p.Totals.Active++
		p.Totals.Unlock()
	default:
		p.Totals.Lock()
		p.Totals.Active++
		p.Totals.Errored--
		p.Totals.Unlock()
	}
}

func (p *ProgressInfo) UpdateBar(status *ProgressStatusUpdate) {
	if statusBarInfo := p.BarsByName[status.Name]; statusBarInfo != nil {
		statusBarInfo.Content.Update(status)
	}
}

func (p *ProgressInfo) CompleteBar(status *ProgressStatusComplete) {
	if statusBarInfo := p.BarsByName[status.Name]; statusBarInfo != nil {
		statusBarInfo.Content.Complete(status)
		p.Totals.Lock()
		p.Totals.Active--
		p.Totals.Completed++
		p.Totals.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (p *ProgressInfo) ErrorBar(status *ProgressStatusError) {
	if statusBarInfo := p.BarsByName[status.Name]; statusBarInfo != nil {
		p.NamesByBar[statusBarInfo] = statusBarInfo.Content.Error(status)
		p.Totals.Lock()
		p.Totals.Active--
		p.Totals.Errored++
		p.Totals.Unlock()
	}
}

type ProgressBarrer interface {
	Complete(*ProgressStatusComplete)
	Update(*ProgressStatusUpdate)
	Error(*ProgressStatusError) string
}

type ProgressBar struct {
	*uiprogress.Bar
}

func (b *ProgressBar) Complete(_ *ProgressStatusComplete) {
	b.Set(b.Total)
}

func (b *ProgressBar) Error(status *ProgressStatusError) string {
	return fmt.Sprintf("[ERROR: %s] %s", status.Err, status.Name)
}

type ProgressBarPercentage struct {
	*ProgressBar
}

func (b *ProgressBarPercentage) Update(status *ProgressStatusUpdate) {
	b.Incr()
	b.Set(status.Increment - 1)
}

type ProgressBarBytes struct {
	*ProgressBar
}

func (b *ProgressBarBytes) Update(status *ProgressStatusUpdate) {
	b.Set(b.Current() + status.Increment)
}

type ProgressBarText struct {
	*ProgressBar
}

func (b *ProgressBarText) Update(status *ProgressStatusUpdate) {
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
