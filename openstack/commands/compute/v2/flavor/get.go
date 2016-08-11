package flavor

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

var _ lib.PipeCommander = new(commandGet)

type commandGet struct {
	openstack.CommandUtil
	FlavorV2Command
}

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "[--id <flavorID> | --name <flavorName> | --stdin id]"),
	Description:  "Retreives information about a flavor",
	Action:       actionGet,
	Flags:        openstack.CommandFlags(new(commandGet)),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsGet) },
}

func actionGet(ctx *cli.Context) {
	c := new(commandGet)
	c.Context = ctx
	lib.Run(ctx, c)
}

func (c *commandGet) Flags() []cli.Flag {
	return flagsGet
}

func (c *commandGet) Fields() []string {
	return []string{""}
}

var flagsGet = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the flavor.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the flavor.",
	},
	cli.StringFlag{
		Name:  "stdin",
		Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
	},
}

func (c *commandGet) HandleFlags() error {
	return nil
}

func (c *commandGet) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandGet) HandleSingle() (interface{}, error) {
	return c.IDOrName(flavors.IDFromName)
}

func (c *commandGet) Execute(in, out chan interface{}) {
	defer close(out)

	for item := range in {
		item := item
		var m map[string]map[string]interface{}
		err := flavors.Get(c.ServiceClient, item.(string)).ExtractInto(&m)
		switch err {
		case nil:
			out <- m["flavor"]
		default:
			out <- err
		}
	}
}

func (c *commandGet) PipeFieldOptions() []string {
	return []string{"id"}
}
