package interfaces

import (
	"io"
	"sync"

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
	Context() *cli.Context
	ServiceClientFunc() ServiceClientFunc
	ServiceType() string
	ServiceVersion() string
	ExecDoneCh() chan interface{}
	AllDoneCh() chan interface{}
	UpdateCh() chan interface{}
	WG() *sync.WaitGroup
}

type PipeCommander interface {
	Commander
	HandleSingle() (interface{}, error)
	HandlePipe(string) (interface{}, error)
	PipeFieldOptions() []string
}

type StreamPipeCommander interface {
	PipeCommander
	HandleStreamPipe(io.Reader) (interface{}, error)
	StreamFieldOptions() []string
}
