package interfaces

import (
	"github.com/gophercloud/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

type ServiceClientFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)

type Commander interface {
	Execute(item interface{}, out chan interface{})
	SetServiceClient(*gophercloud.ServiceClient)
	ServiceClient() *gophercloud.ServiceClient
	SetContext(*cli.Context) error
	Context() *cli.Context
	ServiceClientFunc() ServiceClientFunc
	ServiceType() string
	ServiceVersion() string
}

type Flagser interface {
	Flags() []cli.Flag
	HandleFlags() error
}

type Singler interface {
	HandleSingle() (interface{}, error)
}

type PipeCommander interface {
	Commander
	Singler
	HandlePipe(string) (interface{}, error)
	PipeFieldOptions() []string
	PipeFlags() []cli.Flag
	SetConcurrency(int)
	Concurrency() int
}

type StreamCommander interface {
	Commander
	Singler
	HandleStream() (interface{}, error)
	StreamFieldOptions() []string
}
