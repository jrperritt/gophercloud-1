package files

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/files/v1"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("files")
	switch version {
	default:
		return v1.Get()
	}
}
