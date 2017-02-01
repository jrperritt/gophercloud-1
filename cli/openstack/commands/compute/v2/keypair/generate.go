package keypair

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"gopkg.in/urfave/cli.v1"
)

type CommandGenerate struct {
	KeypairV2Command
	traits.Waitable
	traits.Pipeable
	traits.DataResp
}

var (
	cGenerate                          = new(CommandGenerate)
	_         interfaces.Waiter        = cGenerate
	_         interfaces.PipeCommander = cGenerate

	flagsGenerate = openstack.CommandFlags(cGenerate)
)

var generate = cli.Command{
	Name:         "generate",
	Usage:        util.Usage(commandPrefix, "generate", "[--name <NAME> | --stdin name]"),
	Description:  "Generates a keypair",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGenerate) },
	Flags:        flagsGenerate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGenerate) },
}

func (c *CommandGenerate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name for the keypair.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

func (c *CommandGenerate) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandGenerate) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	err := keypairs.Create(c.ServiceClient, keypairs.CreateOpts{Name: item.(string)}).ExtractInto(&m)
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}

func (c *CommandGenerate) PipeFieldOptions() []string {
	return []string{"name"}
}
