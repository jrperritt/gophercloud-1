package openstack

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uiprogress"
)

// StatusMsg is the type of status message being sent on the channel returned by
// ChannelOfStatuses.
type StatusMsg string

var (
	StatusStarted StatusMsg = "started"
	StatusUpdate  StatusMsg = "update"
	StatusSuccess StatusMsg = "success"
	StatusError   StatusMsg = "error"
)

type ProgressStatus struct {
	Name      string
	TotalSize int
	Increment int
	StartTime time.Time
	MsgType   StatusMsg
	Err       error
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

type Progress struct {
	*uiprogress.Progress
	statusChannel chan ProgressStatus
}

func (pb *Progress) UpdateSummary() {
	pb.Bars[0].Incr()
	pb.Bars[0].Set(pb.Bars[0].Current() - 1)
}

func (pb *Progress) UpdateProgress() {
	statusBarsByName := map[string]*ProgressBarInfo{}
	fileNamesByBar := map[*uiprogress.Bar]string{}

	totalActive, totalCompleted, totalErrored := 0, 0, 0

	summaryBar := pb.AddBar(2).PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("\tActive: %d\tCompleted: %d\tErrorred: %d", totalActive, totalCompleted, totalErrored)
	}).PrependElapsed()
	summaryBar.LeftEnd = ' '
	summaryBar.RightEnd = ' '
	summaryBar.Head = ' '
	summaryBar.Fill = ' '
	summaryBar.Empty = ' '

	pb.Start()

	for status := range pb.statusChannel {
		//status := status.(*TransferStatus)
		switch status.MsgType {
		case StatusStarted:
			statusBarInfo := statusBarsByName[status.Name]
			if statusBarInfo == nil {
				statusBar := pb.AddBar(status.TotalSize).AppendCompleted().PrependElapsed().PrependFunc(func(b *uiprogress.Bar) string {
					return fileNamesByBar[b]
				}).AppendFunc(func(b *uiprogress.Bar) string {
					return fmt.Sprintf("%s/%s", humanize.Bytes(uint64(b.Current())), humanize.Bytes(uint64(b.Total)))
				})
				index := len(pb.Bars) - 1
				statusBarsByName[status.Name] = &ProgressBarInfo{index, statusBar}
				fileNamesByBar[statusBar] = status.Name
				totalActive++
				pb.UpdateSummary()
			} else {
				totalActive++
				totalErrored--
				pb.UpdateSummary()
			}

		case StatusUpdate:
			if statusBarInfo := statusBarsByName[status.Name]; statusBarInfo != nil {
				statusBarInfo.bar.Incr()
				statusBarInfo.bar.Set(statusBarInfo.bar.Current() - 1 + status.Increment)
				pb.UpdateSummary()
			}
		case StatusSuccess:
			if statusBarInfo := statusBarsByName[status.Name]; statusBarInfo != nil {
				statusBarInfo.bar.Set(status.TotalSize)
				delete(fileNamesByBar, statusBarInfo.bar)
				delete(statusBarsByName, status.Name)
				pb.Bars = append(pb.Bars[:statusBarInfo.index], pb.Bars[statusBarInfo.index+1:]...)
				for i, progressBar := range pb.Bars {
					if i != 0 {
						statusBarsByName[fileNamesByBar[progressBar]].index = i
					}
				}
				totalActive--
				totalCompleted++
				pb.UpdateSummary()
			}
		case StatusError:
			if statusBarInfo := statusBarsByName[status.Name]; statusBarInfo != nil {
				fileNamesByBar[statusBarInfo.bar] = fmt.Sprintf("[ERROR: %s, WILL RETRY] %s", status.Err, status.Name)
				totalActive--
				totalErrored++
				pb.UpdateSummary()
			}
		default:
			pb.statusChannel <- ProgressStatus{Err: status.Err}
		}
	}
	close(pb.statusChannel)
}

type ProgressBarInfo struct {
	index int
	bar   *uiprogress.Bar
}

func NewProgressBar(statusCh chan ProgressStatus) *Progress {
	p := uiprogress.New()
	p.RefreshInterval = time.Second * 1
	return &Progress{
		Progress:      p,
		statusChannel: statusCh,
		//gophercloudChannel: gophercloudCh,
	}
}
