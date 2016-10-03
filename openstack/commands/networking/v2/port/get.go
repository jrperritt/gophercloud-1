package port

import (
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	PortV2Command
	traits.Waitable
	traits.DataResp
}

var (
	cGet                         = new(CommandGet)
	_    openstack.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(CommandPrefix, "get", "[--id <ID> | --name <NAME> | --stdin id]"),
	Description:  "Gets a port",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
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

func (c *CommandGet) Fields() []string {
	return []string{""}
}

func (c *CommandGet) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	return nil
}

func (c *CommandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.IDOrName(ports.IDFromName)
}

func (c *CommandGet) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := ports.Get(c.ServiceClient, item.(string)).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["port"]
	default:
		out <- err
	}
}

func (c *CommandGet) PipeFieldOptions() []string {
	return []string{"id"}
}
