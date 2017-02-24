package account

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/accounts"
	"gopkg.in/urfave/cli.v1"
)

type CommandUpdate struct {
	AccountV1Command
	traits.Fieldsable
	opts accounts.UpdateOptsBuilder
}

var (
	cUpdate = new(CommandUpdate)

	flagsUpdate = openstack.CommandFlags(cUpdate)
)

var update = cli.Command{
	Name:         "update",
	Usage:        util.Usage(commandPrefix, "update", ""),
	Description:  "Updates account metadata",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpdate) },
	Flags:        flagsUpdate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpdate) },
}

func (c *CommandUpdate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
		cli.StringFlag{
			Name:  "temp-url-key",
			Usage: "[optional] The first temporary URL key",
		},
		cli.StringFlag{
			Name:  "temp-url-key-2",
			Usage: "[optional] The back-up temporary URL key",
		},
	}
}

func (c *CommandUpdate) HandleFlags() (err error) {
	opts := new(accounts.UpdateOpts)
	opts.TempURLKey = c.Context().String("temp-url-key")
	opts.TempURLKey2 = c.Context().String("temp-url-key-2")

	if c.Context().IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	c.opts = opts
	return
}

func (c *CommandUpdate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := accounts.Update(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}
