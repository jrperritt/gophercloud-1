package interfaces

import (
	"strings"

	"github.com/gophercloud/gophercloud/cli/lib"

	"gopkg.in/urfave/cli.v1"
)

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
	InitProgress(donechout chan interface{})
	ProgDoneChIn() chan interface{}
	ProgDoneChOut() chan interface{}
	ProgUpdateChIn() chan interface{}
	ProgListenCh() chan ProgressStatuser
	BarID(item interface{}) string
	ShowBar(id string)
	SetProgress(bool)
	ShouldProgress() bool
	ProgressFlags() []cli.Flag
	Update()

	StartBar(ProgressStartStatuser)
	CompleteBar(ProgressCompleteStatuser)
	UpdateBar(ProgressUpdateStatuser)
	ErrorBar(ProgressErrorStatuser)
}

type BytesProgresser interface {
	Progresser
	BarSizes() map[string]int64
}

type ProgressBarrer interface {
	Start(ProgressStartStatuser) ProgressBarrer
	Complete(ProgressCompleteStatuser)
	Update(ProgressUpdateStatuser)
	Error(ProgressErrorStatuser) string
	Reset()
	ID() string
}

type ProgressBytesBarrer interface {
	ProgressBarrer
	TotalSize()
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

func HandleInterfaceFlags(cmd Commander) error {
	if w, ok := cmd.(Waiter); ok {
		lib.Log.Debugln("cmd implements Waiter")
		w.SetWait(cmd.Context().IsSet("wait"))
		lib.Log.Debugf("w.ShouldWait(): %+v\n", w.ShouldWait())
	}

	if p, ok := cmd.(Progresser); ok {
		lib.Log.Debugln("cmd implements Progresser")
		p.SetProgress(cmd.Context().IsSet("quiet"))
		lib.Log.Debugln("p.ShouldProgress() : ", p.ShouldProgress())
	}

	if t, ok := cmd.(Tabler); ok {
		lib.Log.Debugln("cmd implements Tabler")
		t.SetTable(cmd.Context().IsSet("table"))
		t.SetHeader(cmd.Context().IsSet("no-header"))
	}

	if f, ok := cmd.(Fieldser); ok {
		lib.Log.Debugln("cmd implements Fieldser")
		f.SetFields(strings.Split(cmd.Context().String("fields"), ","))
	}

	return nil
}
