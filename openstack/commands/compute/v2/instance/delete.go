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

type commandDelete struct {
	openstack.CommandUtil
	InstanceV2Command
	wait bool
	*openstack.Progress
}

var (
	cDelete                   = new(commandDelete)
	_       lib.PipeCommander = cDelete
	_       lib.Progresser    = cDelete
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes a server",
	Action:       actionDelete,
	Flags:        openstack.CommandFlags(flagsDelete, []string{}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDelete) },
}

func actionDelete(ctx *cli.Context) {
	c := new(commandDelete)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsDelete = []cli.Flag{
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
		Name:  "wait",
		Usage: "[optional] If provided, will wait to return until the server has been deleted.",
	},
}

func (c *commandDelete) HandleFlags() error {
	c.wait = c.Context.IsSet("wait")
	return nil
}

func (c *commandDelete) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandDelete) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandDelete) Execute(in, out chan interface{}) {
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
			err := servers.Delete(c.ServiceClient, id).ExtractErr()
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
					_, err := servers.Get(c.ServiceClient, id).Extract()
					if err != nil {
						c.CompleteBar(&openstack.ProgressStatus{
							Name: id,
						})
						ch <- fmt.Sprintf("Deleted server [%s]", id)
						return true, nil
					}

					c.UpdateBar(&openstack.ProgressStatus{
						Name: id,
					})
					return false, nil
				})

				if err != nil {
					c.ErrorBar(&openstack.ProgressStatus{
						Name: item.(string),
						Err:  err,
					})
					ch <- err
				}
			default:
				out <- fmt.Sprintf("Deleting server [%s]", id)
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

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandDelete) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Deleting"
	c.Progress.DoneMsg = "Deleted"
	c.Progress.Start()
}
