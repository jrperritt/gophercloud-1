package snapshotcommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/gophercloud/gophercloud/pagination"
)

type CommandList struct {
	openstack.Command
	Opts *snapshots.ListOpts
}

func (c CommandList) CommandFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "Only list snapshots with this volume ID.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list snapshots with this name.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list snapshots that have this status.",
		},
	}
}

func (c CommandList) Name() string {
	return "list"
}

// List is the command to list snapshots
var List = cli.Command{
	Name:        "list",
	Usage:       util.Usage(commandPrefix, "list", ""),
	Description: "Lists existing snapshots",
	Action:      ActionList,
	Flags:       gophercloudCLI.CommandFlags(flagsList, keysList),
	BashComplete: func(c *cli.Context) {
		gophercloudCLI.CompleteFlags(gophercloudCLI.CommandFlags(flagsList, keysList))
	},
}

var keysList = []string{"ID", "Name", "Size", "Status", "VolumeID", "VolumeType", "SnapshotID", "Bootable"}

func ActionList(c *cli.Context) {
	gophercloudCLI.Run(cliContext, &CommandList{})
}

func (command *CommandList) Keys() []string {
	return keysList
}

func (command *CommandList) ServiceClientType() string {
	return serviceClientType
}

func (command *CommandList) HandleFlags(resource *Resource) error {
	c := command.CLIContext

	command.Opts = &snapshots.ListOpts{
		VolumeID: c.String("volume-id"),
		Name:     c.String("name"),
		Status:   c.String("status"),
	}

	return nil
}

func (command *CommandList) Execute(resource *Resource) {
	opts := command.Opts
	pager := snapshots.List(command.ServiceClient, opts)
	var snapshots []map[string]interface{}
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		info, err := snapshots.ExtractSnapshots(page)
		if err != nil {
			return false, err
		}
		result := make([]map[string]interface{}, len(info))
		for j, snapshot := range info {
			result[j] = snapshotSingle(&snapshot)
		}
		snapshots = append(snapshots, result...)
		return true, nil
	})
	if err != nil {
		resource.Err = err
		return
	}
	if len(snapshots) == 0 {
		resource.Result = nil
	} else {
		resource.Result = snapshots
	}
}
