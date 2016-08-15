package container

import (
	"fmt"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"gopkg.in/urfave/cli.v1"
)

type commandUpdate struct {
	openstack.CommandUtil
	ContainerV1Command
	name string
	opts containers.UpdateOptsBuilder
}

var (
	cUpdate                   = new(commandUpdate)
	_       lib.PipeCommander = cUpdate
	_       lib.Waiter        = cUpdate

	flagsUpdate = openstack.CommandFlags(cUpdate)
)

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", "--name <containerName>"),
	Description:  "Updates a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpdate) },
	Flags:        flagsUpdate,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsUpdate) },
}

func (c *commandUpdate) Flags() []cli.Flag {
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

func (c *commandUpdate) HandleFlags() (err error) {
	c.opts = &containers.UpdateOpts{
		ContainerRead:  c.Context.String("container-read"),
		ContainerWrite: c.Context.String("container-write"),
	}
	return
}

func (c *commandUpdate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandUpdate) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandUpdate) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		name := item.(string)
		r := containers.Update(c.ServiceClient, name, c.opts)
		switch r.Err {
		case nil:
			out <- fmt.Sprintf("Successfully updated container [%s]\n", name)
		default:
			out <- r.Err
		}
	}
}

func (c *commandUpdate) PipeFieldOptions() []string {
	return []string{"name"}
}

func (c *commandUpdate) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}
