package volumeattachment

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	VolumeAttachmentV2Command
	traits.Waitable
	traits.Pipeable
	traits.DataResp
	traits.Tableable
}

var (
	cList                              = new(CommandList)
	_         interfaces.Waiter        = cList
	_         interfaces.PipeCommander = cList
	_         interfaces.Tabler        = cList
	flagsList                          = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists all volume attachments associated with a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "server-id",
			Usage: "[optional; required if `server-name` or `stdin` isn't provided] The server ID of the attachment.",
		},
		cli.StringFlag{
			Name:  "server-name",
			Usage: "[optional; required if `server-id` or `stdin` isn't provided] The server name of the attachment.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `server-id` or `server-name` isn't provided] The field being piped into STDIN. Valid values are: server-id",
		},
	}
}

// DefaultTableFields returns default fields for tabular output.
// Partially satisfies interfaces.Tabler interface
func (c *CommandList) DefaultTableFields() []string {
	return []string{"id", "device", "volume_id", "server_id"}
}

func (c *CommandList) HandleSingle() (interface{}, error) {
	return serverIDorName(c.Context(), c.ServiceClient)
}

func (c *CommandList) Execute(item interface{}, out chan interface{}) {
	p, err := volumeattach.List(c.ServiceClient, item.(string)).AllPages()
	if err != nil {
		out <- err
		return
	}

	var m map[string][]map[string]interface{}
	err = (p.(volumeattach.VolumeAttachmentPage)).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}
	out <- m["volumeAttachments"]
}

func (c *CommandList) PipeFieldOptions() []string {
	return []string{"server-id"}
}
