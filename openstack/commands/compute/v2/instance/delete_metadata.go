package instance

import (
	"fmt"
	"strings"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type commandDeleteMetadata struct {
	openstack.CommandUtil
	InstanceV2Command
	opts []string
}

var (
	cDeleteMetadata                   = new(commandDeleteMetadata)
	_               lib.PipeCommander = cDeleteMetadata
	_               lib.Waiter        = cDeleteMetadata

	flagsDeleteMetadata = openstack.CommandFlags(cDeleteMetadata)
)

var deleteMetadata = cli.Command{
	Name:         "delete-metadata",
	Usage:        util.Usage(commandPrefix, "delete-metadata", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes metadata associated with a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDeleteMetadata) },
	Flags:        flagsDeleteMetadata,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDeleteMetadata) },
}

func (c *commandDeleteMetadata) Flags() []cli.Flag {
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

func (c *commandDeleteMetadata) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	c.opts = strings.Split(c.Context.String("metadata-keys"), ",")
	return nil
}

func (c *commandDeleteMetadata) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandDeleteMetadata) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandDeleteMetadata) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		id := item.(string)
		for _, key := range c.opts {
			err := servers.DeleteMetadatum(c.ServiceClient, id, key).ExtractErr()
			switch err {
			case nil:
				out <- fmt.Sprintf("Deleted metadata [%s] from server [%s]", key, id)
			default:
				out <- err
			}
		}
	}
}

func (c *commandDeleteMetadata) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandDeleteMetadata) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}
