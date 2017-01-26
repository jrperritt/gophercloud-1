package object

import (
	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandGet struct {
	ObjectV1Command
	traits.Waitable
}

var (
	cGet                          = new(commandGet)
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "--container <containerName> [--name <objectName> | --stdin name]"),
	Description:  "Gets an object's metadata",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *commandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the object.",
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
	err := c.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}
	c.container = c.Context.String("container")
	return nil
}

func (c *commandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandGet) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandGet) Execute(item interface{}, out chan interface{}) {
	name := item.(string)
	var m map[string]interface{}
	err := objects.Get(c.ServiceClient, c.container, name, nil).ExtractInto(&m)
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}

func (c *commandGet) PipeFieldOptions() []string {
	return []string{"name"}
}
