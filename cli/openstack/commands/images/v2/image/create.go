package image

import (
	"strings"

	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	ImageV2Command
	opts images.CreateOptsBuilder
}

var (
	cCreate = new(CommandCreate)

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(commandPrefix, "create", "--name NAME"),
	Description:  "Creates an image",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the image",
		},
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional] The ID of the image",
		},
		cli.StringFlag{
			Name:  "visibility",
			Usage: "[optional] Who can see the image. Options: public, private, shared, community",
		},
		cli.StringFlag{
			Name:  "tags",
			Usage: "[optional] Comma-separated values associated with the image",
		},
		cli.StringFlag{
			Name:  "container-format",
			Usage: "[optional] The container format. Options: ami, ari, aki, bare, ovf",
		},
		cli.StringFlag{
			Name:  "disk-format",
			Usage: "[optional] The disk format. Options: ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, iso",
		},
		cli.IntFlag{
			Name:  "min-disk",
			Usage: "[optional] Amount of space required to boot the image (in GB)",
		},
		cli.IntFlag{
			Name:  "min-ram",
			Usage: "[optional] Amount of RAM required to boot the image (in MB)",
		},
		cli.BoolFlag{
			Name:  "protected",
			Usage: "[optional] If provided, the image is not deletable",
		},
		cli.StringFlag{
			Name:  "properties",
			Usage: "[optional] Comma-separated key-value pairs for the image. Example: key1=val1,key2=val2",
		},
	}
}

func (c *CommandCreate) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"name"})
	if err != nil {
		return err
	}

	opts := new(images.CreateOpts)
	opts.Name = c.Context().String("name")
	opts.ID = c.Context().String("id")
	opts.ContainerFormat = c.Context().String("container-format")
	opts.DiskFormat = c.Context().String("disk-format")
	opts.Tags = strings.Split(c.Context().String("tags"), ",")
	opts.MinDisk = c.Context().Int("min-disk")
	opts.MinRAM = c.Context().Int("min-ram")

	if c.Context().IsSet("visibility") {
		v := images.ImageVisibility(c.Context().String("visibility"))
		opts.Visibility = &v
	}

	if c.Context().IsSet("protected") {
		t := true
		opts.Protected = &t
	}

	if c.Context().IsSet("properties") {
		metadata, err := c.ValidateKVFlag("properties")
		if err != nil {
			return err
		}
		opts.Properties = metadata
	}

	c.opts = opts

	return nil
}

func (c *CommandCreate) Execute(item interface{}, out chan interface{}) {
	var m map[string]interface{}
	err := images.Create(c.ServiceClient(), c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m
	default:
		out <- err
	}
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}
