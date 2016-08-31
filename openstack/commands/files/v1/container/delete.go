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
	ContainerV1Command
	purge bool
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
		cli.BoolFlag{
			Name:  "purge",
			Usage: "[optional] Delete all objects in the container, and then delete the container.",
		},
	}
}

func (c *commandDelete) Fields() []string {
	return []string{""}
}

func (c *commandDelete) HandleFlags() error {
	c.purge = c.Context.IsSet("purge")
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
		item := item.(string)
		if c.purge {
			err := handleEmpty(cDelete.ContainerV1Command, item)
			if err != nil {
				out <- fmt.Errorf("Error purging container [%s]: %s", item, err)
				continue
			}
		}
		res := containers.Delete(c.ServiceClient, item)
		switch res.Err {
		case nil:
			out <- fmt.Sprintf("Successfully deleted container [%s]", item)
		default:
			out <- res.Err
		}
	}
}

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"name"}
}
