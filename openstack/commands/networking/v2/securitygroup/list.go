package securitygroup

import (
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	SecurityGroupV2Command
	commands.Waitable
	commands.DataResp
	opts groups.ListOpts
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
	Description:  "Lists security groups",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list security groups with this name.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] Only list security groups that are owned by this tenant ID.",
		},
	}
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "name", "tenant_id"}
}

func (c *CommandList) HandleFlags() error {
	c.opts = groups.ListOpts{
		Name:     c.Context.String("name"),
		TenantID: c.Context.String("tenant-id"),
	}

	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	err := groups.List(c.ServiceClient, c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp map[string][]map[string]interface{}
		err := (page.(groups.SecGroupPage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp["security_groups"]
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
