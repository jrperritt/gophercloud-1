package server

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	ServerV2Command
	traits.Pipeable
	traits.Waitable
	traits.DataResp
}

var (
	cGet                          = new(CommandGet)
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(CommandPrefix, "get", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Gets a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the server.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandGet) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	var m map[string]map[string]interface{}
	err := servers.Get(c.ServiceClient, id).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["server"]
	default:
		out <- err
	}
}