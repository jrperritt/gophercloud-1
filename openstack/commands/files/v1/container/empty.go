package container

import (
	"fmt"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"gopkg.in/urfave/cli.v1"
)

type commandEmpty struct {
	ContainerV1Command
}

var (
	cEmpty                   = new(commandEmpty)
	_      lib.PipeCommander = cEmpty

	flagsEmpty = openstack.CommandFlags(cEmpty)
)

var empty = cli.Command{
	Name:         "empty",
	Usage:        util.Usage(commandPrefix, "empty", "[--name <containerName> | --stdin name]"),
	Description:  "Deletes all objects in a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cEmpty) },
	Flags:        flagsEmpty,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsEmpty) },
}

func (c *commandEmpty) Flags() []cli.Flag {
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

func (c *commandEmpty) Fields() []string {
	return []string{""}
}

func (c *commandEmpty) HandleFlags() error {
	return nil
}

func (c *commandEmpty) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandEmpty) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *commandEmpty) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		item := item.(string)
		err := handleEmpty(cEmpty.ContainerV1Command, item)
		if err != nil {
			out <- fmt.Errorf("Error emptying container [%s]: %s", item, err)
			continue
		}
		out <- fmt.Sprintf("Successfully emptied container [%s]", item)
	}
}

func (c *commandEmpty) PipeFieldOptions() []string {
	return []string{"name"}
}
