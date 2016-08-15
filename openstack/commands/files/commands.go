package files

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/gophercloud/cli/openstack/commands/files/v1"
	"github.com/gophercloud/cli/util"
)

func Get() []cli.Command {
	version := util.GetVersion("files")
	switch version {
	default:
		return v1.Get()
	}
}
