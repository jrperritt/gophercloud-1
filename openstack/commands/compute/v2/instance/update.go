package instance

import (
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type commandUpdate struct {
	ServerV2Command
	id   string
	opts servers.UpdateOptsBuilder
}

var (
	cUpdate                     = new(commandUpdate)
	_       openstack.Commander = cUpdate

	flagsUpdate = openstack.CommandFlags(cUpdate)
)

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", "[--id <serverID> | --name <serverName>]"),
	Description:  "Updates a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpdate) },
	Flags:        flagsUpdate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpdate) },
}

func (c *commandUpdate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The ID of the server",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The name of the server",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] Rename this server",
		},
		cli.StringFlag{
			Name:  "ipv4",
			Usage: "[optional] Change the server's IPv4 address",
		},
		cli.StringFlag{
			Name:  "ipv6",
			Usage: "[optional] Change the server's IPv6 address",
		},
	}
}

func (c *commandUpdate) HandleFlags() (err error) {
	c.opts = &servers.UpdateOpts{
		Name:       c.Context.String("rename"),
		AccessIPv4: c.Context.String("ipv4"),
		AccessIPv6: c.Context.String("ipv6"),
	}
	c.id, err = c.IDOrName(servers.IDFromName)
	return
}

func (c *commandUpdate) Execute(id interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := servers.Update(c.ServiceClient, id.(string), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["server"]
	default:
		out <- err
	}
}
