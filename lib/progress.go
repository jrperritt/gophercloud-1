package lib

import (
	"time"

	"gopkg.in/urfave/cli.v1"
)

type ProgressStatuser interface {
	Error() error
	TimeElapsed() time.Duration
	ID() string
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
	ShowProgress(item interface{}, out chan interface{})
	End()
}

type Waiter interface {
	Commander
	ShouldWait() bool
	WaitFlags() []cli.Flag
	ShouldQuiet() bool
	ExecuteAndWait(in, out chan interface{})
}
