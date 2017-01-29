package interfaces

import (
	"io"

	"github.com/gophercloud/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

type ServiceClientFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)

type Commander interface {
	HandleInterfaceFlags() error
	HandleFlags() error
	Execute(item interface{}, out chan interface{})
	Flags() []cli.Flag
	SetServiceClient(*gophercloud.ServiceClient) error
	SetContext(*cli.Context) error
	ServiceClientFunc() ServiceClientFunc
	ServiceType() string
	ServiceVersion() string
}

type PipeCommander interface {
	Commander
	Waiter
	HandleSingle() (interface{}, error)
	HandlePipe(string) (interface{}, error)
	PipeFieldOptions() []string
}

type StreamPipeCommander interface {
	PipeCommander
	HandleStreamPipe(io.Reader) (interface{}, error)
	StreamFieldOptions() []string
}
