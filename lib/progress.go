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
	ShowProgress(in, out chan interface{})
	StartBar(ProgressStatuser)
	UpdateBar(ProgressStatuser)
	CompleteBar(ProgressStatuser)
	ErrorBar(ProgressStatuser)
}

type Waiter interface {
	ShouldWait() bool
	ExecuteAndWait(in, out chan interface{})
}
