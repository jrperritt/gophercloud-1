package container

import (
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	ContainerV1Command
	traits.Waitable
	traits.DataResp
	opts containers.ListOptsBuilder
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

func (c *CommandList) DefaultTableFields() []string {
	return []string{"name", "count", "bytes"}
}

func (c *CommandList) HandleFlags() error {
	c.opts = &containers.ListOpts{
		Full:      true,
		Prefix:    c.Context.String("prefix"),
		EndMarker: c.Context.String("end-marker"),
		Marker:    c.Context.String("marker"),
		Limit:     c.Context.Int("limit"),
	}
	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	err := containers.List(c.ServiceClient, c.opts).EachPage(func(page pagination.Page) (bool, error) {
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
