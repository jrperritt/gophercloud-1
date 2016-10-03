package subnet

import (
	"fmt"
	"strings"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	SubnetV2Command
	commands.DataResp
	opts subnets.CreateOptsBuilder
}

var (
	cCreate = new(CommandCreate)

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--network-id <ID> --cidr <CIDR>"),
	Description:  "Creates a subnet",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[required] The network ID under which to create this subnet.",
		},
		cli.StringFlag{
			Name:  "cidr",
			Usage: "[required] The CIDR of this subnet.",
		},
		cli.IntFlag{
			Name:  "ip-version",
			Usage: "[required] The IP version this subnet should have. Options are: 4, 6.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the network",
		},
		cli.StringFlag{
			Name:  "gateway-ip",
			Usage: "[optional] The gateway IP address this subnet should have.",
		},
		cli.BoolFlag{
			Name:  "enable-dhcp",
			Usage: "[optional] If set, DHCP will be enabled on this subnet.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant who should own this network.",
		},
		cli.StringSliceFlag{
			Name: "allocation-pool",
			Usage: strings.Join([]string{"[optional] An allocation pool for this subnet. This flag may be provided several times.\n",
				"\tEach one of these flags takes 2 values: start and end.\n",
				"\tExamle: --allocation-pool start=192.0.2.1,end=192.0.2.254 --allocation-pool start:172.20.0.1,end=172.20.0.254"}, ""),
		},
		cli.StringFlag{
			Name:  "dns-nameservers",
			Usage: "[optional] A comma-separated list of DNS Nameservers for this subnet.",
		},
		cli.BoolFlag{
			Name: "no-gateway",
			Usage: "[optional] If provided, no gateway will be associated with the subnet.\n" +
				"\tBy default, a default gateway will be provisioned.",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"network-id", "cidr", "ip-version"})
	if err != nil {
		return err
	}

	opts := &subnets.CreateOpts{
		NetworkID: c.Context.String("network-id"),
		CIDR:      c.Context.String("cidr"),
		Name:      c.Context.String("name"),
		IPVersion: gophercloud.IPVersion(c.Context.Int("ip-version")),
		TenantID:  c.Context.String("tenant-id"),
	}

	if c.Context.IsSet("gateway-ip") && c.Context.IsSet("no-gateway") {
		return fmt.Errorf("Only one of gateway-ip and no-gateway may be provided")
	}

	if c.Context.IsSet("gateway-ip") || c.Context.IsSet("no-gateway") {
		gatewayIP := c.Context.String("gateway-ip")
		opts.GatewayIP = &gatewayIP
	}

	opts.EnableDHCP = gophercloud.Disabled
	if c.Context.IsSet("enable-dhcp") {
		opts.EnableDHCP = gophercloud.Enabled
	}

	if c.Context.IsSet("dns-nameservers") {
		opts.DNSNameservers = strings.Split(c.Context.String("dns-nameservers"), ",")
	}

	if c.Context.IsSet("allocation-pool") {
		allocationPoolsRaw := c.Context.StringSlice("allocation-pool")
		allocationPoolsRawSlice, err := c.ValidateStructFlag(allocationPoolsRaw)
		if err != nil {
			return err
		}
		allocationPools := make([]subnets.AllocationPool, len(allocationPoolsRawSlice))
		for i, allocationPoolMap := range allocationPoolsRawSlice {
			allocationPools[i] = subnets.AllocationPool{
				Start: allocationPoolMap["start"].(string),
				End:   allocationPoolMap["end"].(string),
			}
		}
		opts.AllocationPools = allocationPools
	}

	c.opts = opts

	return nil
}

func (c *CommandCreate) Execute(_ interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := subnets.Create(c.ServiceClient, c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["subnet"]
	default:
		out <- err
	}
}
