package volumecommands

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandDelete struct {
	openstack.Command
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
	return util.Usage(commandPrefix, "delete", "[--id <volumeID> | --name <volumeName> | --stdin id]")
}

func (c commandDelete) Description() string {
	return "Deletes a volume"
}

var flagsDelete = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the volume.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the volume.",
	},
	cli.StringFlag{
		Name:  "stdin",
		Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
	},
	cli.BoolFlag{
		Name:  "wait-for-completion",
		Usage: "[optional] If provided, the command will wait to return until the volume has been deleted.",
	},
}

func (c *commandDelete) HandleFlags() error {
	wait := false
	if c.IsSet("wait-for-completion") {
		wait = true
	}

	return nil
}

func (c *commandDelete) HandleSingle() error {
	id, err := c.IDOrName(volumes.IDFromName)
	if err != nil {
		return err
	}
	c.id = id
	return nil
}

func (c *commandDelete) Execute(r lib.Resourcer) (res lib.Resulter) {
	err := volumes.Delete(c.ServiceClient(), c.id).ExtractErr()
	if err != nil {
		res.SetError(err)
		return
	}

	if c.wait {
		i := 0
		for i < 120 {
			_, err := volumes.Get(c.ServiceClient(), c.id).Extract()
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		res.SetValue(fmt.Sprintf("Deleted volume [%s]\n", volumeID))
	} else {
		res.SetValue(fmt.Sprintf("Deleting volume [%s]\n", volumeID))
	}

	return
}

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"id", "name"}
}
