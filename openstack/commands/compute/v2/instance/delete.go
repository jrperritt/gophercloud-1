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

type commandDelete struct {
	openstack.CommandUtil
	InstanceV2Command
	*openstack.Progress
}

var (
	cDelete                   = new(commandDelete)
	_       lib.PipeCommander = cDelete
	_       lib.Progresser    = cDelete
	_       lib.Waiter        = cDelete
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
	c.Wait = c.Context.IsSet("wait")
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
	for item := range in {
		id := item.(string)
		err := servers.Delete(c.ServiceClient, id).ExtractErr()
		if err != nil {
			out <- err
			return
		}
		switch c.Wait {
		case true:
			out <- id
		default:
			out <- fmt.Sprintf("Deleting server [%s]", id)
		}
	}
}

func (c *commandDelete) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandDelete) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}

func (c *commandDelete) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Deleting"
	c.Progress.DoneMsg = "Deleted"
	c.Progress.Start()
}

func (c *commandDelete) ShowProgress(in, out chan interface{}) {
	for raw := range in {
		id := (raw).(string)

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
				out <- fmt.Sprintf("Deleted server [%s]", id)
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
}
