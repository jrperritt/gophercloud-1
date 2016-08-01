package openstack

import (
	"fmt"
	"time"

	"github.com/gophercloud/cli/lib"
	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
)

type BarType uint8

const (
	BarPecentage BarType = iota
	BarBytes
	BarText
)

type Progress struct {
	*uiprogress.Progress
	SummaryBar                                *ProgressText
	BarType                                   BarType
	StatusBarsByName                          map[string]*ProgressItem
	FileNamesByBar                            map[*ProgressItem]string
	TotalActive, TotalCompleted, TotalErrored int
	RunningMsg, DoneMsg                       string
}

func NewProgress(barType BarType) *Progress {
	p := new(Progress)
	p.Progress = uiprogress.New()
	p.Progress.RefreshInterval = time.Second * 1
	p.BarType = barType
	p.StatusBarsByName = make(map[string]*ProgressItem, 0)
	p.FileNamesByBar = make(map[*ProgressItem]string, 0)

	p.AddBar(2).PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", p.TotalActive, p.TotalCompleted, p.TotalErrored)
	}).PrependElapsed()

	p.SummaryBar = &ProgressText{&ProgressBar{p.Bars[0]}}
	p.SummaryBar.setBarToText()

	return p
}

func (p *Progress) Update() {
	p.Bars[0].Incr()
	p.Bars[0].Set(0)
}

func (p *Progress) StartBar(raw lib.ProgressStatuser) {
	status := raw.(*ProgressStatus)
	statusBarInfo := p.StatusBarsByName[status.Name]
	switch statusBarInfo {
	case nil:
		var bar lib.ProgressContenter
		switch p.BarType {
		case 0:
			statusBar := p.AddBar(status.TotalSize).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
				return status.Name
			})
			bar = &ProgressBarPercentage{&ProgressBar{statusBar}}
		case 1:
			statusBar := p.AddBar(status.TotalSize).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
				return status.Name
			})
			bar = &ProgressBarBytes{&ProgressBar{statusBar}}
		default:
			statusBar := p.AddBar(2).PrependFunc(func(b *uiprogress.Bar) string {
				var msg string
				switch b.Current() {
				case b.Total:
					msg = p.DoneMsg
				default:
					msg = p.RunningMsg
				}
				return fmt.Sprintf("%s\t%s\t%s", status.Name, strutil.PrettyTime(time.Since(status.StartTime)), msg)
			})
			pt := &ProgressText{&ProgressBar{statusBar}}
			pt.setBarToText()
			bar = pt
		}

		pi := &ProgressItem{len(p.Bars) - 1, bar}
		p.StatusBarsByName[status.Name] = pi
		p.FileNamesByBar[pi] = status.Name
		p.TotalActive++
	default:
		p.TotalActive++
		p.TotalErrored--
	}
	p.Update()
}

func (p *Progress) UpdateBar(raw lib.ProgressStatuser) {
	status := raw.(*ProgressStatus)
	if statusBarInfo := p.StatusBarsByName[status.Name]; statusBarInfo != nil {
		statusBarInfo.Content.Update(raw)
		p.Update()
	}
}

func (p *Progress) CompleteBar(raw lib.ProgressStatuser) {
	status := raw.(*ProgressStatus)
	if statusBarInfo := p.StatusBarsByName[status.Name]; statusBarInfo != nil {
		statusBarInfo.Content.Complete(raw)
		p.TotalActive--
		p.TotalCompleted++
		p.Update()
		time.Sleep(1 * time.Second)
	}
}

func (p *Progress) ErrorBar(raw lib.ProgressStatuser) {
	status := raw.(*ProgressStatus)
	if statusBarInfo := p.StatusBarsByName[status.Name]; statusBarInfo != nil {
		p.FileNamesByBar[statusBarInfo] = statusBarInfo.Content.Error(raw)
		p.TotalActive--
		p.TotalErrored++
		p.Update()
	}
}

type ProgressBar struct {
	*uiprogress.Bar
}

func (pb *ProgressBar) Complete(_ lib.ProgressStatuser) {
	pb.Set(pb.Total)
}

func (pb *ProgressBar) Error(raw lib.ProgressStatuser) string {
	status := raw.(*ProgressStatus)
	return fmt.Sprintf("[ERROR: %s] %s", status.Err, status.Name)
}

type ProgressBarPercentage struct {
	*ProgressBar
}

func (pbp *ProgressBarPercentage) Update(raw lib.ProgressStatuser) {
	status := raw.(*ProgressStatus)
	pbp.Incr()
	pbp.Set(status.Increment - 1)
}

type ProgressBarBytes struct {
	*ProgressBar
}

func (pbb *ProgressBarBytes) Update(raw lib.ProgressStatuser) {
	status := raw.(*ProgressStatus)
	pbb.Set(pbb.Current() + status.Increment)
}

type ProgressText struct {
	*ProgressBar
}

func (pt *ProgressText) Update(raw lib.ProgressStatuser) {}

type ProgressItem struct {
	Index   int
	Content lib.ProgressContenter
}

func (pt *ProgressText) setBarToText() {
	pt.LeftEnd = ' '
	pt.RightEnd = ' '
	pt.Head = ' '
	pt.Fill = ' '
	pt.Empty = ' '
}

type ProgressStatus struct {
	Name      string
	TotalSize int
	Increment int
	StartTime time.Time
	Err       error
	Result    interface{}
}

func (ps ProgressStatus) Error() error {
	return ps.Err
}

func (ps ProgressStatus) TimeElapsed() time.Duration {
	return time.Since(ps.StartTime)
}
