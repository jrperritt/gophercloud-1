package flavor

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	FlavorV2Command
	opts flavors.CreateOptsBuilder
}

var (
	cCreate = new(CommandCreate)

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--name NAME --disk DISK --ram RAM --vcpus VCPUS"),
	Description:  "Creates a flavor",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the flavor",
		},
		cli.IntFlag{
			Name:  "disk",
			Usage: "[required] The disk space of the flavor (in GBs)",
		},
		cli.IntFlag{
			Name:  "ram",
			Usage: "[required] The memory of the flavor (in MBs)",
		},
		cli.IntFlag{
			Name:  "vcpus",
			Usage: "[required] The number of virtual CPUs of the flavor",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional] The ID of the flavor.",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"name", "disk", "ram", "vcpus"})
	if err != nil {
		return err
	}

	disk := c.Context().Int("disk")

	opts := new(flavors.CreateOpts)
	opts.Name = c.Context().String("name")
	opts.Disk = &disk
	opts.RAM = c.Context().Int("ram")
	opts.VCPUs = c.Context().Int("vcpus")
	opts.ID = c.Context().String("id")

	c.opts = opts
	return nil
}

func (c *CommandCreate) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	err := flavors.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["flavor"]
	default:
		out <- err
	}
}
