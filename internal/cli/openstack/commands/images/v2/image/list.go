package image

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type commandList struct {
	ImageV2Command
	traits.Fieldsable
	traits.Tableable
	opts images.ListOptsBuilder
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
	Description:  "Lists existing images",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *commandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "tag",
			Usage: "[optional] Only list images that have this tag",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "[optional] Only list images that have this status",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing images at this flavor ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many images at most.",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *commandList) DefaultTableFields() []string {
	return []string{"id", "name", "size", "container_format", "disk_format", "min_disk", "min_ram", "rxtx_factor"}
}

func (c *commandList) HandleFlags() error {
	opts := new(images.ListOpts)
	opts.Tag = c.Context().String("tag")
	opts.Marker = c.Context().String("marker")
	opts.Limit = c.Context().Int("limit")
	opts.Status = images.ImageStatus(c.Context().String("status"))

	c.opts = opts
	return nil
}

func (c *commandList) Execute(_ interface{}, out chan interface{}) {
	err := images.List(c.ServiceClient(), c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp struct {
			Images []map[string]interface{} `json:"images"`
		}
		err := (page.(images.ImagePage)).ExtractInto(&tmp)
		if err != nil {
			return false, err
		}
		out <- tmp.Images
		return true, nil
	})
	if err != nil {
		out <- err
	}
}
