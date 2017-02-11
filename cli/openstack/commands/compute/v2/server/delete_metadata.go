package server

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type CommandDeleteMetadata struct {
	ServerV2Command
	traits.Waitable
	opts []string
}

var (
	cDeleteMetadata                          = new(CommandDeleteMetadata)
	_               interfaces.PipeCommander = cDeleteMetadata

	flagsDeleteMetadata = openstack.CommandFlags(cDeleteMetadata)
)

var deleteMetadata = cli.Command{
	Name:         "delete-metadata",
	Usage:        util.Usage(CommandPrefix, "delete-metadata", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes metadata associated with a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDeleteMetadata) },
	Flags:        flagsDeleteMetadata,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDeleteMetadata) },
}

func (c *CommandDeleteMetadata) Flags() []cli.Flag {
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
		cli.StringFlag{
			Name:  "metadata-keys",
			Usage: "[required] A comma-separated string of keys of the metadata to delete from the server.",
		},
	}
}

func (c *CommandDeleteMetadata) HandleFlags() error {
	c.opts = strings.Split(c.Context().String("metadata-keys"), ",")
	return nil
}

func (c *CommandDeleteMetadata) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandDeleteMetadata) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandDeleteMetadata) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	for _, key := range c.opts {
		err := servers.DeleteMetadatum(c.ServiceClient(), id, key).ExtractErr()
		switch err {
		case nil:
			out <- fmt.Sprintf("Deleted metadata [%s] from server [%s]", key, id)
		default:
			out <- err
		}
	}
}

func (c *CommandDeleteMetadata) PipeFieldOptions() []string {
	return []string{"id"}
}
