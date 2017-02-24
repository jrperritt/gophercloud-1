package server

import (
	"fmt"

	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/mvalkon/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

type CommandReboot struct {
	ServerV2Command
	traits.Waitable
	traits.Pipeable
	traits.TextProgressable
	traits.Fieldsable
	opts servers.RebootOptsBuilder
}

var (
	cReboot                          = new(CommandReboot)
	_       interfaces.PipeCommander = cReboot
	_       interfaces.Progresser    = cReboot

	flagsReboot = openstack.CommandFlags(cReboot)
)

var reboot = cli.Command{
	Name:         "reboot",
	Usage:        util.Usage(CommandPrefix, "reboot", "[--id <serverID> | --name <serverName> | --stdin id] [--soft | --hard]"),
	Description:  "Reboots a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cReboot) },
	Flags:        flagsReboot,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsReboot) },
}

func (c *CommandReboot) Flags() []cli.Flag {
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
		cli.BoolFlag{
			Name:  "soft",
			Usage: "[optional; required if 'hard' is not provided] Ask the OS to restart under its own procedures.",
		},
		cli.BoolFlag{
			Name:  "hard",
			Usage: "[optional; required if 'soft' is not provided] Cut power to the machine and then restore it after a brief while.",
		},
	}
}

func (c *CommandReboot) HandleFlags() error {
	switch c.Context().IsSet("hard") {
	case true:
		switch c.Context().IsSet("soft") {
		case true:
			return fmt.Errorf("Only one of either --soft or --hard may be provided.")
		default:
			c.opts = &servers.RebootOpts{Type: servers.HardReboot}
		}
	default:
		switch c.Context().IsSet("soft") {
		case true:
			c.opts = &servers.RebootOpts{Type: servers.SoftReboot}
		default:
			return fmt.Errorf("One of either --soft or --hard must be provided.")
		}
	}

	return nil
}

func (c *CommandReboot) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandReboot) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandReboot) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	err := servers.Reboot(c.ServiceClient(), id, c.opts).ExtractErr()
	if err != nil {
		out <- err
		return
	}
	switch c.ShouldWait() || c.ShouldProgress() {
	case true:
		out <- id
	default:
		out <- fmt.Sprintf("Rebooting server [%s]", id)
	}
}

func (c *CommandReboot) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *CommandReboot) WaitFor(raw interface{}, out chan<- interface{}) {
	id := raw.(string)

	err := gophercloud.WaitFor(900, func() (bool, error) {
		var m map[string]map[string]interface{}
		err := servers.Get(c.ServiceClient(), id).ExtractInto(&m)
		if err != nil {
			return false, err
		}
		switch m["server"]["status"].(string) {
		case "ACTIVE":
			out <- fmt.Sprintf("Rebooted server [%s]", id)
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

func (c *CommandReboot) InitProgress() {
	c.TextProgressable.InitProgress()
}

func (c *CommandReboot) RunningMsg() string {
	return "Rebooting"
}

func (c *CommandReboot) DoneMsg() string {
	return "Rebooted"
}
