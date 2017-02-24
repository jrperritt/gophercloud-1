package object

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type commandList struct {
	ObjectV1Command
	traits.Pipeable
	traits.Waitable
	traits.Fieldsable
	traits.Tableable
	opts objects.ListOptsBuilder
}

var (
	cList                          = new(commandList)
	_     interfaces.PipeCommander = cList
	_     interfaces.Tabler        = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists objects in a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *commandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `container` isn't provided] The field being piped into STDIN. Valid values are: container",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "[optional] Only return objects with this prefix.",
		},
		cli.StringFlag{
			Name:  "end-marker",
			Usage: "[optional] Only return objects with name less than this value.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing objects at this object name.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many objects at most.",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *commandList) DefaultTableFields() []string {
	return []string{"name", "bytes", "hash", "content_type", "last_modified"}
}

func (c *commandList) HandleFlags() error {
	c.opts = &objects.ListOpts{
		Full:      true,
		Prefix:    c.Context().String("prefix"),
		EndMarker: c.Context().String("end-marker"),
		Marker:    c.Context().String("marker"),
		Limit:     c.Context().Int("limit"),
	}
	return nil
}

func (c *commandList) HandleSingle() (interface{}, error) {
	return c.Context().String("container"), c.CheckFlagsSet([]string{"container"})
}

func (c *commandList) Execute(item interface{}, out chan interface{}) {
	err := objects.List(c.ServiceClient(), item.(string), c.opts).EachPage(func(page pagination.Page) (bool, error) {
		var tmp []map[string]interface{}
		err := (page.(objects.ObjectPage)).ExtractInto(&tmp)
		switch err {
		case nil:
			if len(tmp) > 0 {
				out <- tmp
			}
			return true, nil
		}
		return false, err
	})
	if err != nil {
		out <- err
	}
}

func (c *commandList) PipeFieldOptions() []string {
	return []string{"container"}
}
