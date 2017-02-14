package interfaces

import (
	"github.com/gophercloud/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

type ServiceClientFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)

type Commander interface {
	HandleFlags() error
	Execute(item interface{}, out chan interface{})
	Flags() []cli.Flag
	SetServiceClient(*gophercloud.ServiceClient)
	ServiceClient() *gophercloud.ServiceClient
	SetContext(*cli.Context) error
	Context() *cli.Context
	ServiceClientFunc() ServiceClientFunc
	ServiceType() string
	ServiceVersion() string
}

type PipeCommander interface {
	Commander
	HandleSingle() (interface{}, error)
	HandlePipe(string) (interface{}, error)
	PipeFieldOptions() []string
}
