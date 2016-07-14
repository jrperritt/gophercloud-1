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
	InitProgress()
	Started(interface{})
	Updated(interface{})
	Completed(interface{})
	Errored(interface{})
}
