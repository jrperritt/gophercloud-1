package blockstorage

import (
	"github.com/gophercloud/cli/openstack/commands/blockstorage/v2"
	"github.com/gophercloud/cli/util"
	"gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("block-storage")
	switch version {
	default:
		return v2.Get()
	}
}
