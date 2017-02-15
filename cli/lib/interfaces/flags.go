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
	FieldsFlags() []cli.Flag
	SetFields([]string)
	Fields() []string
}

// Tabler is the interface a command implements if it offers tabular output.
// `TableFlags` and `ShouldHeader` are common to all `Tabler`s, so a command
// need only have `DefaultTableFields` method
type Tabler interface {
	//Fieldser
	TableFlags() []cli.Flag
	DefaultTableFields() []string
	SetTable(bool)
	ShouldTable() bool
	SetHeader(bool)
	ShouldHeader() bool
}
