package instance

import (
	"fmt"
	"time"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
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

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDelete) },
}

func (c *commandDelete) Flags() []cli.Flag {
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
	}
}

func (c *commandDelete) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	c.Quiet = c.Context.IsSet("quiet")
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
	c.ProgressChan = make(chan *openstack.ProgressStatus)
	go c.Progress.Listen(c.ProgressChan)
	if !c.Quiet {
		c.Progress.Start()
	}
}

func (c *commandDelete) ShowProgress(raw interface{}, out chan interface{}) {
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
			out <- fmt.Sprintf("Deleted server [%s]", id)
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
