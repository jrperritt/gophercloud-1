package snapshotcommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

var (
	cg *commandGet
	_  lib.PipeCommander = cg
)

type commandGet struct {
	openstack.PipeCommand
	id string
}

func (c *commandGet) Name() string {
	return "get"
}

func (c *commandGet) Usage() string {
	return util.Usage(commandPrefix, "get", "[--id <snapshotID> | --name <snapshotName> | --stdin id]")
}

func (c *commandGet) Description() string {
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

func (c *commandGet) HandleSingle(r lib.Resourcer) error {
	id, err := c.IDOrName(snapshots.IDFromName)
	if err != nil {
		return err
	}
	c.id = id
	return nil
}

func (c *commandGet) Execute(r lib.Resourcer) (res lib.Resulter) {
	snapshot, err := snapshots.Get(c.ServiceClient(), c.id).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = snapshotSingle(snapshot)
}

func (c *commandGet) PipeField() string {
	return c.stdinField
}

func (c *commandGet) PipeFieldOptions() []string {
	return []string{"id", "name"}
}

func (c *commandGet) SetPipeField(s string) {

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
