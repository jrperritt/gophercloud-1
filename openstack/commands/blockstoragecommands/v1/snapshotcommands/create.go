package snapshotcommands

import (
	"reflect"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/snapshots"
)

var (
	cc *commandCreate
	_  lib.Commander = cc
)

type commandCreate struct {
	openstack.Command
	openstack.Progress
	opts snapshots.CreateOptsBuilder
	wait bool
}

func (c *commandCreate) Name() string {
	return "create"
}

func (c *commandCreate) Usage() string {
	return util.Usage(commandPrefix, "create", "--volume-id <volumeID>")
}

func (c *commandCreate) Description() string {
	return "Creates a snapshot of a volume"
}

var flagsCreate = []cli.Flag{
	cli.StringFlag{
		Name:  "volume-id",
		Usage: "[required] The volume ID from which to create this snapshot",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional] A name for this snapshot",
	},
	cli.StringFlag{
		Name:  "description",
		Usage: "[optional] A description for this snapshot",
	},
	cli.BoolFlag{
		Name:  "wait-for-completion",
		Usage: "[optional] If provided, the command will wait to return until the snapshot is available",
	},
}

func (c *commandCreate) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"volume-id"})
	if err != nil {
		return err
	}

	if c.IsSet("wait-for-completion") {
		c.wait = true
	}

	c.opts = &snapshots.CreateOpts{
		VolumeID:    c.String("volume-id"),
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	return nil
}

func (c *commandCreate) Execute(_ lib.Resourcer) (r lib.Resulter) {
	var m map[string]interface{}
	err := snapshots.Create(c.ServiceClient(), c.opts).ExtractInto(m)
	if err != nil {
		r.SetError(err)
		return
	}

	if c.wait {
		err = snapshots.WaitForStatus(c.ServiceClient(), m["id"].(string), "available", 600)
		if err != nil {
			r.SetError(err)
			return
		}
	}

	r.SetValue(m)
	return
}

func (c *commandCreate) ResponseType() reflect.Type {
	return reflect.TypeOf(&snapshots.Snapshot{})
}

/*
func (c *commandCreate) PreJSON(resource lib.Resourcer) error {
	m := make(map[string]interface{})
	err = resource.Result.(gophercloud.Result).ExtractInto(m)
	if err != nil {
		return err
	}
	resource.Result = m
	return nil
}

func (c *commandCreate) PreTable(resource lib.Resourcer) error {
	resource.FlattenMap("Metadata")
}
*/
