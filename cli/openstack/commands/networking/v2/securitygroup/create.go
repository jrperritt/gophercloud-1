package securitygroup

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	SecurityGroupV2Command
	traits.Pipeable
	traits.Waitable
	traits.DataResp
	opts groups.CreateOptsBuilder
}

var (
	cCreate                          = new(CommandCreate)
	_       interfaces.Waiter        = cCreate
	_       interfaces.PipeCommander = cCreate

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "[--name <NAME> | --stdin name]"),
	Description:  "Creates a security group",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name for the security group.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A description for the security group.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant who should own this network.",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	c.opts = &groups.CreateOpts{
		Description: c.Context().String("description"),
		TenantID:    c.Context().String("tenant-id"),
	}
	return nil
}

func (c *CommandCreate) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandCreate) Execute(_ interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := groups.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["security_group"]
	default:
		out <- err
	}
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}
