package lib

import "time"

/*
type StatusUpdater interface {
	Commander
	StatusChannel(*Resourcer) (chan interface{}, chan interface{}, error)
}

type ProgressBarUpdater interface {
	StatusUpdater
	NewProgress(chan interface{}, chan interface{}) Progresser
}
*/

type ProgressStatuser interface {
	Error() error
	TimeElapsed() time.Duration
	PercentComplete() int
}

type Progresser interface {
	Commander
	StatusChannel() (chan ProgressStatuser, error)
	UpdateSummary()
	UpdateProgress()
}
