package compute

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/gophercloud/cli/openstack/commands/compute/v2"
	"github.com/gophercloud/cli/util"
)

func Get() []cli.Command {
	version := util.GetVersion("compute")
	switch version {
	default:
		return v2.Get()
	}
}
