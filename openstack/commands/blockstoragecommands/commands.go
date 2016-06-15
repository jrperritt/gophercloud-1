package blockstoragecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack/commands/blockstoragecommands/v1"
	"github.com/gophercloud/cli/util"
)

func Get() []cli.Command {
	version := util.GetVersion("block-storage")
	switch version {
	case "1":
		return v1.Get()
	case "2":
		return nil
	default:
		return nil
	}
}
