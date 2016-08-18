package object

import (
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type commandList struct {
	openstack.CommandUtil
	ObjectV1Command
	opts     objects.ListOptsBuilder
	allPages bool
}

var (
	cList                   = new(commandList)
	_     lib.PipeCommander = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists objects in a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsList) },
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
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all objects. Default is to paginate.",
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

func (c *commandList) Fields() []string {
	return []string{"name", "bytes", "hash", "content_type", "last_modified"}
}

func (c *commandList) HandleFlags() error {
	c.opts = &objects.ListOpts{
		Full:      c.Context.Bool("full"),
		Prefix:    c.Context.String("prefix"),
		EndMarker: c.Context.String("end-marker"),
		Marker:    c.Context.String("marker"),
		Limit:     c.Context.Int("limit"),
	}
	c.allPages = c.Context.IsSet("all-pages")
	return nil
}

func (c *commandList) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandList) HandleSingle() (interface{}, error) {
	return c.Context.String("container"), c.CheckFlagsSet([]string{"container"})
}

func (c *commandList) Execute(in, out chan interface{}) {
	defer close(out)
	c.opts.(*objects.ListOpts).Full = true
	for item := range in {
		pager := objects.List(c.ServiceClient, item.(string), c.opts)
		switch c.allPages {
		case true:
			page, err := pager.AllPages()
			if err != nil {
				out <- err
				return
			}
			var tmp []map[string]interface{}
			err = (page.(objects.ObjectPage)).ExtractInto(&tmp)
			switch err {
			case nil:
				out <- tmp
			default:
				out <- err
			}
		default:
			err := pager.EachPage(func(page pagination.Page) (bool, error) {
				var tmp []map[string]interface{}
				err := (page.(objects.ObjectPage)).ExtractInto(&tmp)
				if err != nil {
					return false, err
				}
				out <- tmp
				return true, nil
			})
			if err != nil {
				out <- err
			}
		}
	}
}

func (c *commandList) PipeFieldOptions() []string {
	return []string{"container"}
}

func (c *commandList) PreTable(_ interface{}) error {

	return nil
}
