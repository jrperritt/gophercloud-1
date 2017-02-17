package images

import (
	"github.com/gophercloud/gophercloud/cli/openstack/commands/images/v2"
	"github.com/gophercloud/gophercloud/cli/util"
	cli "gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("images")
	switch version {
	default:
		return v2.Get()
	}
}
