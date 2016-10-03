package object

import (
	"fmt"

	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandDelete struct {
	ObjectV1Command
	traits.Waitable
	traits.Pipeable
}

var (
	cDelete                         = new(commandDelete)
	_       openstack.Waiter        = cDelete
	_       openstack.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "--container <containerName> [--name <objectName> | --stdin name]"),
	Description:  "Deletes an object",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *commandDelete) Flags() []cli.Flag {
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

func (c *commandDelete) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}
	c.container = c.Context.String("container")
	return nil
}

func (c *commandDelete) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandDelete) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := objects.Delete(c.ServiceClient, c.container, item.(string), nil).ExtractInto(&m)
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted object [%s] from container [%s]", item.(string), c.container)
	default:
		out <- err
	}
}

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"name"}
}
