package container

import (
	"fmt"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"gopkg.in/urfave/cli.v1"
)

type commandDelete struct {
	openstack.CommandUtil
	ContainerV1Command
}

var (
	cDelete                   = new(commandDelete)
	_       lib.PipeCommander = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--name <containerName> | --stdin name]"),
	Description:  "Deletes a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDelete) },
}

func (c *commandDelete) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name of the container.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
	}
}

func (c *commandDelete) Fields() []string {
	return []string{""}
}

func (c *commandDelete) HandleFlags() error {
	return nil
}

func (c *commandDelete) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandDelete) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandDelete) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		var m map[string]interface{}
		err := containers.Delete(c.ServiceClient, item.(string)).ExtractInto(&m)
		switch err {
		case nil:
			out <- fmt.Sprintf("Successfully deleted container [%s]", item.(string))
		default:
			out <- err
		}
	}
}

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"name"}
}
