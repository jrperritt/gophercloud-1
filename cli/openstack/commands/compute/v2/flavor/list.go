package flavor

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type commandList struct {
	FlavorCommand
	traits.Waitable
	traits.Fieldsable
	traits.Tableable
	opts flavors.ListOptsBuilder
}

var (
	cList                      = new(commandList)
	_     interfaces.Commander = cList
	_     interfaces.Tabler    = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists existing flavors",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *commandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  "min-disk",
			Usage: "[optional] Only list flavors that have at least this much disk storage (in GB).",
		},
		cli.IntFlag{
			Name:  "min-ram",
			Usage: "[optional] Only list flavors that have at least this much RAM (in GB).",
		}, cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing flavors at this flavor ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many flavors at most.",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *commandList) DefaultTableFields() []string {
	return []string{"id", "name", "ram", "disk", "swap", "vcpus", "rxtx_factor"}
}

func (c *commandList) HandleFlags() error {
	c.opts = &flavors.ListOpts{
		MinDisk: c.Context().Int("min-disk"),
		MinRAM:  c.Context().Int("min-ram"),
		Marker:  c.Context().String("marker"),
		Limit:   c.Context().Int("limit"),
	}
	return nil
}

func (c *commandList) Execute(_ interface{}, out chan interface{}) {
	err := flavors.ListDetail(c.ServiceClient(), c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp map[string][]map[string]interface{}
		err := (page.(flavors.FlavorPage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp["flavors"]
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
