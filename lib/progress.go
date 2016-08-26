package lib

import (
	"time"

	cli "gopkg.in/urfave/cli.v1"
)

type Waiter interface {
	Commander
	ShouldWait() bool
	WaitFlags() []cli.Flag
	ExecuteAndWait(in, out chan interface{})
}

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
	Waiter
	InitProgress()
	ShowProgress(item interface{}, out chan interface{})
	EndProgress()
	ShouldProgress() bool
	ProgressFlags() []cli.Flag
}
