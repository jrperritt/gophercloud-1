package service

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type commandList struct {
	ServiceV3Command
	traits.Fieldsable
	traits.Tableable
	opts services.ListOptsBuilder
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
	Description:  "Lists existing services",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *commandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "[optional] Only list services that have this type",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *commandList) DefaultTableFields() []string {
	return []string{"id", "name", "type"}
}

func (c *commandList) HandleFlags() error {
	opts := new(services.ListOpts)
	opts.ServiceType = c.Context().String("type")
	c.opts = opts
	return nil
}

func (c *commandList) Execute(_ interface{}, out chan interface{}) {
	err := services.List(c.ServiceClient(), c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp struct {
			Services []map[string]interface{} `json:"services"`
		}
		err := (page.(services.ServicePage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp.Services
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
