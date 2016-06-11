package snapshotcommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/jrperritt/gophercloud/openstack/blockstorage/v1/snapshots"
)

type commandGet struct {
	openstack.Command
	id string
}

var get = func() cli.Command {
	c := new(commandGet)
	c.SetFlags(flagsGet)
	c.SetDefaultFields()
	return openstack.NewCommand(c)
}

func (c commandGet) Name() string {
	return "get"
}

func (c commandGet) Usage() string {
	return util.Usage(commandPrefix, "get", "[--id <snapshotID> | --name <snapshotName> | --stdin id]")
}

func (c commandGet) Description() string {
	return "Gets a snapshot"
}

var flagsGet = []cli.Flag{
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
}

func (command *commandGet) HandleFlags() error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(item string) error {
	resource.Params.(*paramsGet).snapshotID = item
	return nil
}

func (command *commandGet) HandleSingle() error {
	snapshotID, err := command.Ctx.IDOrName(snapshots.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).snapshotID = snapshotID
	return nil
}

func (c *commandGet) Execute() {
	snapshotID := resource.Params.(*paramsGet).snapshotID
	snapshot, err := snapshots.Get(c.ServiceClient(), snapshotID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = snapshotSingle(snapshot)
}

func (c *commandGet) StdinField() string {
	return "id"
}

/*
func (command *commandGet) PreCSV() error {
	resource.FlattenMap("Metadata")
	return nil
}

func (command *commandGet) PreTable() error {
	return command.PreCSV(resource)
}
*/
