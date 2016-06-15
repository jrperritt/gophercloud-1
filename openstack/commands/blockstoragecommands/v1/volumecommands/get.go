package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

type commandGet struct {
	openstack.Command
	id string
}

func (c *commandGet) Name() string {
	return "get"
}

func (c *commandGet) Usage() string {
	return util.Usage(commandPrefix, "get", "[--id <volumeID> | --name <volumeName> | --stdin id]")
}

func (c *commandGet) Description() string {
	return "Gets a volume"
}

var flagsGet = []cli.Flag{
	cli.StringFlag{
		Name:  "id",
		Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the volume.",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the volume.",
	},
	cli.StringFlag{
		Name:  "stdin",
		Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
	},
}

func (c *commandGet) HandleFlags() error {
	return nil
}

func (c *commandGet) HandlePipe(item string) error {
	c.id = item
	return nil
}

func (c *commandGet) HandleSingle() error {
	id, err := c.IDOrName(volumes.IDFromName)
	if err != nil {
		return err
	}
	c.id = id
	return nil
}

func (c *commandGet) Execute(_ lib.Resourcer) (r lib.Resulter) {
	var m map[string]interface{}
	err := volumes.Get(c.ServiceClient(), c.id).ExtractInto(m)
	if err != nil {
		r.SetError(err)
		return
	}
	r.SetValue(m)
	return
}

func (c *commandGet) StdinField() string {
	return "id"
}

/*
func (c *commandGet) PreCSV() error {
	resource.FlattenMap("Metadata")
	resource.FlattenMap("Attachments")
	return nil
}

func (c *commandGet) PreTable() error {
	return command.PreCSV(resource)
}
*/
