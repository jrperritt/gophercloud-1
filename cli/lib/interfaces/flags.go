package interfaces

import "gopkg.in/urfave/cli.v1"

// Waiter should be implemented by commands that launch background operations
type Waiter interface {
	WaitFor(item interface{})
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
