package securitygroup

import (
	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	SecurityGroupV2Command
	traits.Pipeable
	traits.Waitable
	traits.DataResp
}

var (
	cGet                          = new(CommandGet)
	_    interfaces.Waiter        = cGet
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--id <ID> | --name <NAME> | --stdin name]"),
	Description:  "Gets a security group",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the security group.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the security group.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.IDOrName(groups.IDFromName)
}

func (c *CommandGet) Execute(raw interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := groups.Get(c.ServiceClient, raw.(string)).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["security_group"]
	default:
		out <- err
	}
}
