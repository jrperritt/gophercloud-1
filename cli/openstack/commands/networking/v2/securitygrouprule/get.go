package securitygrouprule

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	SecurityGroupRuleV2Command
	traits.Pipeable
	traits.Waitable
	traits.Fieldsable
}

var (
	cGet                          = new(CommandGet)
	_    interfaces.Waiter        = cGet
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--id <ID> | --stdin id]"),
	Description:  "Gets a security group rule",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` isn't provided] The ID of the security group rule.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.Context().String("id"), c.CheckFlagsSet([]string{"id"})
}

func (c *CommandGet) Execute(raw interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := rules.Get(c.ServiceClient(), raw.(string)).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["security_group_rule"]
	default:
		out <- err
	}
}
