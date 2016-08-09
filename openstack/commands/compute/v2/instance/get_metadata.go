package instance

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type commandGetMetadata struct {
	openstack.CommandUtil
	InstanceV2Command
	opts []string
}

var (
	cGetMetadata                   = new(commandGetMetadata)
	_            lib.PipeCommander = cGetMetadata
	_            lib.Waiter        = cGetMetadata
)

var getMetadata = cli.Command{
	Name:         "get-metadata",
	Usage:        util.Usage(commandPrefix, "get-metadata", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Gets metadata associated with a server",
	Action:       actionGetMetadata,
	Flags:        openstack.CommandFlags(cGetMetadata),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsGetMetadata) },
}

func actionGetMetadata(ctx *cli.Context) {
	c := new(commandGetMetadata)
	c.Context = ctx
	lib.Run(ctx, c)
}

func (c *commandGetMetadata) Flags() []cli.Flag {
	return flagsGetMetadata
}

func (c *commandGetMetadata) Fields() []string {
	return []string{""}
}

var flagsGetMetadata = []cli.Flag{
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

func (c *commandGetMetadata) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	return nil
}

func (c *commandGetMetadata) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandGetMetadata) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandGetMetadata) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		id := item.(string)
		m, err := servers.Metadata(c.ServiceClient, id).Extract()
		switch err {
		case nil:
			out <- m
		default:
			out <- err
		}
	}
}

func (c *commandGetMetadata) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandGetMetadata) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}
