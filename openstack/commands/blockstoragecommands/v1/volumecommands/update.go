package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandUpdate struct {
	openstack.Command
	id   string
	opts volumes.UpdateOptsBuilder
}

var update = func() cli.Command {
	c := new(commandUpdate)
	c.SetFlags(flagsUpdate)
	c.SetDefaultFields()
	return openstack.NewCommand(c)
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

func (command *commandUpdate) HandleFlags() error {
	volumeID, err := command.Ctx.IDOrName(osVolumes.IDFromName)
	if err != nil {
		return err
	}

	c := command.Ctx.CLIContext

	opts := &osVolumes.UpdateOpts{
		Name:        c.String("rename"),
		Description: c.String("description"),
	}

	resource.Params = &paramsUpdate{
		volumeID: volumeID,
		opts:     opts,
	}

	return nil
}

func (command *commandUpdate) Execute() {
	opts := resource.Params.(*paramsUpdate).opts
	volumeID := resource.Params.(*paramsUpdate).volumeID
	volume, err := osVolumes.Update(command.Ctx.ServiceClient, volumeID, opts).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = volumeSingle(volume)
}

func (command *commandUpdate) PreCSV() error {
	resource.FlattenMap("Attachments")
	return nil
}

func (command *commandUpdate) PreTable() error {
	return command.PreCSV(resource)
}
