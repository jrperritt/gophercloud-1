package interfaces

import (
	"io"

	cli "gopkg.in/urfave/cli.v1"
)

// Progresser should be implemented by commands that allow progress updates
// during execution
type Progresser interface {
	SetProgress(bool)
	ShouldProgress() bool
	ProgressFlags() []cli.Flag

	InitProgress()
	AddSummaryBar()
	ProgStartCh() chan ProgressItemer

	CreateBar(ProgressItemer) ProgressBarrer
	StartBar()
	CompleteBar()
	ErrorBar()
}

type ProgressItemer interface {
	SetUpCh(chan interface{})
	UpCh() chan interface{}
	SetEndCh(chan interface{})
	EndCh() chan interface{}
	SetID(string)
	ID() string
	Size() int64
}

type BytesProgressItemer interface {
	ProgressItemer
	SetSize(int64)
}

type ReadBytesProgressItemer interface {
	BytesProgressItemer
	SetReader(io.Reader)
	Reader() io.Reader
}

type WriteBytesProgressItemer interface {
	BytesProgressItemer
	SetWriter(io.Writer)
	Writer() io.Writer
}

type BytesProgresser interface {
	Progresser
}

type PercentProgresser interface {
	Progresser
}

type TextProgresser interface {
	Progresser
	RunningMsg() string
	DoneMsg() string
}

type ProgressBarrer interface {
	Start(ProgressStartStatuser) ProgressBarrer
	Complete(ProgressCompleteStatuser)
	Update(ProgressUpdateStatuser)
	Error(ProgressErrorStatuser) string
	ID() string
}

type ProgressStatuser interface {
	SetBarID(barid string)
	BarID() string
}

type ProgressStartStatuser interface {
	ProgressStatuser
	BarSize() int64
	SetBarSize(size int64)
}

type ProgressUpdateStatuser interface {
	ProgressStatuser
	Change() interface{}
	SetChange(change interface{})
}

type ProgressCompleteStatuser interface {
	ProgressStatuser
	Result() interface{}
	SetResult(result interface{})
}

type ProgressErrorStatuser interface {
	ProgressStatuser
	Err() error
	SetErr(err error)
}
