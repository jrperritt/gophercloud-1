package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandCreate struct {
	openstack.CommandUtil
	VolumeV1Command
	opts volumes.CreateOptsBuilder
	wait bool
}

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--size <size>"),
	Description:  "Creates a volume",
	Action:       actionCreate,
	Flags:        openstack.CommandFlags(flagsCreate, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsCreate) },
}

func actionCreate(ctx *cli.Context) {
	c := new(commandCreate)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsCreate = []cli.Flag{
	cli.IntFlag{
		Name:  "size",
		Usage: "[required] The size of this volume (in gigabytes). Valid values are between 75 and 1024.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional] A name for this volume.",
	},
	cli.StringFlag{
		Name:  "description",
		Usage: "[optional] A description for this volume.",
	},
	cli.StringFlag{
		Name:  "volume-type",
		Usage: "[optional] The volume type of this volume.",
	},
	cli.BoolFlag{
		Name:  "wait",
		Usage: "[optional] If provided, the command will wait to return until the volume is available.",
	},
}

func (c *commandCreate) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"size"})
	if err != nil {
		return err
	}

	c.wait = c.Context.IsSet("wait")

	c.opts = &volumes.CreateOpts{
		Size:        c.Context.Int("size"),
		Name:        c.Context.String("name"),
		Description: c.Context.String("description"),
		VolumeType:  c.Context.String("volume-type"),
	}

	return nil
}

func (c *commandCreate) Execute(_ interface{}, out chan interface{}) {
	defer func() {
		close(out)
	}()
	var m map[string]map[string]interface{}
	err := volumes.Create(c.ServiceClient, c.opts).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}

	id := m["volume"]["id"].(string)

	if c.wait {
		err = volumes.WaitForStatus(c.ServiceClient, id, "available", 600)
		if err != nil {
			out <- err
			return
		}

		_, err = volumes.Get(c.ServiceClient, id).Extract()
		if err != nil {
			out <- err
			return
		}
	}

	out <- m["volume"]
}

/*
func (c *commandCreate) PreCSV() error {
	resource.FlattenMap("Metadata")
	resource.FlattenMap("Attachments")
	return nil
}

func (c *commandCreate) PreTable() error {
	return c.PreCSV(resource)
}
*/
