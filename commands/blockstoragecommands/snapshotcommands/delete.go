package snapshotcommands

import (
	"fmt"
	"time"

	gophercloudCLI "github.com/gophercloud/cli"
	gophercloudLib "github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

var remove = cli.Command{
	Name:        "delete",
	Usage:       util.Usage(commandPrefix, "delete", "[--id <snapshotID> | --name <snapshotName> | --stdin id]"),
	Description: "Deletes a snapshot",
	Action:      ActionDelete,
	Flags:       gophercloudCLI.CommandFlags(flagsDelete, keysDelete),
	BashComplete: func(c *cli.Context) {
		gophercloudCLI.CompleteFlags(gophercloudCLI.CommandFlags(flagsDelete, keysDelete))
	},
}

func flagsDelete() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the snapshot.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the snapshot.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the snapshot is available.",
		},
	}
}

var keysDelete = []string{}

type commandDelete struct {
	gophercloudCLI.Command
	wait       bool
	snapshotID string
	stdinField string
}

type ParamsDelete struct {
	SnapshotID   string
	SnapshotName string
}

func ActionDelete(cliContext *cli.Context) {
	gophercloudLib.Run(cliContext, &commandDelete{})
}

func (command commandDelete) Keys() []string {
	return keysDelete
}

func (command commandDelete) ServiceClientType() string {
	return serviceClientType
}

func (command *commandDelete) HandleFlags() error {
	if command.CLIContext.IsSet("wait-for-completion") {
		command.wait = true
	}
	return nil
}

func (command commandDelete) HandlePipe(resource gophercloduLib.Resourcer, v interface{}) error {
	var err error
	switch command.stdinField {
	case "id":
		resource.StdIn.(ParamsDelete).SnapshotID = v.(string)
	case "name":
		resource.StdIn.(ParamsDelete).SnapshotID, err = snapshots.IDFromName(command.ServiceClient, v.(string))
	}
	return err
}

func (command commandDelete) Execute(resource lib.Resourcer) lib.Resulter {
	result := resource.NewResult()
	err := snapshots.Delete(command.ServiceClient, resource.GetStdInParams().(ParamsDelete).SnapshotID).ExtractErr()
	if err != nil {
		result.SetError(err)
		return result
	}

	if command.wait {
		i := 0
		for i < 120 {
			_, err := snapshots.Get(command.ServiceClient, snapshotID).Extract()
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		result.SetValue(fmt.Sprintf("Deleted snapshot [%s]\n", snapshotID))
		return result
	}
	result.SetValue(fmt.Sprintf("Deleting snapshot [%s]\n", snapshotID))
	return result
}

func (command commandDelete) StdinFields() []string {
	return []string{"id", "name"}
}
