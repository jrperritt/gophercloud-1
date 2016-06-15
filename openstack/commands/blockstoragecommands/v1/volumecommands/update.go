package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandUpdate struct {
	openstack.Command
	id   string
	opts volumes.UpdateOptsBuilder
}

func (c commandUpdate) Name() string {
	return "delete"
}

func (c commandUpdate) Usage() string {
	return util.Usage(commandPrefix, "update", "[--id <volumeID> | --name <volumeName>]")
}

func (c commandUpdate) Description() string {
	return "Updates a volume"
}

var flagsUpdate = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `name` isn't provided] The ID of the volume.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `id` isn't provided] The name of the volume.",
	},
	cli.StringFlag{
		Name:  "rename",
		Usage: "[optional] A new name for this volume.",
	},
	cli.StringFlag{
		Name:  "description",
		Usage: "[optional] A new description for this volume.",
	},
}

func (c *commandUpdate) HandleFlags() (err error) {
	c.opts = &volumes.UpdateOpts{
		Name:        c.String("rename"),
		Description: c.String("description"),
	}
	c.id, err = c.IDOrName(volumes.IDFromName)
	return
}

func (c *commandUpdate) Execute(_ lib.Resourcer) (r lib.Resulter) {
	var m map[string]interface{}
	err := volumes.Update(c.ServiceClient(), c.id, c.opts).ExtractInto(&m)
	if err != nil {
		r.SetError(err)
		return
	}
	r.SetValue(m)
	return
}

/*
func (c *commandUpdate) PreCSV() error {
	resource.FlattenMap("Attachments")
	return nil
}

func (c *commandUpdate) PreTable() error {
	return command.PreCSV(resource)
}
*/
