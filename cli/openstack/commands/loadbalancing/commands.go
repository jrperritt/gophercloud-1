package loadbalancing

import (
	"github.com/gophercloud/cli/openstack/commands/loadbalancing/v2"
	"github.com/gophercloud/cli/util"
	"gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("load-balancing")
	switch version {
	default:
		return v2.Get()
	}
}
