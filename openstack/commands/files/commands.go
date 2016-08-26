package files

import (
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands/files/v1"
	"github.com/gophercloud/cli/util"
	"gopkg.in/urfave/cli.v1"
)

type FilesCommand struct {
	openstack.CommandUtil
}

func Get() []cli.Command {
	version := util.GetVersion("files")
	switch version {
	default:
		return v1.Get()
	}
}
