package container

import (
	"fmt"

	"github.com/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"gopkg.in/urfave/cli.v1"
)

type CommandEmpty struct {
	ContainerV1Command
	traits.Waitable
}

var (
	cEmpty                          = new(CommandEmpty)
	_      interfaces.PipeCommander = cEmpty

	flagsEmpty = openstack.CommandFlags(cEmpty)
)

var empty = cli.Command{
	Name:         "empty",
	Usage:        util.Usage(commandPrefix, "empty", "[--name <containerName> | --stdin name]"),
	Description:  "Deletes all objects in a container",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cEmpty) },
	Flags:        flagsEmpty,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsEmpty) },
}

func (c *CommandEmpty) Flags() []cli.Flag {
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

func (c *CommandEmpty) Fields() []string {
	return []string{""}
}

func (c *CommandEmpty) HandleFlags() error {
	return nil
}

func (c *CommandEmpty) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandEmpty) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandEmpty) Execute(item interface{}, out chan interface{}) {
	err := handleEmpty(cEmpty.ContainerV1Command, item.(string))
	if err != nil {
		out <- fmt.Errorf("Error emptying container [%s]: %s", item.(string), err)
		return
	}
	out <- fmt.Sprintf("Successfully emptied container [%s]", item.(string))
}

func (c *CommandEmpty) PipeFieldOptions() []string {
	return []string{"name"}
}
