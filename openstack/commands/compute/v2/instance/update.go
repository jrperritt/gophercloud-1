package instance

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type commandUpdate struct {
	openstack.CommandUtil
	InstanceV2Command
	id   string
	opts servers.UpdateOptsBuilder
}

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", "[--id <serverID> | --name <serverName>]"),
	Description:  "Updates a server",
	Action:       actionUpdate,
	Flags:        openstack.CommandFlags(flagsUpdate, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsUpdate) },
}

func actionUpdate(ctx *cli.Context) {
	c := new(commandUpdate)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsUpdate = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `name` isn't provided] The ID of the volume to server",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `id` isn't provided] The name of the volume to server",
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

func (c *commandUpdate) HandleFlags() (err error) {
	c.opts = &servers.UpdateOpts{
		Name:       c.Context.String("rename"),
		AccessIPv4: c.Context.String("ipv4"),
		AccessIPv6: c.Context.String("ipv6"),
	}
	c.id, err = c.IDOrName(servers.IDFromName)
	return
}

func (c *commandUpdate) Execute(_ interface{}, out chan interface{}) {
	defer func() {
		close(out)
	}()
	var m map[string]interface{}
	err := servers.Update(c.ServiceClient, c.id, c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["server"]
	default:
		out <- err
	}
}
