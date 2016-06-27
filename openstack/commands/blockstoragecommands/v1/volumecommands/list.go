package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/gophercloud/pagination"
)

type commandList struct {
	openstack.CommandUtil
	VolumeV1Command
	opts volumes.ListOptsBuilder
}

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists existing volumes",
	Action:       actionList,
	Flags:        openstack.CommandFlags(flagsList, []string{""}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsList) },
}

func actionList(ctx *cli.Context) {
	c := new(commandList)
	c.Context = ctx
	lib.Run(ctx, c)
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

func (c *commandList) HandleFlags() error {
	c.opts = &volumes.ListOpts{
		Name:   c.Context.String("name"),
		Status: c.Context.String("status"),
	}
	return nil
}

func (c *commandList) Execute(_ interface{}, out chan interface{}) {
	pager := volumes.List(c.ServiceClient, c.opts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		var tmp map[string][]map[string]interface{}
		err := (page.(volumes.VolumePage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp["volumes"]
		return true, nil
	})
	if err != nil {
		out <- err
	}
	close(out)
}

func (c *commandList) Fields() []string {
	return []string{"id", "display_name", "bootable", "size", "status"}
}
