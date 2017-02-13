package keypair

import (
	"fmt"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	KeypairV2Command
	traits.Waitable
	traits.Pipeable
}

var (
	cDelete                          = new(CommandDelete)
	_       interfaces.Waiter        = cDelete
	_       interfaces.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--name <NAME> | --stdin name]"),
	Description:  "Deletes a keypair",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
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

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	err := keypairs.Delete(c.ServiceClient(), item.(string)).ExtractErr()
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully deleted keypair [%s]", item.(string))
	default:
		out <- err
	}
}

func (c *CommandDelete) PipeFieldOptions() []string {
	return []string{"name"}
}
