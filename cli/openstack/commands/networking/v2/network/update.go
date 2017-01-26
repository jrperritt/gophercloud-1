package network

import (
	"fmt"

	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"gopkg.in/urfave/cli.v1"
)

type CommandUpdate struct {
	NetworkV2Command
	traits.Waitable
	opts networks.UpdateOptsBuilder
}

var (
	cUpdate                          = new(CommandUpdate)
	_       interfaces.PipeCommander = cUpdate

	flagsUpdate = openstack.CommandFlags(cUpdate)
)

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", "[--id <ID> | --name <NAME> | --stdin id]"),
	Description:  "Updates a network",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpdate) },
	Flags:        flagsUpdate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpdate) },
}

func (c *CommandUpdate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the network to update.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the network to update.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.StringFlag{
			Name:  "new-name",
			Usage: "[optional] The name that the network should have.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "[optional] If provided, the network will be up. Options are: true, false",
		},
		cli.StringFlag{
			Name:  "shared",
			Usage: "[optional] If provided, the network is shared among all tenants. Options are: true, false",
		},
	}
}

func (c *CommandUpdate) Fields() []string {
	return []string{""}
}

func (c *CommandUpdate) HandleFlags() error {
	opts := &networks.UpdateOpts{
		Name: c.Context.String("new-name"),
	}

	if c.Context.IsSet("up") {
		switch c.Context.String("up") {
		case "true":
			opts.AdminStateUp = gophercloud.Enabled
		case "false":
			opts.AdminStateUp = gophercloud.Disabled
		default:
			return fmt.Errorf("Invalid value for flag `up`: %s. Options are: true, false", c.Context.String("up"))
		}
	}

	if c.Context.IsSet("shared") {
		switch c.Context.String("shared") {
		case "true":
			opts.Shared = gophercloud.Enabled
		case "false":
			opts.Shared = gophercloud.Disabled
		default:
			return fmt.Errorf("Invalid value for flag `shared`: %s. Options are: true, false", c.Context.String("shared"))
		}
	}

	c.opts = opts
	c.Wait = c.Context.IsSet("wait")

	return nil
}

func (c *CommandUpdate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandUpdate) HandleSingle() (interface{}, error) {
	return c.IDOrName(networks.IDFromName)
}

func (c *CommandUpdate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := networks.Update(c.ServiceClient, item.(string), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["network"]
	default:
		out <- err
	}
}

func (c *CommandUpdate) PipeFieldOptions() []string {
	return []string{"id"}
}
