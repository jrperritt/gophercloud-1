package compute

import (
	"github.com/gophercloud/gophercloud/cli/openstack/commands/compute/v2"
	"github.com/gophercloud/gophercloud/cli/util"
	"gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("compute")
	switch version {
	default:
		return v2.Get()
	}
}
