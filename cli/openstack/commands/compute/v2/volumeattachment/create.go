package volumeattachment

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	VolumeAttachmentV2Command
	traits.Waitable
	traits.Pipeable
	traits.DataResp
	opts     volumeattach.CreateOptsBuilder
	serverID string
}

var (
	cCreate                              = new(CommandCreate)
	_           interfaces.Waiter        = cCreate
	_           interfaces.PipeCommander = cCreate
	flagsCreate                          = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "[--server-id <ID> | --server-name <NAME>] [--volume-id <ID> | --volume-name <NAME> | --stdin volume-id]"),
	Description:  "Creates a volume attachment on a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "volume-id",
			Usage: "[optional; required if `stdin` or volume-name isn't provided] The ID of the volume to attach.",
		},
		cli.StringFlag{
			Name:  "volume-name",
			Usage: "[optional; required if `stdin` or `volume-id` isn't provided] The name of the volume to attach.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `volume-id` or `volume-name` isn't provided] The field being piped into STDIN. Valid values are: volume-id",
		},
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` isn't provided] The server ID to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` isn't provided] The server name to which attach the volume.",
		},
		cli.StringFlag{
			Name:  "device",
			Usage: "[optional] The name of the device to which the volume will attach. Default is 'auto'.",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	serverID, err := serverIDorName(c.Context, c.ServiceClient)
	if err != nil {
		return err
	}
	c.serverID = serverID

	c.opts = &volumeattach.CreateOpts{
		Device: c.Context.String("device"),
	}
	return nil
}

func (c *CommandCreate) HandleSingle() (interface{}, error) {
	return volumeIDorName(c.Context, c.ServiceClient)
}

func (c *CommandCreate) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	opts := *c.opts.(*volumeattach.CreateOpts)
	opts.VolumeID = item.(string)
	err := volumeattach.Create(c.ServiceClient, c.serverID, opts).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}
	out <- m["volumeAttachment"]
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"volume-id"}
}
