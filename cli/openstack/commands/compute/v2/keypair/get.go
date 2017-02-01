package keypair

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	KeypairV2Command
	traits.Waitable
	traits.Pipeable
	traits.DataResp
}

var (
	cGet                          = new(CommandGet)
	_    interfaces.Waiter        = cGet
	_    interfaces.PipeCommander = cGet

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--name <NAME> | --stdin name]"),
	Description:  "Gets a keypair",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the keypair.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

func (c *CommandGet) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandGet) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	err := keypairs.Get(c.ServiceClient, item.(string)).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["keypair"]
	default:
		out <- err
	}
}

func (c *CommandGet) PipeFieldOptions() []string {
	return []string{"name"}
}
