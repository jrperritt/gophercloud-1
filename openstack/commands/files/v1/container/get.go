package container

import (
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"gopkg.in/urfave/cli.v1"
)

type commandGet struct {
	openstack.CommandUtil
	ContainerV1Command
}

var (
	cGet                   = new(commandGet)
	_    lib.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--name <containerName> | --stdin name]"),
	Description:  "Gets a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsGet) },
}

func (c *commandGet) Flags() []cli.Flag {
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

func (c *commandGet) Fields() []string {
	return []string{""}
}

func (c *commandGet) HandleFlags() error {
	return nil
}

func (c *commandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandGet) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandGet) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
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
}

func (c *commandGet) PipeFieldOptions() []string {
	return []string{"name"}
}
