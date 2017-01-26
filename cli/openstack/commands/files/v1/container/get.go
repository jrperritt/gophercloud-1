package container

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	ContainerV1Command
	traits.Waitable
}

var (
	cGet                          = new(CommandGet)
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--name <containerName> | --stdin name]"),
	Description:  "Gets a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

func (c *CommandGet) Fields() []string {
	return []string{""}
}

func (c *CommandGet) HandleFlags() error {
	return nil
}

func (c *CommandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandGet) Execute(item interface{}, out chan interface{}) {
	name := item.(string)
	var m map[string]interface{}
	err := containers.Get(c.ServiceClient, name).ExtractInto(&m)
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}

func (c *CommandGet) PipeFieldOptions() []string {
	return []string{"name"}
}
