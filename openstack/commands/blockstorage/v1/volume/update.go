package volume

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandUpdate struct {
	openstack.CommandUtil
	VolumeV1Command
	id   string
	opts volumes.UpdateOptsBuilder
}

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", "[--id <volumeID> | --name <volumeName>]"),
	Description:  "Updates a volume",
	Action:       actionUpdate,
	Flags:        openstack.CommandFlags(flagsUpdate, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsUpdate) },
}

func actionUpdate(ctx *cli.Context) {
	c := new(commandUpdate)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsUpdate = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `name` isn't provided] The ID of the volume to update",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `id` isn't provided] The name of the volume to update",
	},
	cli.StringFlag{
		Name:  "rename",
		Usage: "[optional] A new name for this volume",
	},
	cli.StringFlag{
		Name:  "description",
		Usage: "[optional] A new description for this volume",
	},
}

func (c *commandUpdate) HandleFlags() (err error) {
	c.opts = &volumes.UpdateOpts{
		Name:        c.Context.String("rename"),
		Description: c.Context.String("description"),
	}
	c.id, err = c.IDOrName(volumes.IDFromName)
	return
}

func (c *commandUpdate) Execute(_ interface{}, out chan interface{}) {
	defer func() {
		close(out)
	}()
	var m map[string]interface{}
	err := volumes.Update(c.ServiceClient, c.id, c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["volume"]
	default:
		out <- err
	}
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
