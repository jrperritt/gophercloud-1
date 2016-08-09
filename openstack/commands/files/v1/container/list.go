package container

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/pagination"
)

type commandList struct {
	openstack.CommandUtil
	ContainerV1Command
	opts     containers.ListOptsBuilder
	allPages bool
}

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists existing containers",
	Action:       actionList,
	Flags:        openstack.CommandFlags(new(commandList)),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsList) },
}

func actionList(ctx *cli.Context) {
	c := new(commandList)
	c.Context = ctx
	lib.Run(ctx, c)
}

func (c *commandList) Flags() []cli.Flag {
	return flagsList
}

func (c *commandList) Fields() []string {
	return []string{"name", "count", "bytes"}
}

var flagsList = []cli.Flag{
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

func (c *commandList) HandleFlags() error {
	c.opts = &containers.ListOpts{
		Prefix:    c.Context.String("prefix"),
		EndMarker: c.Context.String("end-marker"),
		Marker:    c.Context.String("marker"),
		Limit:     c.Context.Int("limit"),
	}
	c.allPages = c.Context.IsSet("all-pages")
	return nil
}

func (c *commandList) Execute(_, out chan interface{}) {
	defer close(out)
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

func (c *commandList) PreTable(rawServers interface{}) error {
	if rawServers, ok := rawServers.([]map[string]interface{}); ok {
		for i, rawServer := range rawServers {
			for k, v := range rawServer {
				switch k {
				case "image":
					if imageMap, ok := v.(map[string]interface{}); ok {
						rawServer["image"] = imageMap["id"]
						rawServers[i] = rawServer
					}
				case "flavor":
					if flavorMap, ok := v.(map[string]interface{}); ok {
						rawServer["flavor"] = flavorMap["id"]
						rawServers[i] = rawServer
					}
				}
			}
		}
	}
	return nil
}
