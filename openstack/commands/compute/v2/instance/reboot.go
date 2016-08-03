package instance

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type commandReboot struct {
	openstack.CommandUtil
	InstanceV2Command
	opts servers.RebootOptsBuilder
	*openstack.Progress
}

var (
	cReboot                   = new(commandReboot)
	_       lib.PipeCommander = cReboot
	_       lib.Progresser    = cReboot
	_       lib.Waiter        = cReboot
)

var reboot = cli.Command{
	Name:         "reboot",
	Usage:        util.Usage(commandPrefix, "reboot", "[--id <serverID> | --name <serverName> | --stdin id] [--soft | --hard]"),
	Description:  "Reboots a server",
	Action:       actionReboot,
	Flags:        openstack.CommandFlags(flagsReboot, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsReboot) },
}

func actionReboot(ctx *cli.Context) {
	c := new(commandReboot)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsReboot = []cli.Flag{
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
	cli.BoolFlag{
		Name:  "wait",
		Usage: "[optional] If provided, will wait to return until the server has been rebooted.",
	},
}

func (c *commandReboot) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")

	switch c.Context.IsSet("hard") {
	case true:
		switch c.Context.IsSet("soft") {
		case true:
			return fmt.Errorf("Only one of either --soft or --hard may be provided.")
		default:
			c.opts = &servers.RebootOpts{servers.HardReboot}
		}
	default:
		switch c.Context.IsSet("soft") {
		case true:
			c.opts = &servers.RebootOpts{servers.SoftReboot}
		default:
			return fmt.Errorf("One of either --soft or --hard must be provided.")
		}
	}

	return nil
}

func (c *commandReboot) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandReboot) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandReboot) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		item := item
		go func() {
			id := item.(string)
			err := servers.Reboot(c.ServiceClient, id, c.opts).ExtractErr()
			if err != nil {
				out <- err
				return
			}
			switch c.Wait {
			case true:
				out <- id
			default:
				out <- fmt.Sprintf("Rebooting server [%s]", id)
			}
		}()
	}
}

func (c *commandReboot) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandReboot) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}

func (c *commandReboot) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Rebooting"
	c.Progress.DoneMsg = "Rebooted"
	c.Progress.Start()
}

func (c *commandReboot) ShowProgress(in, out chan interface{}) {
	id := (<-in).(string)

	c.StartBar(&openstack.ProgressStatus{
		Name:      id,
		StartTime: time.Now(),
	})

	err := util.WaitFor(900, func() (bool, error) {
		_, err := servers.Get(c.ServiceClient, id).Extract()
		if err != nil {
			c.CompleteBar(&openstack.ProgressStatus{
				Name: id,
			})
			out <- fmt.Sprintf("Rebooted server [%s]", id)
			return true, nil
		}

		c.UpdateBar(&openstack.ProgressStatus{
			Name: id,
		})
		return false, nil
	})

	if err != nil {
		c.ErrorBar(&openstack.ProgressStatus{
			Name: id,
			Err:  err,
		})
		out <- err
	}
}
