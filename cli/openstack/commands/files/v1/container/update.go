package container

import (
	"fmt"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"gopkg.in/urfave/cli.v1"
)

type CommandUpdate struct {
	ContainerV1Command
	traits.Waitable
	name string
	opts containers.UpdateOptsBuilder
}

var (
	cUpdate                          = new(CommandUpdate)
	_       interfaces.PipeCommander = cUpdate

	flagsUpdate = openstack.CommandFlags(cUpdate)
)

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", "--name <containerName>"),
	Description:  "Updates a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpdate) },
	Flags:        flagsUpdate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpdate) },
}

func (c *CommandUpdate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.StringFlag{
			Name:  "container-read",
			Usage: "[optional] Comma-separated list of users for whom to grant read access to the container",
		},
		cli.StringFlag{
			Name:  "container-write",
			Usage: "[optional] Comma-separated list of users for whom to grant write access to the container",
		},
	}
}

func (c *CommandUpdate) HandleFlags() (err error) {
	c.opts = &containers.UpdateOpts{
		ContainerRead:  c.Context.String("container-read"),
		ContainerWrite: c.Context.String("container-write"),
	}
	return
}

func (c *CommandUpdate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandUpdate) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandUpdate) Execute(item interface{}, out chan interface{}) {
	name := item.(string)
	r := containers.Update(c.ServiceClient, name, c.opts)
	switch r.Err {
	case nil:
		out <- fmt.Sprintf("Successfully updated container [%s]", name)
	default:
		out <- r.Err
	}
}

func (c *CommandUpdate) PipeFieldOptions() []string {
	return []string{"name"}
}
