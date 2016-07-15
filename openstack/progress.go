package openstack

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosuri/uiprogress"
)

// StatusMsg is the type of status message being sent on the channel returned by
// ChannelOfStatuses.
type StatusMsg string

var (
	StatusStarted   StatusMsg = "started"
	StatusUpdated   StatusMsg = "update"
	StatusCompleted StatusMsg = "success"
	StatusErrored   StatusMsg = "error"
)

type ProgressSummary struct {
	*uiprogress.Progress
	TotalActive, TotalCompleted, TotalErrored int
	StatusBarsByName                          map[string]*ProgressBarInfo
	FileNamesByBar                            map[*uiprogress.Bar]string
	SummaryBar                                *uiprogress.Bar
	WaitGroup                                 *sync.WaitGroup
}

type ProgressBarInfo struct {
	Index int
	Bar   *uiprogress.Bar
}

func setBarToText(b *uiprogress.Bar) {
	b.LeftEnd = ' '
	b.RightEnd = ' '
	b.Head = ' '
	b.Fill = ' '
	b.Empty = ' '
}

func NewTextBar(s string) *uiprogress.Bar {
	b := new(uiprogress.Bar).PrependFunc(func(b *uiprogress.Bar) string {
		return s
	})
	setBarToText(b)
	return b
}

func NewProgressSummary() *ProgressSummary {
	ps := &ProgressSummary{
		WaitGroup:        new(sync.WaitGroup),
		Progress:         uiprogress.New(),
		StatusBarsByName: make(map[string]*ProgressBarInfo, 0),
		FileNamesByBar:   make(map[*uiprogress.Bar]string, 0),
	}
	ps.Progress.RefreshInterval = time.Second * 1
	ps.AddBar(2).PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrored: %d", ps.TotalActive, ps.TotalCompleted, ps.TotalErrored)
	}).PrependElapsed()
	setBarToText(ps.Bars[0])
	return ps
}

func (ps *ProgressSummary) Update() {
	ps.Bars[0].Incr()
	ps.Bars[0].Set(0)
}

type ProgressStatus struct {
	Name      string
	TotalSize int
	Increment int
	StartTime time.Time
	MsgType   StatusMsg
	Err       error
	Result    interface{}
}

func (ps ProgressStatus) Error() error {
	return ps.Err
}

func (ps ProgressStatus) TimeElapsed() time.Duration {
	return time.Since(ps.StartTime)
}

func (ps ProgressStatus) PercentComplete() int {
	return 0
}
