package traits

import (
	"fmt"
	"sync"
	"time"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gosuri/uiprogress"
)

type ProgressInfo struct {
	Totals struct {
		*sync.RWMutex
		Active    int
		Completed int
		Errored   int
	}
	*uiprogress.Progress
	SummaryBar          *ProgressBarText
	BarType             BarType
	BarsByName          map[string]interfaces.ProgressBarrer
	NamesByBar          map[interfaces.ProgressBarrer]string
	RunningMsg, DoneMsg string
}

func NewProgressInfo(barType BarType) *ProgressInfo {
	p := new(ProgressInfo)

	p.Totals.RWMutex = new(sync.RWMutex)

	p.Progress = uiprogress.New()
	p.Progress.RefreshInterval = time.Second * 1
	p.BarType = barType
	p.BarsByName = make(map[string]interfaces.ProgressBarrer, 0)
	p.NamesByBar = make(map[interfaces.ProgressBarrer]string, 0)

	p.AddBar(2).PrependFunc(func(b *uiprogress.Bar) string {
		p.Totals.Lock()
		defer p.Totals.Unlock()
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", p.Totals.Active, p.Totals.Completed, p.Totals.Errored)
	}).PrependElapsed()

	b := new(ProgressBarText)
	b.ProgressBar = new(ProgressBar)
	b.ProgressBar.Bar = p.Bars[0]
	b.id = "summary"
	p.SummaryBar = b
	p.SummaryBar.setBarToText()

	return p
}

func (p *ProgressInfo) Listen(statusch chan interfaces.ProgressStatuser) {
	for s := range statusch {
		switch s.(type) {
		case *ProgressStatusStart:
			//lib.Log.Debugln("recvd ProgressStatusStart")
			p.StartBar(s.(*ProgressStatusStart))
		case *ProgressStatusUpdate:
			//lib.Log.Debugln("recvd ProgressStatusUpdate")
			p.UpdateBar(s.(*ProgressStatusUpdate))
		case *ProgressStatusComplete:
			//lib.Log.Debugln("recvd ProgressStatusComplete")
			p.CompleteBar(s.(*ProgressStatusComplete))
		case *ProgressStatusError:
			//lib.Log.Debugln("recvd ProgressStatusError")
			p.ErrorBar(s.(*ProgressStatusError))
		}
		p.Update()
	}
}

func (p *ProgressInfo) Update() {
	p.Bars[0].Incr()
	p.Bars[0].Set(0)
}

func (p *ProgressInfo) StartBar(status *ProgressStatusStart) {
	if statusBarInfo, ok := p.BarsByName[status.Name]; ok {
		p.Totals.Lock()
		p.Totals.Active++
		p.Totals.Errored--
		p.Totals.Unlock()
		statusBarInfo.Reset()
		return
	}

	var bar interfaces.ProgressBarrer
	switch p.BarType {
	case 0:
		statusBar := p.AddBar(100).PrependElapsed().AppendCompleted().PrependFunc(func(b *uiprogress.Bar) string {
			return status.Name
		})
		b := new(ProgressBarPercentage)
		b.ProgressBar = new(ProgressBar)
		b.ProgressBar.Bar = statusBar
		b.id = status.Name
		bar = b
		//bar = &ProgressBarPercentage{&ProgressBar{statusBar}}
	case 1:
		statusBar := p.AddBar(status.TotalSize).PrependElapsed().AppendCompleted().PrependFunc(func(b *uiprogress.Bar) string {
			return status.Name
		})
		b := new(ProgressBarBytes)
		b.ProgressBar = new(ProgressBar)
		b.ProgressBar.Bar = statusBar
		b.id = status.Name
		bar = b
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
		b := new(ProgressBarText)
		b.ProgressBar = new(ProgressBar)
		b.ProgressBar.Bar = statusBar
		b.id = status.Name
		b.setBarToText()
		bar = b
	}

	p.BarsByName[status.Name] = bar
	p.NamesByBar[bar] = status.Name
	p.Totals.Lock()
	p.Totals.Active++
	p.Totals.Unlock()
}

func (p *ProgressInfo) UpdateBar(status *ProgressStatusUpdate) {
	if statusBarInfo := p.BarsByName[status.Name]; statusBarInfo != nil {
		statusBarInfo.Update(status)
	}
}

func (p *ProgressInfo) CompleteBar(status *ProgressStatusComplete) {
	if statusBarInfo := p.BarsByName[status.Name]; statusBarInfo != nil {
		statusBarInfo.Complete()
		p.Totals.Lock()
		p.Totals.Active--
		p.Totals.Completed++
		p.Totals.Unlock()
	}
}

func (p *ProgressInfo) ErrorBar(status *ProgressStatusError) {
	if statusBarInfo := p.BarsByName[status.Name]; statusBarInfo != nil {
		p.NamesByBar[statusBarInfo] = statusBarInfo.Error(status.Err)
		p.Totals.Lock()
		p.Totals.Active--
		p.Totals.Errored++
		p.Totals.Unlock()
	}
}
