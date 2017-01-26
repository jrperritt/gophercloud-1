package network

import (
	"fmt"

	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	NetworkV2Command
	traits.Waitable
	traits.Pipeable
}

var (
	cDelete                          = new(CommandDelete)
	_       interfaces.Waiter        = cDelete
	_       interfaces.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--id <ID> | --name <NAME> | --stdin id]"),
	Description:  "Deletes a network",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the network",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the network.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.IDOrName(networks.IDFromName)
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	err := networks.Delete(c.ServiceClient, item.(string)).ExtractErr()
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted network [%s]", item.(string))
	default:
		out <- err
	}
}
