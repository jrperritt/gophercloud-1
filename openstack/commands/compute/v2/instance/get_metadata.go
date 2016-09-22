package instance

import (
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type CommandGetMetadata struct {
	ServerV2Command
	commands.WaitCommand
	opts []string
}

var (
	cGetMetadata                         = new(CommandGetMetadata)
	_            openstack.PipeCommander = cGetMetadata

	flagsGetMetadata = openstack.CommandFlags(cGetMetadata)
)

var getMetadata = cli.Command{
	Name:         "get-metadata",
	Usage:        util.Usage(commandPrefix, "get-metadata", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Gets metadata associated with a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGetMetadata) },
	Flags:        flagsGetMetadata,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGetMetadata) },
}

func (c *CommandGetMetadata) Flags() []cli.Flag {
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

func (c *CommandGetMetadata) Fields() []string {
	return []string{""}
}

func (c *CommandGetMetadata) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	return nil
}

func (c *CommandGetMetadata) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandGetMetadata) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandGetMetadata) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	m, err := servers.Metadata(c.ServiceClient, id).Extract()
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}

func (c *CommandGetMetadata) PipeFieldOptions() []string {
	return []string{"id"}
}
