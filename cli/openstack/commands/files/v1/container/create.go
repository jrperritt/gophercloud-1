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

type CommandCreate struct {
	ContainerV1Command
	traits.Waitable
	opts containers.CreateOptsBuilder
}

var (
	cCreate                          = new(CommandCreate)
	_       interfaces.PipeCommander = cCreate

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "[--name <name> | --stdin name]"),
	Description:  "Creates a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
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
			Name:  "metadata",
			Usage: "[optional] Comma-separated key-value pairs for the container. Example: key1=val1,key2=val2",
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

func (c *CommandCreate) Fields() []string {
	return []string{""}
}

func (c *CommandCreate) HandleFlags() error {
	opts := &containers.CreateOpts{
		ContainerRead:  c.Context().String("container-read"),
		ContainerWrite: c.Context().String("container-write"),
	}

	if c.Context().IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	c.opts = opts

	return nil
}

func (c *CommandCreate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandCreate) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandCreate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := containers.Create(c.ServiceClient, item.(string), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully created container [%s]", item.(string))
	default:
		out <- err
	}
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}
