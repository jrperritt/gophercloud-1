package volumeattachment

import (
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"gopkg.in/urfave/cli.v1"
)

type CommandGet struct {
	VolumeAttachmentV2Command

	serverID     string
	attachmentID string
}

var (
	cGet     = new(CommandGet)
	flagsGet = openstack.CommandFlags(cGet)
)

var get = cli.Command{
	Name:         "get",
	Usage:        util.Usage(commandPrefix, "get", "--id <ID> [--server-id <ID> | --server-name <NAME>]"),
	Description:  "Gets a volume attachment on a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cGet) },
	Flags:        flagsGet,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsGet) },
}

func (c *CommandGet) Flags() []cli.Flag {
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

func (c *CommandGet) HandleFlags() error {
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

func (c *CommandGet) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	err := volumeattach.Get(c.ServiceClient(), c.serverID, c.attachmentID).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}
	out <- m["volumeAttachment"]
}
