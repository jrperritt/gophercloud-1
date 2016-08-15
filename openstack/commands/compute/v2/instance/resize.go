package instance

import (
	"fmt"
	"time"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type commandResize struct {
	openstack.CommandUtil
	InstanceV2Command
	opts servers.ResizeOptsBuilder
	*openstack.Progress
}

var (
	cResize                   = new(commandResize)
	_       lib.PipeCommander = cResize
	_       lib.Progresser    = cResize
	_       lib.Waiter        = cResize

	flagsResize = openstack.CommandFlags(cResize)
)

var resize = cli.Command{
	Name:         "resize",
	Usage:        util.Usage(commandPrefix, "resize", "[--id <serverID> | --name <serverName> | --stdin id] [--flavor-id | --flavor-name]"),
	Description:  "Resizes a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cResize) },
	Flags:        flagsResize,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsResize) },
}

func (c *commandResize) Flags() []cli.Flag {
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
		cli.BoolFlag{
			Name:  "wait",
			Usage: "[optional] If provided, will wait to return until the server has been resizeed.",
		},
	}
}

func (c *commandResize) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	c.Quiet = c.Context.IsSet("quiet")

	opts := new(servers.ResizeOpts)

	if c.Context.IsSet("flavor-id") {
		opts.FlavorRef = c.Context.String("flavor-id")
		c.opts = opts
		return nil
	}

	if c.Context.IsSet("flavor-name") {
		id, err := flavors.IDFromName(c.ServiceClient, c.Context.String("flavor-name"))
		if err != nil {
			return err
		}
		opts.FlavorRef = id
		c.opts = opts
		return nil
	}

	return fmt.Errorf("One and only one of flavor-name and flavor-id must be provided")
}

func (c *commandResize) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandResize) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandResize) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		id := item.(string)
		err := servers.Resize(c.ServiceClient, id, c.opts).ExtractErr()
		if err != nil {
			out <- err
			return
		}
		switch c.Wait {
		case true:
			out <- id
		default:
			out <- fmt.Sprintf("Resizing server [%s]", id)
		}
	}
}

func (c *commandResize) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandResize) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}

func (c *commandResize) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Resizing"
	c.Progress.DoneMsg = "Resized"
	c.ProgressChan = make(chan *openstack.ProgressStatus)
	go c.Progress.Listen(c.ProgressChan)
	if !c.Quiet {
		c.Progress.Start()
	}
}

func (c *commandResize) ShowProgress(raw interface{}, out chan interface{}) {
	id := (raw).(string)

	c.ProgressChan <- &openstack.ProgressStatus{
		Name:      id,
		StartTime: time.Now(),
		Type:      "start",
	}

	err := util.WaitFor(900, func() (bool, error) {
		_, err := servers.Get(c.ServiceClient, id).Extract()
		if err != nil {
			c.ProgressChan <- &openstack.ProgressStatus{
				Name: id,
				Type: "complete",
			}
			out <- fmt.Sprintf("Resized server [%s]", id)
			return true, nil
		}

		c.ProgressChan <- &openstack.ProgressStatus{
			Name: id,
			Type: "update",
		}
		return false, nil
	})

	if err != nil {
		c.ProgressChan <- &openstack.ProgressStatus{
			Name: id,
			Err:  err,
			Type: "error",
		}
		out <- err
	}
}
