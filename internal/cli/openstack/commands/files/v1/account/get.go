package account

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/accounts"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	AccountV1Command
	traits.Fieldsable
}

var (
	cGet = new(CommandGet)

	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", ""),
	Description:  "Gets account metadata",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Execute(_ interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := accounts.Get(c.ServiceClient(), nil).ExtractInto(&m)
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}
