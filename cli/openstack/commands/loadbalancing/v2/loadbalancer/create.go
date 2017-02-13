package loadbalancer

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	LoadbalancerV2Command
	traits.Fieldsable
	opts loadbalancers.CreateOptsBuilder
}

var (
	cCreate = new(CommandCreate)

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--vip-subnet-id <ID>"),
	Description:  "Creates a load balancer",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "vip-subnet-id",
		},
		cli.StringFlag{
			Name: "name",
		},
		cli.StringFlag{
			Name: "description",
		},
		cli.StringFlag{
			Name: "vip-address",
		},
		cli.BoolFlag{
			Name: "up",
		},
		cli.StringFlag{
			Name: "tenant-id",
		},
		cli.StringFlag{
			Name: "flavor",
		},
		cli.StringFlag{
			Name: "provider",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	opts := &loadbalancers.CreateOpts{
		VipSubnetID: c.Context().String("vip-subnet-id"),
		Name:        c.Context().String("name"),
		Description: c.Context().String("description"),
		VipAddress:  c.Context().String("vip-address"),
		TenantID:    c.Context().String("tenant-id"),
		Flavor:      c.Context().String("flavor"),
		Provider:    c.Context().String("provider"),
	}

	opts.AdminStateUp = gophercloud.Disabled
	if c.Context().IsSet("up") {
		opts.AdminStateUp = gophercloud.Enabled
	}

	c.opts = opts

	return nil
}

func (c *CommandCreate) Execute(_ interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := loadbalancers.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["loadbalancer"]
	default:
		out <- err
	}
}
