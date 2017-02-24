package loadbalancing

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/loadbalancing/v2"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("load-balancing")
	switch version {
	default:
		return v2.Get()
	}
}
