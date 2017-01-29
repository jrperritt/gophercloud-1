package interfaces

import "gopkg.in/urfave/cli.v1"

type Waiter interface {
	WaitFor(item interface{})
	SetWait(bool)
	ShouldWait() bool
	WaitFlags() []cli.Flag
}

type Fieldser interface {
	Fields() []string
}

type Progresser interface {
	Waiter
	InitProgress()
	BarID(item interface{}) string
	ShowBar(id string)
	SetProgress(bool)
	ShouldProgress() bool
	ProgressFlags() []cli.Flag
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
