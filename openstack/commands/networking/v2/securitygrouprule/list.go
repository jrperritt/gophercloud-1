package securitygrouprule

import (
	"fmt"

	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	SecurityGroupRuleV2Command
	traits.Waitable
	traits.DataResp
	opts rules.ListOpts
}

var (
	cList                     = new(CommandList)
	_     openstack.Waiter    = cList
	_     openstack.Commander = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists security groups rules",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "direction",
			Usage: "[optional] Only list security group rules with this direction. Options are: ingress, egress.",
		},
		cli.StringFlag{
			Name:  "ether-type",
			Usage: "[optional] Only list security group rules with this ether type. Options are: ipv4, ipv6.",
		},
		cli.IntFlag{
			Name:  "port-range-min",
			Usage: "[optional] Only list security group rules that have the low port greater than this.",
		},
		cli.IntFlag{
			Name:  "port-range-max",
			Usage: "[optional] Only list security group rules that have the high port less than this.",
		},
		cli.StringFlag{
			Name:  "protocol",
			Usage: "[optional] Only list security group rules with this protocol. Examples: tcp, udp, icmp.",
		},
		cli.StringFlag{
			Name:  "security-group-id",
			Usage: "[optional] Only list security group rules with this security group ID.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list security group rules that are owned by the tenant with this tenant ID.",
		},
	}
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "direction", "ethertype", "port_range_min", "port_range_max", "protocol", "security_group_id"}
}

func (c *CommandList) HandleFlags() error {
	opts := rules.ListOpts{
		Direction:    c.Context.String("direction"),
		PortRangeMax: c.Context.Int("port-range-max"),
		PortRangeMin: c.Context.Int("port-range-min"),
		Protocol:     c.Context.String("protocol"),
		SecGroupID:   c.Context.String("security-group-id"),
		TenantID:     c.Context.String("tenant-id"),
	}

	if c.Context.IsSet("ether-type") {
		etherType := c.Context.String("ether-type")
		switch etherType {
		case "ipv4":
			opts.EtherType = string(rules.EtherType4)
		case "ipv6":
			opts.EtherType = string(rules.EtherType6)
		default:
			return fmt.Errorf("Invalid value for `ether-type`: %s. Options are: ipv4, ipv6", etherType)
		}
	}

	c.opts = opts

	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	err := rules.List(c.ServiceClient, c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp map[string][]map[string]interface{}
		err := (page.(rules.SecGroupRulePage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp["security_group_rules"]
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
