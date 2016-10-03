package port

import (
	"fmt"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	PortV2Command
	commands.Waitable
	commands.DataResp
	opts ports.ListOptsBuilder
}

var (
	cList                     = new(CommandList)
	_     openstack.Commander = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(CommandPrefix, "list", ""),
	Description:  "Lists existing ports",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list ports with this name.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "[optional] Only list ports that are up or not. Options are: true, false.",
		},
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[optional] Only list ports with this network ID.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "[optional] Only list ports that have this status.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list ports that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "device-id",
			Usage: "[optional] Only list ports with this device ID.",
		},
	}
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "name", "network_id", "status", "mac_address", "device_id"}
}

func (c *CommandList) HandleFlags() error {
	opts := &ports.ListOpts{
		Name:      c.Context.String("name"),
		NetworkID: c.Context.String("network-id"),
		DeviceID:  c.Context.String("device-id"),
		TenantID:  c.Context.String("tenant-id"),
		Status:    c.Context.String("status"),
	}

	if c.Context.IsSet("up") {
		switch c.Context.String("up") {
		case "true":
			opts.AdminStateUp = gophercloud.Enabled
		case "false":
			opts.AdminStateUp = gophercloud.Disabled
		default:
			return fmt.Errorf("Invalid value for flag `up`: %s. Options are: true, false", c.Context.String("up"))
		}
	}

	c.Wait = c.Context.IsSet("wait")
	c.opts = opts

	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	err := ports.List(c.ServiceClient, c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp map[string][]map[string]interface{}
		err := (page.(ports.PortPage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp["ports"]
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
