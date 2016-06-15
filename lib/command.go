package lib

import (
	"reflect"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/gophercloud"
)

// Commander is an interface that all commands implement.
type Commander interface {
	Name() string
	Usage() string
	Description() string
	Action(*cli.Context)
	Flags() []cli.Flag
	SetFlags([]cli.Flag)
	BashComplete(*cli.Context)
	// Fields returns the fields available for the command output.
	Fields() []string
	SetFields([]string)

	ServiceClient() *gophercloud.ServiceClient
	SetServiceClient(*gophercloud.ServiceClient)
	// ServiceClientType returns the type of the service client to use.
	ServiceClientType() string
	SetServiceClientType(string)
	// HandleFlags processes flags for the command that are relevant for both piped
	// and non-piped commands.
	HandleFlags() error
	RunCommand(chan Resulter) error
	// Execute executes the command's HTTP request.
	Execute(Resourcer) Resulter

	ReturnType() reflect.Type
	SetDebugChannel(chan string)
}

// PipeCommander is an interface that commands implement if they can accept input
// from STDIN.
type PipeCommander interface {
	// Commander is an interface that all commands will implement.
	Commander
	// HandleSingle contains logic for processing a single resource. This method
	// will be used if input isn't sent to STDIN, so it will contain, for example,
	// logic for handling flags that would be mandatory if otherwise not piped in.
	HandleSingle(Resourcer) error
	// HandlePipe is a method that commands implement for processing piped input.
	//HandlePipe(string) error
	// StdinFieldOptions is a slice of the fields that the command accepts on STDIN.
	PipeFieldOptions() []string
	PipeField() string
	//SetPipeField(string)
}

// StreamPipeCommander is an interface that commands implement if they can stream input
// from STDIN.
type StreamPipeCommander interface {
	// PipeHandler is an interface that commands implement if they can accept input
	// from STDIN.
	PipeCommander
	// HandleStreamPipe is a method that commands implement for processing streaming, piped input.
	//HandleStreamPipe() error
	// StreamFieldOptions is a slice of the fields that the command accepts for streaming input on STDIN.
	StreamFieldOptions() []string
}

type CommandInfoer interface {
	CommandInfo() string
}
