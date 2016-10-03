package port

import (
	"fmt"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	PortV2Command
	commands.Waitable
}

var (
	cDelete                         = new(CommandDelete)
	_       openstack.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(CommandPrefix, "delete", "[--id <ID> | --name <NAME> | --stdin id]"),
	Description:  "Deletes a port",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` or `stdin` isn't provided] The ID of the port",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` or `stdin` isn't provided] The name of the port.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` or `id` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandDelete) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	return nil
}

func (c *CommandDelete) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.IDOrName(ports.IDFromName)
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	err := ports.Delete(c.ServiceClient, item.(string)).ExtractErr()
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted port [%s]", item.(string))
	default:
		out <- err
	}
}

func (c *CommandDelete) PipeFieldOptions() []string {
	return []string{"id"}
}
