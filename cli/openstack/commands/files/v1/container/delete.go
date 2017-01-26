package container

import (
	"fmt"

	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	ContainerV1Command
	traits.Waitable
	purge bool
}

var (
	cDelete                          = new(CommandDelete)
	_       interfaces.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--name <NAME> | --stdin name]"),
	Description:  "Deletes a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.BoolFlag{
			Name:  "purge",
			Usage: "[optional] Delete all objects in the container, and then delete the container.",
		},
	}
}

func (c *CommandDelete) Fields() []string {
	return []string{""}
}

func (c *CommandDelete) HandleFlags() error {
	c.purge = c.Context.IsSet("purge")
	return nil
}

func (c *CommandDelete) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	if c.purge {
		err := handleEmpty(cDelete.ContainerV1Command, item.(string))
		if err != nil {
			out <- fmt.Errorf("Error purging container [%s]: %s", item.(string), err)
			return
		}
	}
	res := containers.Delete(c.ServiceClient, item.(string))
	switch res.Err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted container [%s]", item.(string))
	default:
		out <- res.Err
	}
}

func (c *CommandDelete) PipeFieldOptions() []string {
	return []string{"name"}
}
