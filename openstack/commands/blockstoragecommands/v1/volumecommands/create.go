package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandCreate struct {
	openstack.Command
	opts volumes.CreateOptsBuilder
	wait bool
}

var create = func() cli.Command {
	c := new(commandCreate)
	c.SetFlags(flagsCreate)
	c.SetDefaultFields()
	return openstack.NewCommand(c)
}

func (c commandCreate) Name() string {
	return "create"
}

func (c commandCreate) Usage() string {
	return util.Usage(commandPrefix, "create", "--size <size>")
}

func (c commandCreate) Description() string {
	return "Creates a volume"
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
		Name:  "wait-for-completion",
		Usage: "[optional] If provided, the command will wait to return until the volume is available.",
	},
}

func (c *commandCreate) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"size"})
	if err != nil {
		return err
	}

	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	c.opts = volumes.CreateOptsBuilder(
		&volumes.CreateOpts{
			Size:        c.Int("size"),
			Name:        c.String("name"),
			Description: c.String("description"),
			VolumeType:  c.String("volume-type"),
		})

	return nil
}

func (c *commandCreate) Execute(r lib.Resourcer) (res lib.Resulter) {
	volume, err := volumes.Create(c.ServiceClient(), c.opts).Extract()
	if err != nil {
		res.SetError(err)
		return
	}

	if c.wait {
		err = volumes.WaitForStatus(c.ServiceClient(), volume.ID, "available", 600)
		if err != nil {
			res.SetError(err)
			return
		}

		volume, err = volumes.Get(c.ServiceClient(), volume.ID).Extract()
		if err != nil {
			res.SetError(err)
			return
		}
	}

	res.SetValue(volumeSingle(volume))
	return
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
