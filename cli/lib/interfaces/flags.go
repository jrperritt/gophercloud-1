package interfaces

import "gopkg.in/urfave/cli.v1"

// Waiter should be implemented by commands that launch background operations
// that will continue even if the command ends
type Waiter interface {
	WaitFor(item interface{}, out chan<- interface{})
	SetWait(bool)
	ShouldWait() bool
	WaitFlags() []cli.Flag
}

// Fieldser should be implemented by commands that return fields in the output
type Fieldser interface {
	SetFields([]string)
	Fields() []string
}

// Progresser should be implemented by commands that allow progress updates
// during execution
type Progresser interface {
	SetProgress(bool)
	ShouldProgress() bool
	ProgressFlags() []cli.Flag

	InitProgress()
	ProgUpdateCh() chan interface{}
	ProgStartCh() chan ProgressItemer

	CreateBar(ProgressItemer) ProgressBarrer
	StartBar()
	CompleteBar()
	ErrorBar()
}

type ProgressItemer interface {
	UpCh() chan interface{}
	SetEndCh(chan interface{})
	EndCh() chan interface{}
	ID() string
	Size() int64
}

type ReadBytesProgresser interface {
	Progresser
}

type WriteBytesProgresser interface {
	Progresser
}

type BytesProgresser interface {
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
	TotalSize() int64
}

type ProgressBytesBarrer interface {
	ProgressBarrer
}

// Tabler is the interface a command implements if it offers tabular output.
// `TableFlags` and `ShouldHeader` are common to all `Tabler`s, so a command
// need only have `DefaultTableFields` method
type Tabler interface {
	TableFlags() []cli.Flag
	DefaultTableFields() []string
	SetTable(bool)
	ShouldTable() bool
	SetHeader(bool)
	ShouldHeader() bool
}
