package port

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"gopkg.in/urfave/cli.v1"
)

type CommandUpdate struct {
	PortV2Command
	traits.Waitable
	traits.DataResp
	opts ports.UpdateOptsBuilder
}

var (
	cUpdate                          = new(CommandUpdate)
	_       interfaces.PipeCommander = cUpdate

	flagsUpdate = openstack.CommandFlags(cUpdate)
)

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(CommandPrefix, "update", "[--id <ID> | --name <NAME> | --stdin id]"),
	Description:  "Updates a port",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpdate) },
	Flags:        flagsUpdate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpdate) },
}

func (c *CommandUpdate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the port to update.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the port to update.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id.",
		},
		cli.StringFlag{
			Name:  "new-name",
			Usage: "[optional] The name that the port should have.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "[optional] If provided, the port will be up. Options are: true, false",
		},
		cli.StringFlag{
			Name:  "security-groups",
			Usage: "[optional] A comma-separated list of security group IDs for this port.",
		},
		cli.StringFlag{
			Name:  "device-id",
			Usage: "[optional] A device ID to associate with the port.",
		},
	}
}

func (c *CommandUpdate) HandleFlags() error {
	opts := &ports.UpdateOpts{
		Name:     c.Context.String("new-name"),
		DeviceID: c.Context.String("device-id"),
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

	if c.Context.IsSet("security-groups") {
		opts.SecurityGroups = strings.Split(c.Context.String("security-groups"), ",")
	}

	c.opts = opts
	c.Wait = c.Context.IsSet("wait")

	return nil
}

func (c *CommandUpdate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandUpdate) HandleSingle() (interface{}, error) {
	return c.IDOrName(ports.IDFromName)
}

func (c *CommandUpdate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := ports.Update(c.ServiceClient, item.(string), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["port"]
	default:
		out <- err
	}
}

func (c *CommandUpdate) PipeFieldOptions() []string {
	return []string{"id"}
}
