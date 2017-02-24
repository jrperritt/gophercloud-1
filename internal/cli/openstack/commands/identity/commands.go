package identity

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack/commands/identity/v3"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	cli "gopkg.in/urfave/cli.v1"
)

func Get() []cli.Command {
	version := util.GetVersion("identity")
	switch version {
	default:
		return v3.Get()
	}
}
