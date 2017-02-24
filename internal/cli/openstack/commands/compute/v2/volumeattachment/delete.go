package volumeattachment

import (
	"fmt"

	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"gopkg.in/urfave/cli.v1"
)

type CommandDelete struct {
	VolumeAttachmentV2Command

	serverID     string
	attachmentID string
}

var (
	cDelete     = new(CommandDelete)
	flagsDelete = openstack.CommandFlags(cDelete)
)

var remove = cli.Command{
	Name:         "delete",
	Usage:        util.Usage(commandPrefix, "delete", "--id <ID> [--server-id <ID> | --server-name <NAME>]"),
	Description:  "Deletes a volume attachment from a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDelete) },
	Flags:        flagsDelete,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDelete) },
}

func (c *CommandDelete) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[required] The ID of the attachment",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The server ID to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The server name to which attach the volume.",
		},
	}
}

func (c *CommandDelete) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"id"})
	if err != nil {
		return err
	}
	c.attachmentID = c.Context().String("id")

	c.serverID, err = serverIDorName(c.Context(), c.ServiceClient())
	if err != nil {
		return err
	}

	return nil
}

func (c *CommandDelete) Execute(item interface{}, out chan interface{}) {
	err := volumeattach.Delete(c.ServiceClient(), c.serverID, c.attachmentID).ExtractErr()
	if err != nil {
		out <- err
		return
	}
	out <- fmt.Sprintf("Successfully deleted volume attachment [%s] from server [%s]", c.attachmentID, c.serverID)
}
