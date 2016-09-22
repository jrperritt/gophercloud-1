package container

import (
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	ContainerV1Command
	opts     containers.ListOptsBuilder
	allPages bool
}

var (
	cList                     = new(CommandList)
	_     openstack.Commander = cList

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists existing containers",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all containers. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "prefix",
			Usage: "[optional] Only return containers with this prefix.",
		},
		cli.StringFlag{
			Name:  "end-marker",
			Usage: "[optional] Only return containers with name less than this value.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing containers at this container name.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many containers at most.",
		},
	}
}

func (c *CommandList) Fields() []string {
	return []string{""}
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"name", "count", "bytes"}
}

func (c *CommandList) HandleFlags() error {
	c.opts = &containers.ListOpts{
		Prefix:    c.Context.String("prefix"),
		EndMarker: c.Context.String("end-marker"),
		Marker:    c.Context.String("marker"),
		Limit:     c.Context.Int("limit"),
	}
	c.allPages = c.Context.IsSet("all-pages")
	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	c.opts.(*containers.ListOpts).Full = true
	pager := containers.List(c.ServiceClient, c.opts)
	switch c.allPages {
	case true:
		page, err := pager.AllPages()
		if err != nil {
			out <- err
			return
		}
		var tmp []map[string]interface{}
		err = (page.(containers.ContainerPage)).ExtractInto(&tmp)
		switch err {
		case nil:
			out <- tmp
		default:
			out <- err
		}
	default:
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			var tmp []map[string]interface{}
			err := (page.(containers.ContainerPage)).ExtractInto(&tmp)
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
