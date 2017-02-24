package blockstorage

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/blockstorage/v2"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("block-storage")
	switch version {
	default:
		return v2.Get()
	}
}
