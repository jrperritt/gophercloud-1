package network

import (
	"fmt"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	NetworkV2Command
	opts     networks.ListOptsBuilder
	allPages bool
}

var (
	cList                     = new(CommandList)
	_     openstack.Commander = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists existing networks",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all networks. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Only list networks with this name.",
		},
		cli.StringFlag{
			Name:  "up",
			Usage: "[optional] If provided, the network will be up. Options are: true, false",
		},
		cli.StringFlag{
			Name:  "shared",
			Usage: "[optional] If provided, the network is shared among all tenants. Options are: true, false",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "Only list networks that are owned by the tenant with this tenant ID.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "Only list networks that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing networks at this network ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many networks at most.",
		},
	}
}

func (c *CommandList) Fields() []string {
	return []string{""}
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "name", "admin_state_up", "status", "shared", "tenant_id"}
}

func (c *CommandList) HandleFlags() error {
	opts := &networks.ListOpts{
		Name:     c.Context.String("name"),
		TenantID: c.Context.String("tenant-id"),
		Status:   c.Context.String("status"),
		Marker:   c.Context.String("marker"),
		Limit:    c.Context.Int("limit"),
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

	if c.Context.IsSet("shared") {
		switch c.Context.String("shared") {
		case "true":
			opts.Shared = gophercloud.Enabled
		case "false":
			opts.Shared = gophercloud.Disabled
		default:
			return fmt.Errorf("Invalid value for flag `shared`: %s. Options are: true, false", c.Context.String("shared"))
		}
	}

	c.allPages = c.Context.IsSet("all-pages")
	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	pager := networks.List(c.ServiceClient, c.opts)
	switch c.allPages {
	case true:
		page, err := pager.AllPages()
		if err != nil {
			out <- err
			return
		}
		var tmp map[string][]map[string]interface{}
		err = (page.(networks.NetworkPage)).ExtractInto(&tmp)
		switch err {
		case nil:
			out <- tmp["networks"]
		default:
			out <- err
		}
	default:
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			var tmp map[string][]map[string]interface{}
			err := (page.(networks.NetworkPage)).ExtractInto(&tmp)
			if err != nil {
				return false, err
			}
			out <- tmp["networks"]
			return true, nil
		})
		if err != nil {
			out <- err
		}
	}
}
