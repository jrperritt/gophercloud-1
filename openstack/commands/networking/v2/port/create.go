package port

import (
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	PortV2Command
	commands.Waitable
	commands.DataResp
	opts ports.CreateOptsBuilder
}

var (
	cCreate                         = new(CommandCreate)
	_       openstack.PipeCommander = cCreate

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(CommandPrefix, "create", "--network-id <ID>"),
	Description:  "Creates a port",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "network-id",
			Usage: "[required] The network ID of the port.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional] The name of the port.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional] The field being piped into STDIN. Valid values are: name",
		},
		cli.BoolFlag{
			Name:  "up",
			Usage: "[optional] If provided, the port will be up upon creation.",
		},
		cli.StringFlag{
			Name:  "tenant-id",
			Usage: "[optional] The ID of the tenant that will own the port.",
		},
		cli.StringFlag{
			Name:  "device-id",
			Usage: "[optional] The device ID to associate with the port.",
		},
	}
}

func (c *CommandCreate) Fields() []string {
	return []string{""}
}

func (c *CommandCreate) HandleFlags() error {
	opts := &ports.CreateOpts{
		NetworkID: c.Context.String("network-id"),
		TenantID:  c.Context.String("tenant-id"),
		DeviceID:  c.Context.String("device-id"),
	}

	if c.Context.IsSet("up") {
		opts.AdminStateUp = gophercloud.Enabled
	}

	c.opts = opts
	c.Wait = c.Context.IsSet("wait")

	return nil
}

func (c *CommandCreate) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandCreate) HandleSingle() (interface{}, error) {
	return c.Context.String("name"), c.CheckFlagsSet([]string{"network-id"})
}

func (c *CommandCreate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	opts := *c.opts.(*ports.CreateOpts)
	opts.Name = item.(string)
	err := ports.Create(c.ServiceClient, opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["port"]
	default:
		out <- err
	}
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}
