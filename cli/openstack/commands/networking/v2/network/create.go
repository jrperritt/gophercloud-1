package network

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	NetworkV2Command
	traits.Pipeable
	traits.Waitable
	opts networks.CreateOptsBuilder
}

var (
	cCreate                          = new(CommandCreate)
	_       interfaces.PipeCommander = cCreate

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "[--name <name> | --stdin name]"),
	Description:  "Creates a network",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the network",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.BoolFlag{
			Name:  "up",
			Usage: "[optional] If provided, the network will be up upon creation.",
		},
		cli.BoolFlag{
			Name:  "shared",
			Usage: "[optional] If provided, the network is shared among all tenants.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant who should own this network.",
		},
	}
}

func (c *CommandCreate) Fields() []string {
	return []string{""}
}

func (c *CommandCreate) HandleFlags() error {
	opts := &networks.CreateOpts{
		TenantID: c.Context().String("tenant-id"),
	}

	if c.Context().IsSet("up") {
		opts.AdminStateUp = gophercloud.Enabled
	}

	if c.Context().IsSet("shared") {
		opts.Shared = gophercloud.Enabled
	}

	c.opts = opts

	return nil
}

func (c *CommandCreate) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandCreate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	opts := *c.opts.(*networks.CreateOpts)
	opts.Name = item.(string)
	err := networks.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["network"]
	default:
		out <- err
	}
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}
