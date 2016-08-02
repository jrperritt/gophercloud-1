package instance

import (
	"fmt"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type commandDeleteMetadata struct {
	openstack.CommandUtil
	InstanceV2Command
	wait bool
	opts []string
	*openstack.Progress
}

var (
	cDeleteMetadata                   = new(commandDeleteMetadata)
	_               lib.PipeCommander = cDeleteMetadata
)

var deleteMetadata = cli.Command{
	Name:         "delete-metadata",
	Usage:        util.Usage(commandPrefix, "delete-metadata", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes metadata associated with a server",
	Action:       actionDeleteMetadata,
	Flags:        openstack.CommandFlags(flagsDeleteMetadata, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDeleteMetadata) },
}

func actionDeleteMetadata(ctx *cli.Context) {
	c := new(commandDeleteMetadata)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsDeleteMetadata = []cli.Flag{
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
	cli.BoolFlag{
		Name: "wait",
		Usage: "[optional] If provided, will wait to return until the metadata has been deleted from all servers\n" +
			"\tarriving on STDIN.",
	},
}

func (c *commandDeleteMetadata) HandleFlags() error {
	c.wait = c.Context.IsSet("wait")
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

	var wg sync.WaitGroup

	ch := make(chan interface{})

	for item := range in {
		wg.Add(1)
		item := item
		go func() {
			defer wg.Done()
			id := item.(string)
			for _, key := range c.opts {
				err := servers.DeleteMetadatum(c.ServiceClient, id, key).ExtractErr()
				var v interface{}
				switch err {
				case nil:
					v = fmt.Sprintf("Deleted metadata [%s] from server [%s]", key, id)
				default:
					v = err
				}
				switch c.wait {
				case true:
					ch <- v
				case false:
					out <- v
				}
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	msgs := make([]string, 0)

	for raw := range ch {
		switch msg := raw.(type) {
		case error:
			out <- msg
		case string:
			msgs = append(msgs, msg)
		}
	}

	for _, msg := range msgs {
		out <- msg
	}
}

func (c *commandDeleteMetadata) PipeFieldOptions() []string {
	return []string{"id"}
}
