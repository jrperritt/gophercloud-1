package lib

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/gophercloud"
)

// Commander is an interface that all commands implement.
type Commander interface {
	//Flags() []cli.Flag
	SetServiceClient(*gophercloud.ServiceClient)
	// ServiceClientType returns the type of the service client to use.
	ServiceClientType() string
	// HandleFlags processes flags for the command that are relevant for both piped
	// and non-piped commands.
	HandleFlags() error
	// Execute executes the command's HTTP request.
	Execute(interface{}, chan interface{})
	Ctx() *cli.Context
}

// PipeCommander is an interface that commands implement if they can accept input
// from STDIN.
type PipeCommander interface {
	// Commander is an interface that all commands will implement.
	Commander
	// HandleSingle contains logic for processing a single resource. This method
	// will be used if input isn't sent to STDIN, so it will contain, for example,
	// logic for handling flags that would be mandatory if otherwise not piped in.
	HandleSingle() (interface{}, error)
	// HandlePipe is a method that commands implement for processing piped input.
	HandlePipe(string) (interface{}, error)
	// StdinFieldOptions is a slice of the fields that the command accepts on STDIN.
	PipeFieldOptions() []string
}

// StreamPipeCommander is an interface that commands implement if they can stream input
// from STDIN.
type StreamPipeCommander interface {
	// PipeHandler is an interface that commands implement if they can accept input
	// from STDIN.
	PipeCommander
	// HandleStreamPipe is a method that commands implement for processing streaming, piped input.
	HandleStreamPipe(string) (interface{}, error)
	// StreamFieldOptions is a slice of the fields that the command accepts for streaming input on STDIN.
	StreamFieldOptions() []string
}

type CommandInfoer interface {
	CommandInfo() string
}

type Fieldser interface {
	Fields() []string
}
