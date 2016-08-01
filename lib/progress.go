package lib

import "time"

type ProgressStatuser interface {
	Error() error
	TimeElapsed() time.Duration
	//PercentComplete() int
}

type ProgressContenter interface {
	Update(ProgressStatuser)
	Complete(ProgressStatuser)
	Error(ProgressStatuser) string
}

type Progresser interface {
	Commander
	InitProgress()
	StartBar(ProgressStatuser)
	UpdateBar(ProgressStatuser)
	CompleteBar(ProgressStatuser)
	ErrorBar(ProgressStatuser)
}
