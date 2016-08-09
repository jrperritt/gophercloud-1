package container

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
)

type commandCreate struct {
	openstack.CommandUtil
	ContainerV1Command
	opts containers.CreateOptsBuilder
}

var (
	cCreate                   = new(commandCreate)
	_       lib.PipeCommander = cCreate
	_       lib.Waiter        = cCreate
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "[--name <name> | --stdin name]"),
	Description:  "Creates a container",
	Action:       actionCreate,
	Flags:        openstack.CommandFlags(cCreate),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsCreate) },
}

func actionCreate(ctx *cli.Context) {
	c := new(commandCreate)
	c.Context = ctx
	lib.Run(ctx, c)
}

func (c *commandCreate) Flags() []cli.Flag {
	return flagsCreate
}

func (c *commandCreate) Fields() []string {
	return []string{""}
}

var flagsCreate = []cli.Flag{
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `stdin` isn't provided] The name of the container",
	},
	cli.StringFlag{
		Name:  "stdin",
		Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
	},
	cli.StringFlag{
		Name:  "metadata",
		Usage: "[optional] Comma-separated key-value pairs for the container. Example: key1=val1,key2=val2",
	},
	cli.StringFlag{
		Name:  "container-read",
		Usage: "[optional] Comma-separated list of users for whom to grant read access to the container",
	},
	cli.StringFlag{
		Name:  "container-write",
		Usage: "[optional] Comma-separated list of users for whom to grant write access to the container",
	},
}

func (c *commandCreate) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	c.Quiet = c.Context.IsSet("quiet")

	opts := &containers.CreateOpts{
		ContainerRead:  c.Context.String("container-read"),
		ContainerWrite: c.Context.String("container-write"),
	}

	if c.Context.IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	c.opts = opts

	return nil
}

func (c *commandCreate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandCreate) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandCreate) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		var m map[string]interface{}
		err := containers.Create(c.ServiceClient, item.(string), c.opts).ExtractInto(&m)
		if err != nil {
			out <- err
			return
		}
		out <- m
	}
}

func (c *commandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}

func (c *commandCreate) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}
