package instance

import (
	"fmt"
	"reflect"
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
	cd                   = new(commandDelete)
	_  lib.PipeCommander = cd
	_  lib.Progresser    = cd
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

	deletedServersChan := make(chan string)

	for item := range in {
		wg.Add(1)
		item := item
		go func() {
			defer wg.Done()
			id, ok := item.(string)
			if !ok {
				panic(fmt.Sprintf("expected string for id of server to delete but got %v [%v]", reflect.TypeOf(item), item))
			}
			err := servers.Delete(c.ServiceClient, id).ExtractErr()
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

				i := 0
				for i < 120 {
					_, err := servers.Get(c.ServiceClient, id).Extract()
					if err != nil {
						c.CompleteBar(&openstack.ProgressStatus{
							Name: id,
						})
						deletedServersChan <- fmt.Sprintf("Deleted server [%s]", id)
						break
					}
					time.Sleep(2 * time.Second)
					c.UpdateBar(&openstack.ProgressStatus{
						Name: id,
					})
					i++
				}
			default:
				out <- fmt.Sprintf("Deleting server [%s]", id)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(deletedServersChan)
	}()

	deletedServers := make([]string, 0)

	for deletedServer := range deletedServersChan {
		deletedServers = append(deletedServers, deletedServer)
	}

	for _, deletedServer := range deletedServers {
		out <- deletedServer
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
