package lib

// Commander is an interface that all commands implement.
type Commander interface {
	// Keys returns the keys available for the command output.
	Keys() []string
	// ServiceClientType returns the type of the service client to use.
	ServiceClientType() string
	// HandleFlags processes flags for the command that are relevant for both piped
	// and non-piped commands.
	HandleFlags() error
	// Execute executes the command's HTTP request.
	Execute(*Resourcer)
	ResultsChan() chan *Result
}

// PipeHandler is an interface that commands implement if they can accept input
// from STDIN.
type PipeHandler interface {
	// Commander is an interface that all commands will implement.
	Commander
	// HandleSingle contains logic for processing a single resource. This method
	// will be used if input isn't sent to STDIN, so it will contain, for example,
	// logic for handling flags that would be mandatory if otherwise not piped in.
	HandleSingle(*Resourcer) error
	// HandlePipe is a method that commands implement for processing piped input.
	HandlePipe(*Resourcer, string) error
	// StdinField is a slice of the fields that the command accepts on STDIN.
	StdinFields() []string
}

// StreamPipeHandler is an interface that commands implement if they can stream input
// from STDIN.
type StreamPipeHandler interface {
	// PipeHandler is an interface that commands implement if they can accept input
	// from STDIN.
	PipeHandler
	// StreamField is a slice of the fields that the command accepts for streaming input on STDIN.
	StreamFields() []string
	// HandleStreamPipe is a method that commands implement for processing streaming, piped input.
	HandleStreamPipe(*Resourcer) error
}

type Command struct {
	keys []string
}

func (c Command) Keys() []string {
	return c.keys
}
