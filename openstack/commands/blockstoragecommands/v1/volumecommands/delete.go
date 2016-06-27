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
	openstack.CommandUtil
	VolumeV1Command
	wait bool
}

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--id <volumeID> | --name <volumeName> | --stdin id]"),
	Description:  "Deletes a volume",
	Action:       actionDelete,
	Flags:        openstack.CommandFlags(flagsDelete, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDelete) },
}

func actionDelete(ctx *cli.Context) {
	c := new(commandDelete)
	c.Context = ctx
	lib.Run(ctx, c)
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
		Name:  "wait",
		Usage: "[optional] If provided, the command will wait to return until the volume has been deleted.",
	},
}

func (c *commandDelete) HandleFlags() error {
	c.wait = c.Context.IsSet("wait")
	return nil
}

func (c *commandDelete) HandleSingle() (interface{}, error) {
	id, err := c.IDOrName(volumes.IDFromName)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (c *commandDelete) Execute(item interface{}, out chan interface{}) {
	defer func() {
		close(out)
	}()
	id := item.(string)
	err := volumes.Delete(c.ServiceClient, id).ExtractErr()
	if err != nil {
		out <- err
		return
	}

	switch c.wait {
	case true:
		i := 0
		for i < 120 {
			_, err := volumes.Get(c.ServiceClient, id).Extract()
			if err != nil {
				break
			}
			time.Sleep(5 * time.Second)
			i++
		}
		out <- fmt.Sprintf("Deleted volume [%s]\n", id)
	default:
		out <- fmt.Sprintf("Deleting volume [%s]\n", id)
	}
}

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"id"}
}
