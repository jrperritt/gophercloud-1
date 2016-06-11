package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/gophercloud/pagination"
)

type commandList struct {
	openstack.Command
	opts volumes.ListOptsBuilder
}

var list = func() cli.Command {
	c := new(commandList)
	c.SetFlags(flagsList)
	c.SetDefaultFields()
	return openstack.NewCommand(c)
}

func (c commandList) Name() string {
	return "list"
}

func (c commandList) Usage() string {
	return util.Usage(commandPrefix, "list", "")
}

func (c commandList) Description() string {
	return "Lists existing volumes"
}

var flagsList = []cli.Flag{
	cli.StringFlag{
		Name:  "name",
		Usage: "Only list volumes with this name.",
	},
	cli.StringFlag{
		Name:  "status",
		Usage: "Only list volumes that have this status.",
	},
}

//var keysList = []string{"ID", "Name", "Bootable", "Size", "Status", "VolumeType", "SnapshotID"}

func (command *commandList) HandleFlags() error {
	c := command.Ctx.CLIContext

	opts := &osVolumes.ListOpts{
		Name:   c.String("name"),
		Status: c.String("status"),
	}

	resource.Params = &paramsList{
		opts: opts,
	}

	return nil
}

func (command *commandList) Execute() {
	opts := resource.Params.(*paramsList).opts
	pager := osVolumes.List(command.Ctx.ServiceClient, opts)
	var volumes []map[string]interface{}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		info, err := osVolumes.ExtractVolumes(page)
		if err != nil {
			return false, err
		}
		result := make([]map[string]interface{}, len(info))
		for j, volume := range info {
			result[j] = volumeSingle(&volume)
		}
		volumes = append(volumes, result...)
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	if len(volumes) == 0 {
		resource.Result = nil
	} else {
		resource.Result = volumes
	}
}
