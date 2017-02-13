package loadbalancer

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	LoadbalancerV2Command
	traits.Waitable
	traits.Fieldsable
	traits.Tableable
	opts loadbalancers.ListOptsBuilder
}

var (
	cList                       = new(CommandList)
	_         interfaces.Waiter = cList
	_         interfaces.Tabler = cList
	flagsList                   = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists load balancers",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name: "id",
		},
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
		cli.StringFlag{
			Name: "vip-port-id",
		},
		cli.StringFlag{
			Name: "admin-state-up",
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
		cli.StringFlag{
			Name: "provisioning-status",
		},
		cli.StringFlag{
			Name: "operating-status",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "name", "vip-address", "vip-port-id"}
}

func (c *CommandList) HandleFlags() error {
	opts := &loadbalancers.ListOpts{
		VipSubnetID:        c.Context().String("vip-subnet-id"),
		Name:               c.Context().String("name"),
		Description:        c.Context().String("description"),
		VipAddress:         c.Context().String("vip-address"),
		VipPortID:          c.Context().String("vip-port-id"),
		TenantID:           c.Context().String("tenant-id"),
		Flavor:             c.Context().String("flavor"),
		Provider:           c.Context().String("provider"),
		ID:                 c.Context().String("id"),
		ProvisioningStatus: c.Context().String("provisioning-status"),
		OperatingStatus:    c.Context().String("operating-status"),
	}

	if c.Context().IsSet("admin-state-up") {
		switch c.Context().String("admin-state-up") {
		case "true":
			opts.AdminStateUp = gophercloud.Enabled
		case "false":
			opts.AdminStateUp = gophercloud.Disabled
		default:
		}
	}

	c.opts = opts

	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	err := loadbalancers.List(c.ServiceClient(), c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp map[string][]map[string]interface{}
		err := (page.(loadbalancers.LoadBalancerPage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp["loadbalancers"]
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
