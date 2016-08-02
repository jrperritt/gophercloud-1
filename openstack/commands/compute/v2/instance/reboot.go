package instance

import (
	"fmt"
	"sync"
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
	wait bool
	opts servers.RebootOptsBuilder
	*openstack.Progress
}

var (
	cReboot                   = new(commandReboot)
	_       lib.PipeCommander = cReboot
	_       lib.Progresser    = cReboot
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
	c.wait = c.Context.IsSet("wait")

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

	var wg sync.WaitGroup
	var once sync.Once

	ch := make(chan interface{})

	for item := range in {
		wg.Add(1)
		item := item
		go func() {
			defer wg.Done()
			id := item.(string)
			err := servers.Reboot(c.ServiceClient, id, c.opts).ExtractErr()
			if err != nil {
				out <- err
				return
			}

			switch c.wait {
			case true:
				once.Do(c.InitProgress)
				c.StartBar(&openstack.ProgressStatus{
					Name:      id,
					StartTime: time.Now(),
				})

				err := util.WaitFor(900, func() (bool, error) {
					var m map[string]map[string]interface{}
					err := servers.Get(c.ServiceClient, id).ExtractInto(&m)
					if err != nil {
						return false, err
					}

					switch m["server"]["status"].(string) {
					case "ACTIVE":
						c.CompleteBar(&openstack.ProgressStatus{
							Name: item.(string),
						})
						ch <- fmt.Sprintf("Rebooted server [%s]", id)
						return true, nil
					default:
						c.UpdateBar(&openstack.ProgressStatus{
							Name: item.(string),
						})
						return false, nil
					}
				})

				if err != nil {
					c.ErrorBar(&openstack.ProgressStatus{
						Name: item.(string),
						Err:  err,
					})
					ch <- err
				}
			default:
				out <- fmt.Sprintf("Rebooting server [%s]", id)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	msgs := make([]string, 0)

	for raw := range ch {
		switch msg := raw.(type) {
		case error:
			out <- msg
		case string:
			msgs = append(msgs, msg)
		}
	}

	for _, msg := range msgs {
		out <- msg
	}
}

func (c *commandReboot) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandReboot) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Rebooting"
	c.Progress.DoneMsg = "Rebooted"
	c.Progress.Start()
}
