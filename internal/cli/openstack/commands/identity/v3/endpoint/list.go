package endpoint

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type commandList struct {
	EndpointV3Command
	traits.Fieldsable
	traits.Tableable
	opts endpoints.ListOptsBuilder
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
	Description:  "Lists service endpoints",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *commandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "service-id",
			Usage: "[optional] Only list endpoints that have this service ID",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *commandList) DefaultTableFields() []string {
	return []string{"id", "name", "service_id", "url", "interface", "region"}
}

func (c *commandList) HandleFlags() error {
	opts := new(endpoints.ListOpts)
	opts.ServiceID = c.Context().String("service-id")
	c.opts = opts
	return nil
}

func (c *commandList) Execute(_ interface{}, out chan interface{}) {
	err := endpoints.List(c.ServiceClient(), c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp struct {
			Endpoints []map[string]interface{} `json:"endpoints"`
		}
		err := (page.(endpoints.EndpointPage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp.Endpoints
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
