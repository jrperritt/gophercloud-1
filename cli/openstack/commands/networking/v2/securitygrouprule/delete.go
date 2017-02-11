package securitygrouprule

import (
	"fmt"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	SecurityGroupRuleV2Command
	traits.Pipeable
	traits.Waitable
}

var (
	cDelete                          = new(CommandDelete)
	_       interfaces.Waiter        = cDelete
	_       interfaces.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--id <ID> | --stdin id]"),
	Description:  "Deletes a security group rule",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
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

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.Context().String("id"), c.CheckFlagsSet([]string{"id"})
}

func (c *CommandDelete) Execute(raw interface{}, out chan interface{}) {
	err := rules.Delete(c.ServiceClient(), raw.(string)).ExtractErr()
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted security group rule [%s]", raw.(string))
	default:
		out <- err
	}
}
