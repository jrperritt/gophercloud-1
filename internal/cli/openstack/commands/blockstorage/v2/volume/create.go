package volume

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"gopkg.in/urfave/cli.v1"
)

type Create struct {
	command
	traits.Pipeable
	traits.Waitable
	opts volumes.CreateOptsBuilder
}

var (
	tCreate                          = new(Create)
	_       interfaces.Waiter        = tCreate
	_       interfaces.PipeCommander = tCreate

	flagsCreate = openstack.CommandFlags(tCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--size <SIZE> [--name <NAME> | --stdin name]"),
	Description:  "Creates a network",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, tCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *Create) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "size",
			Usage: "[required] The size of the volume, in GB",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] The volume name",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional] The field being piped into STDIN. Valid values are: name",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] The volume description",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant who should own this network.",
		},
	}
}

func (c *Create) HandleFlags() error {
	opts := &volumes.CreateOpts{
		AvailabilityZone: c.Context().String("availability-zone"),
		Description:      c.Context().String("description"),
	}

	c.opts = opts

	return nil
}

func (c *Create) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"size"})
}

func (c *Create) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	opts := *c.opts.(*volumes.CreateOpts)
	opts.Name = item.(string)
	err := volumes.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["volume"]
	default:
		out <- err
	}
}

func (c *Create) PipeFieldOptions() []string {
	return []string{"name"}
}
