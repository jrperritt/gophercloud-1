package securitygroup

import (
	"fmt"

	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	SecurityGroupV2Command
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
	Usage:        util.Usage(commandPrefix, "delete", "[--id <ID> | --name <NAME> | --stdin name]"),
	Description:  "Deletes a security group",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
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

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.IDOrName(groups.IDFromName)
}

func (c *CommandDelete) Execute(raw interface{}, out chan interface{}) {
	err := groups.Delete(c.ServiceClient(), raw.(string)).ExtractErr()
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted security group [%s]", raw.(string))
	default:
		out <- err
	}
}
