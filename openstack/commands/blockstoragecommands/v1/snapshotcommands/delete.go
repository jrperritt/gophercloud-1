package snapshotcommands

import (
	"fmt"
	"reflect"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

var (
	cd *commandDelete
	_  lib.PipeCommander = cd
)

type commandDelete struct {
	openstack.PipeCommand
	id   string
	wait bool
}

var remove = func() cli.Command {
	c := new(commandDelete)
	c.SetFlags(flagsDelete)
	c.SetDefaultFields()
	return openstack.NewCommand(c)
}

func (c commandDelete) Name() string {
	return "delete"
}

func (c commandDelete) Usage() string {
	return util.Usage(commandPrefix, c.Name(), "[--id <snapshotID> | --name <snapshotName> | --stdin id]")
}

func (c commandDelete) Description() string {
	return "Deletes a snapshot"
}

var flagsDelete = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "The ID of the snapshot.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "The name of the snapshot.",
	},
	cli.StringFlag{
		Name:  "stdin",
		Usage: "The field being piped into STDIN. Valid values are: id, name",
	},
	cli.BoolFlag{
		Name:  "wait-for-completion",
		Usage: "If provided, the command will wait to return until the snapshot is available.",
	},
}

func (c *commandDelete) HandleFlags() (_ error) {
	if c.IsSet("wait-for-completion") {
		c.wait = true
	}
	return
}

func (c commandDelete) HandlePipe(resource lib.Resourcer) (err error) {
	switch c.stdinField {
	case "id":
		c.id = v.(string)
	case "name":
		c.id, err = snapshots.IDFromName(c.ServiceClient(), v.(string))
	}
	return
}

func (c commandDelete) HandleSingle(resource lib.Resourcer) (err error) {
	c.id, err = c.IDOrName(snapshots.IDFromName)
	return
}

func (c commandDelete) Execute(_ lib.Resourcer) lib.Resulter {
	result := resource.NewResult()
	err := snapshots.Delete(c.ServiceClient(), resource.StdInParams().(*stdinDelete).SnapshotID).ExtractErr()
	if err != nil {
		result.SetError(err)
		return result
	}

	if c.wait {
		i := 0
		for i < 120 {
			_, err := snapshots.Get(c.ServiceClient(), snapshotID).Extract()
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

func (c commandDelete) PipeField() string {
	return c.stdinField
}

func (c commandDelete) PipeFieldOptions() []string {
	return []string{"id", "name"}
}

func (c commandDelete) ResponseType() reflect.Type {
	return reflect.TypeOf("")
}
