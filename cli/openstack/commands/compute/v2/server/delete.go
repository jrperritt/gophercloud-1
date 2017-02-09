package server

import (
	"fmt"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	ServerV2Command
	traits.Waitable
	traits.TextProgressable
	traits.MsgResp
}

var (
	cDelete                          = new(CommandDelete)
	_       interfaces.PipeCommander = cDelete
	_       interfaces.Progresser    = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(CommandPrefix, "delete", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
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

func (c *CommandDelete) HandleFlags() error {
	return nil
}

func (c *CommandDelete) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	err := servers.Delete(c.ServiceClient, id).ExtractErr()
	if err != nil {
		out <- err
		return
	}
	switch c.ShouldWait() || c.ShouldProgress() {
	case true:
		out <- id
	default:
		out <- fmt.Sprintf("Deleting server [%s]", id)
	}
}

func (c *CommandDelete) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *CommandDelete) WaitFor(raw interface{}, out chan<- interface{}) {
	id := raw.(string)

	err := util.WaitFor(900, func() (bool, error) {
		_, err := servers.Get(c.ServiceClient, id).Extract()
		if err != nil {
			out <- fmt.Sprintf("Deleted server [%s]", id)
			return true, nil
		}
		c.ProgUpdateChIn() <- c.RunningMsg
		return false, nil
	})

	if err != nil {
		out <- err
	}
}

func (c *CommandDelete) InitProgress(donech chan interface{}) {
	c.RunningMsg = "Deleting"
	c.DoneMsg = "Deleted"
	c.TextProgressable.InitProgress(donech)
}
