package securitygrouprule

import (
	"fmt"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	SecurityGroupRuleV2Command
	traits.Waitable
	traits.DataResp
	opts rules.CreateOptsBuilder
}

var (
	cCreate                   = new(CommandCreate)
	_       interfaces.Waiter = cCreate
	//	_       openstack.PipeCommander = cCreate

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--security-group-id <ID>"),
	Description:  "Creates a security group rule",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "security-group-id",
			Usage: "[required] The security group ID with which to associate this security group rule.",
		},
		cli.StringFlag{
			Name:  "direction",
			Usage: "[required] The direction of the security group rule. Options are: ingress, egress.",
		},
		cli.StringFlag{
			Name:  "ether-type",
			Usage: "[required] The ether type of the security group rule. Options are: ipv4, ipv6.",
		},
		cli.IntFlag{
			Name:  "port-range-min",
			Usage: "[optional] The minimum port of security group rule.",
		},
		cli.IntFlag{
			Name:  "port-range-max",
			Usage: "[optional] The maximum port of security group rule.",
		},
		cli.StringFlag{
			Name:  "protocol",
			Usage: "[optional] The protocol of the security group rule. Examples: tcp, udp, icmp.",
		},
		cli.StringFlag{
			Name:  "remote-ip-prefix",
			Usage: "[optional] The remote IP prefix to associate with this security group rule",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	if err := c.CheckFlagsSet([]string{"security-group-id", "ether-type", "direction"}); err != nil {
		return err
	}

	opts := &rules.CreateOpts{
		PortRangeMax:   c.Context().Int("port-range-max"),
		PortRangeMin:   c.Context().Int("port-range-min"),
		Protocol:       rules.RuleProtocol(c.Context().String("protocol")),
		SecGroupID:     c.Context().String("security-group-id"),
		RemoteIPPrefix: c.Context().String("remote-ip-prefix"),
	}

	direction := c.Context().String("direction")
	switch direction {
	case "ingress":
		opts.Direction = rules.DirIngress
	case "egress":
		opts.Direction = rules.DirEgress
	default:
		return fmt.Errorf("Invalid value for `direction`: %s. Options are: ingress, egress", direction)
	}

	if c.Context().IsSet("ether-type") {
		etherType := c.Context().String("ether-type")
		switch etherType {
		case "ipv4":
			opts.EtherType = rules.EtherType4
		case "ipv6":
			opts.EtherType = rules.EtherType6
		default:
			return fmt.Errorf("Invalid value for `ether-type`: %s. Options are: ipv4, ipv6", etherType)
		}
	}

	c.opts = opts

	return nil
}

func (c *CommandCreate) Execute(_ interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := rules.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["security_group_rule"]
	default:
		out <- err
	}
}
