package volumecommands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/gophercloud/gophercloud/pagination"
)

type commandList struct {
	openstack.Command
	opts volumes.ListOptsBuilder
}

func (c commandList) Name() string {
	return "list"
}

func (c commandList) Usage() string {
	return util.Usage(commandPrefix, "list", "")
}

func (c commandList) Description() string {
	return "Lists existing volumes"
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
		Name:   c.String("name"),
		Status: c.String("status"),
	}
	return nil
}

func (c *commandList) Execute(_ lib.Resourcer) lib.Resulter {
	r := new(openstack.Result)
	pager := volumes.List(c.ServiceClient(), c.opts)
	m := make([]map[string]interface{}, 0)
	//allPages, err := pager.AllPages()
	//if err != nil {
	//	r.SetError(err)
	//	return r
	//}
	fmt.Printf("pager :%+v\n", pager)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		var tmp []map[string]interface{}
		//info, err := volumes.ExtractVolumes(page)
		err := (page.(volumes.VolumePage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		fmt.Sprintf("tmp: %+v", tmp)
		//result := make([]volumes.Volume, len(info))
		//for j, volume := range info {
		//	result[j] = volume
		//}
		m = append(m, tmp...)
		return true, nil
	})
	if err != nil {
		r.SetError(err)
		return r
	}
	r.SetValue(m)
	return r
}
