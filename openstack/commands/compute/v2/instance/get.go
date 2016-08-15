package instance

import (
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type commandGet struct {
	openstack.CommandUtil
	InstanceV2Command
}

var (
	cGet                   = new(commandGet)
	_    lib.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Gets a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsGet) },
}

func (c *commandGet) Flags() []cli.Flag {
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

func (c *commandGet) Fields() []string {
	return []string{""}
}

func (c *commandGet) HandleFlags() error {
	return nil
}

func (c *commandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandGet) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandGet) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
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
}

func (c *commandGet) PipeFieldOptions() []string {
	return []string{"id"}
}
