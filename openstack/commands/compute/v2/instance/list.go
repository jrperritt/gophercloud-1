package instance

import (
	"reflect"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	ServerV2Command
	opts     servers.ListOptsBuilder
	allPages bool
}

var (
	cList     = new(CommandList)
	flagsList = openstack.CommandFlags(cList)
	list      = cli.Command{
		Name:         "list",
		Usage:        util.Usage(commandPrefix, "list", ""),
		Description:  "Lists existing servers",
		Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
		Flags:        flagsList,
		BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
	}
)

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "all-pages",
			Usage: "[optional] Return all servers. Default is to paginate.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] Only list servers with this name.",
		},
		cli.StringFlag{
			Name:  "changes-since",
			Usage: "[optional] Only list servers that have been changed since this time/date stamp.",
		},
		cli.StringFlag{
			Name:  "image",
			Usage: "[optional] Only list servers that have this image ID.",
		},
		cli.StringFlag{
			Name:  "flavor",
			Usage: "[optional] Only list servers that have this flavor ID.",
		},
		cli.StringFlag{
			Name:  "status",
			Usage: "[optional] Only list servers that have this status.",
		},
		cli.StringFlag{
			Name:  "marker",
			Usage: "[optional] Start listing servers at this server ID.",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "[optional] Only return this many servers at most.",
		},
		cli.BoolFlag{
			Name:  "all-tenants",
			Usage: "[optional] If provided, will show servers from all tenants",
		},
	}
}

func (c *CommandList) Fields() []string {
	return []string{""}
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "name", "status", "accessIPv4", "image", "flavor"}
}

func (c *CommandList) HandleFlags() error {
	c.opts = &servers.ListOpts{
		ChangesSince: c.Context.String("changes-since"),
		Image:        c.Context.String("image-name"),
		Flavor:       c.Context.String("flavor-name"),
		Name:         c.Context.String("name"),
		Status:       c.Context.String("status"),
		Host:         c.Context.String("host"),
		Marker:       c.Context.String("marker"),
		Limit:        c.Context.Int("limit"),
		AllTenants:   c.Context.IsSet("all-tenants"),
	}
	c.allPages = c.Context.IsSet("all-pages")
	return nil
}

func (c *CommandList) Execute(_ interface{}, out chan interface{}) {
	pager := servers.List(c.ServiceClient, c.opts)
	switch c.allPages {
	case true:
		page, err := pager.AllPages()
		if err != nil {
			out <- err
			return
		}
		var tmp map[string][]map[string]interface{}
		err = (page.(servers.ServerPage)).ExtractInto(&tmp)
		switch err {
		case nil:
			out <- tmp["servers"]
		default:
			out <- err
		}
	default:
		err := pager.EachPage(func(page pagination.Page) (bool, error) {
			var tmp map[string][]map[string]interface{}
			err := (page.(servers.ServerPage)).ExtractInto(&tmp)
			if err != nil {
				return false, err
			}
			out <- tmp["servers"]
			return true, nil
		})
		if err != nil {
			out <- err
		}
	}
}

func (c *CommandList) PreTable(rawServers interface{}) error {
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
				switch reflect.ValueOf(v).Kind() {
				case reflect.Map, reflect.Slice, reflect.Struct:
					delete(rawServer, k)
				}
			}
		}
	}
	return nil
}
