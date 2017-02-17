package image

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/pagination"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	ImageV2Command
	//traits.Waitable
	//traits.TextProgressable
}

var (
	cDelete                          = new(CommandDelete)
	_       interfaces.PipeCommander = cDelete
	//_       interfaces.Progresser    = cDelete

	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "[--id <serverID> | --name <serverName> | --stdin id]"),
	Description:  "Deletes an image",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the image.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the image.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
	}
}

func (c *CommandDelete) HandleFlags() error {
	return nil
}

func (c *CommandDelete) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandDelete) HandleSingle() (interface{}, error) {
	return c.IDOrName(IDFromName)
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	err := images.Delete(c.ServiceClient(), id).ExtractErr()
	if err != nil {
		out <- err
		return
	}
	out <- fmt.Sprintf("Deleting image [%s]", id)
}

func (c *CommandDelete) PipeFieldOptions() []string {
	return []string{"id"}
}

/*
func (c *CommandDelete) WaitFor(raw interface{}, out chan<- interface{}) {
	id := raw.(string)

	err := util.WaitFor(900, func() (bool, error) {
		_, err := servers.Get(c.ServiceClient(), id).Extract()
		if err != nil {
			out <- fmt.Sprintf("Deleted server [%s]", id)
			return true, nil
		}
		//c.ProgUpdateChIn() <- c.RunningMsg()
		return false, nil
	})

	if err != nil {
		out <- err
	}
}
*/

// IDFromName is a convienience function that returns an image's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	err := images.List(client, images.ListOpts{Name: name}).EachPage(func(page pagination.Page) (bool, error) {
		images, err := images.ExtractImages(page)
		if err != nil {
			return false, err
		}
		count = len(images)
		if len(images) > 0 {
			id = images[0].ID
		}
		return false, nil
	})
	if err != nil {
		return "", err
	}

	switch count {
	case 0:
		err := &gophercloud.ErrResourceNotFound{}
		err.ResourceType = "image"
		err.Name = name
		return "", err
	case 1:
		return id, nil
	default:
		err := &gophercloud.ErrMultipleResourcesFound{}
		err.ResourceType = "image"
		err.Name = name
		err.Count = count
		return "", err
	}
}
