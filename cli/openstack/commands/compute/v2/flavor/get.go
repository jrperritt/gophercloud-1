package flavor

import (
	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	FlavorCommand
	traits.Waitable
	traits.Pipeable
	traits.DataResp
}

var (
	cGet                          = new(CommandGet)
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--id <flavorID> | --name <flavorName> | --stdin id]"),
	Description:  "Retreives information about a flavor",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the flavor.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the flavor.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.IDOrName(flavors.IDFromName)
}

func (c *CommandGet) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	err := flavors.Get(c.ServiceClient, item.(string)).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["flavor"]
	default:
		out <- err
	}
}
