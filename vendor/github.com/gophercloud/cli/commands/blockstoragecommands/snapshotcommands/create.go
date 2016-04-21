package snapshotcommands

import (
	"reflect"

	gophercloudCLI "github.com/gophercloud/cli"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

var Create = cli.Command{
	Name:        "create",
	Usage:       util.Usage(commandPrefix, "create", "--volume-id <volumeID>"),
	Description: "Creates a snapshot of a volume",
	Action:      ActionCreate,
	Flags:       openstackCLI.CommandFlags(FlagsCreate, KeysCreate),
	BashComplete: func(c *cli.Context) {
		keys, err := KeysCreateFunc()
		if err != nil {
			keys = make([]string, 0)
		}
		openstackCLI.CompleteFlags(openstackCLI.CommandFlags(FlagsCreate, keys))
	},
}

func ActionCreate(cliContext *cli.Context) {
	gophercloudCLI.Run(cliContext, &CommandCreate{})
}

func FlagsCreate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "[required] The volume ID from which to create this snapshot.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] A name for this snapshot.",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "[optional] A description for this snapshot.",
		},
		cli.BoolFlag{
			Name:  "wait-for-completion",
			Usage: "[optional] If provided, the command will wait to return until the snapshot is available.",
		},
	}
}

var KeysCreateFunc = func() ([]string, error) {
	return util.BuildKeys(snapshots.Snapshot)
}

type CommandCreate struct {
	gophercloudCLI.Command
	Opts *snapshots.CreateOptsBuilder
	Wait bool
}

func (command *CommandCreate) Keys() ([]string, error) {
	return KeysCreateFunc()
}

func (command *CommandCreate) HandleFlags() error {
	err := command.CheckFlagsSet([]string{"volume-id"})
	if err != nil {
		return err
	}

	c := command.CLIContext
	if c.IsSet("wait-for-completion") {
		command.Wait = true
	}

	command.Opts = &snapshots.CreateOpts{
		VolumeID:    c.String("volume-id"),
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	return nil
}

func (command *CommandCreate) Execute(resource *openstackCLI.Resource) {
	var m map[string]interface{}
	res := snapshots.Create(command.ServiceClient, command.Opts)

	resExtract := res
	snapshot, err := resExtract.Extract()
	if err != nil {
		resource.Err = err
		return
	}

	if command.Wait {
		err = snapshots.WaitForStatus(command.ServiceClient, snapshot.ID, "available", 600)
		if err != nil {
			resource.Err = err
			return
		}
	}

	resource.Result = res
}

func (command *CommandCreate) ResponseType(resource *openstackCLI.Resource) reflect.Type {
	return reflect.TypeOf(*snapshots.Snapshot)
}

func (command *CommandCreate) PreJSON(resource *openstackCLI.Resource) error {
	m := make(map[string]interface{})
	err = resource.Result.(gophercloud.Result).ExtractInto(m)
	if err != nil {
		return err
	}
	resource.Result = m
	return nil
}

func (command *CommandCreate) PreTable(resource *openstackCLI.Resource) error {
	resource.FlattenMap("Metadata")
}
