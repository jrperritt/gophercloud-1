package instance

import (
	"fmt"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type commandResize struct {
	openstack.CommandUtil
	InstanceV2Command
	wait bool
	opts servers.ResizeOptsBuilder
	*openstack.Progress
}

var (
	cResize                   = new(commandResize)
	_       lib.PipeCommander = cResize
	_       lib.Progresser    = cResize
)

var resize = cli.Command{
	Name:         "resize",
	Usage:        util.Usage(commandPrefix, "resize", "[--id <serverID> | --name <serverName> | --stdin id] [--flavor-id | --flavor-name]"),
	Description:  "Resizes a server",
	Action:       actionResize,
	Flags:        openstack.CommandFlags(flagsResize, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsResize) },
}

func actionResize(ctx *cli.Context) {
	c := new(commandResize)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsResize = []cli.Flag{
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

func (c *commandResize) HandleFlags() error {
	c.wait = c.Context.IsSet("wait")

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

	var wg sync.WaitGroup
	var once sync.Once

	ch := make(chan interface{})

	for item := range in {
		wg.Add(1)
		item := item
		go func() {
			defer wg.Done()
			id := item.(string)
			err := servers.Resize(c.ServiceClient, id, c.opts).ExtractErr()
			if err != nil {
				switch c.wait {
				case true:
					ch <- err
				case false:
					out <- err
				}
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
						ch <- fmt.Sprintf("Resized server [%s]", id)
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
				out <- fmt.Sprintf("Resizing server [%s]", id)
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

func (c *commandResize) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandResize) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Resizing"
	c.Progress.DoneMsg = "Resized"
	c.Progress.Start()
}
