package server

import (
	"fmt"

	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/mvalkon/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

type CommandResize struct {
	ServerV2Command
	traits.Waitable
	traits.Pipeable
	traits.TextProgressable
	opts servers.ResizeOptsBuilder
}

var (
	cResize                          = new(CommandResize)
	_       interfaces.PipeCommander = cResize
	_       interfaces.Progresser    = cResize

	flagsResize = openstack.CommandFlags(cResize)
)

var resize = cli.Command{
	Name:         "resize",
	Usage:        util.Usage(CommandPrefix, "resize", "[--id <serverID> | --name <serverName> | --stdin id] [--flavor-id | --flavor-name]"),
	Description:  "Resizes a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cResize) },
	Flags:        flagsResize,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsResize) },
}

func (c *CommandResize) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the server.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[optional; required if `flavor-name` is not provided] The ID of the flavor that the resized server should have.",
		},
		cli.StringFlag{
			Name:  "flavor-name",
			Usage: "[optional; required if `flavor-id` is not provided] The name of the flavor that the resized server should have.",
		},
	}
}

func (c *CommandResize) HandleFlags() error {
	opts := new(servers.ResizeOpts)

	if c.Context().IsSet("flavor-id") {
		opts.FlavorRef = c.Context().String("flavor-id")
		c.opts = opts
		return nil
	}

	if c.Context().IsSet("flavor-name") {
		id, err := flavors.IDFromName(c.ServiceClient(), c.Context().String("flavor-name"))
		if err != nil {
			return err
		}
		opts.FlavorRef = id
		c.opts = opts
		return nil
	}

	return fmt.Errorf("One and only one of flavor-name and flavor-id must be provided")
}

func (c *CommandResize) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandResize) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandResize) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	err := servers.Resize(c.ServiceClient(), id, c.opts).ExtractErr()
	if err != nil {
		out <- err
		return
	}
	switch c.ShouldWait() {
	case true:
		out <- id
	default:
		out <- fmt.Sprintf("Resizing server [%s]", id)
	}
}

func (c *CommandResize) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *CommandResize) WaitFor(raw interface{}, out chan<- interface{}) {
	id := raw.(string)

	err := gophercloud.WaitFor(900, func() (bool, error) {
		var m map[string]map[string]interface{}
		err := servers.Get(c.ServiceClient(), id).ExtractInto(&m)
		if err != nil {
			return false, err
		}
		switch m["server"]["status"].(string) {
		case "ACTIVE":
			out <- fmt.Sprintf("Resized server [%s]", id)
			return true, nil
		default:
			//c.ProgUpdateChIn() <- m["server"]["status"]
			return false, nil
		}
	})

	if err != nil {
		out <- err
	}
}

func (c *CommandResize) InitProgress() {
	c.TextProgressable.InitProgress()
}

func (c *CommandResize) RunningMsg() string {
	return "Resizing"
}

func (c *CommandResize) DoneMsg() string {
	return "Resized"
}
