package volumecommands

import (
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
)

type commandGet struct {
	openstack.Command
	id string
}

var get = func() cli.Command {
	c := new(commandGet)
	c.SetFlags(flagsGet)
	c.SetDefaultFields()
	return openstack.NewCommand(c)
}

func (c commandGet) Name() string {
	return "get"
}

func (c commandGet) Usage() string {
	return util.Usage(commandPrefix, "get", "[--id <volumeID> | --name <volumeName> | --stdin id]")
}

func (c commandGet) Description() string {
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

func (command *commandGet) HandleFlags() error {
	resource.Params = &paramsGet{}
	return nil
}

func (command *commandGet) HandlePipe(item string) error {
	resource.Params.(*paramsGet).volumeID = item
	return nil
}

func (command *commandGet) HandleSingle() error {
	volumeID, err := command.Ctx.IDOrName(osVolumes.IDFromName)
	if err != nil {
		return err
	}
	resource.Params.(*paramsGet).volumeID = volumeID
	return nil
}

func (command *commandGet) Execute() {
	volumeID := resource.Params.(*paramsGet).volumeID
	volume, err := osVolumes.Get(command.Ctx.ServiceClient, volumeID).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = volumeSingle(volume)
}

func (command *commandGet) StdinField() string {
	return "id"
}

func (command *commandGet) PreCSV() error {
	resource.FlattenMap("Metadata")
	resource.FlattenMap("Attachments")
	return nil
}

func (command *commandGet) PreTable() error {
	return command.PreCSV(resource)
}
