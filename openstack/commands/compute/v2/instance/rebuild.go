package instance

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type CommandRebuild struct {
	ServerV2Command
	commands.ProgressCommand
	opts servers.RebuildOptsBuilder
}

var (
	cRebuild                         = new(CommandRebuild)
	_        openstack.PipeCommander = cRebuild
	_        openstack.Progresser    = cRebuild

	flagsRebuild = openstack.CommandFlags(cRebuild)
)

var rebuild = cli.Command{
	Name:         "rebuild",
	Usage:        util.Usage(commandPrefix, "rebuild", "[--id <serverID> | --name <serverName> | --stdin id] [--image-id | --image-name]"),
	Description:  "Rebuilds a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cRebuild) },
	Flags:        flagsRebuild,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsRebuild) },
}

func (c *CommandRebuild) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `stdin` or `name` isn't provided] The ID of the server",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` or `id` isn't provided] The name of the server",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `id` or `name` isn't provided] The field being piped into STDIN. Valid values are: id",
		},
		cli.StringFlag{
			Name:  "image-id",
			Usage: "[optional; required if `image-name`is not provided] The image ID from which to rebuild the server.",
		},
		cli.StringFlag{
			Name:  "image-name",
			Usage: "[optional; required if `image-id` is not provided] The name of the image from which to rebuild the server.",
		},
		cli.StringFlag{
			Name:  "admin-pass",
			Usage: "[optional; required if `generate-pass` is not provided] The new server's admin password",
		},
		cli.BoolFlag{
			Name:  "generate-pass",
			Usage: "[optional; required if `admin-pass` is not provided] If provided, a password will be generated for the new server.",
		},
		cli.StringFlag{
			Name:  "rename",
			Usage: "[optional] The name for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "ipv4",
			Usage: "[optional] The IPv4 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "ipv6",
			Usage: "[optional] The IPv6 address for the rebuilt server.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string a key=value pairs.",
		},
		cli.StringFlag{
			Name: "personality",
			Usage: "[optional] A comma-separated list of key=value pairs. The key is the\n" +
				"\tdestination to inject the file on the created server; the value is the its local location.\n" +
				"\tExample: --personality \"C:\\cloud-automation\\bootstrap.cmd=open_hatch.cmd\"",
		},
	}
}

func (c *CommandRebuild) Fields() []string {
	return []string{""}
}

func (c *CommandRebuild) HandleFlags() error {
	c.Wait = c.Context.IsSet("wait")
	c.Quiet = c.Context.IsSet("quiet")

	opts := &servers.RebuildOpts{
		ImageID:       c.Context.String("image-id"),
		ImageName:     c.Context.String("image-name"),
		Name:          c.Context.String("rename"),
		AccessIPv4:    c.Context.String("ipv4"),
		AccessIPv6:    c.Context.String("ipv6"),
		ServiceClient: c.ServiceClient,
	}

	if c.Context.IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	if c.Context.IsSet("personality") {

		filesToInjectMap, err := c.ValidateKVFlag("personality")
		if err != nil {
			return err
		}

		if len(filesToInjectMap) > 5 {
			return fmt.Errorf("A maximum of 5 files may be provided for the `personality` flag")
		}

		filesToInject := make(servers.Personality, 0)
		for destinationPath, localPath := range filesToInjectMap {
			localAbsFilePath, err := filepath.Abs(localPath)
			if err != nil {
				return err
			}

			fileData, err := ioutil.ReadFile(localAbsFilePath)
			if err != nil {
				return err
			}

			if len(fileData)+len(destinationPath) > 1000 {
				return fmt.Errorf("The maximum length of a file-path-and-content pair for `personality` is 1000 bytes."+
					" Current pair size: path (%s): %d, content: %d", len(destinationPath), len(fileData))
			}

			filesToInject = append(filesToInject, &servers.File{
				Path:     destinationPath,
				Contents: fileData,
			})
		}
		opts.Personality = filesToInject
	}

	switch c.Context.IsSet("generate-pass") {
	case true:
		switch c.Context.IsSet("admin-pass") {
		case true:
			return fmt.Errorf("Only one of `generate-pass` and `admin-pass` may be provided")
		case false:
		}
	case false:
		switch c.Context.IsSet("admin-pass") {
		case true:
			opts.AdminPass = c.Context.String("admin-pass")
		case false:
			return fmt.Errorf("One of `generate-pass` and `admin-pass` must be provided")
		}
	}

	c.opts = opts

	return nil
}

func (c *CommandRebuild) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *CommandRebuild) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *CommandRebuild) Execute(item interface{}, out chan interface{}) {
	id := item.(string)
	m := make(map[string]map[string]interface{})
	err := servers.Rebuild(c.ServiceClient, id, c.opts).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}
	switch c.Wait || !c.Quiet {
	case true:
		out <- m
	default:
		out <- fmt.Sprintf("Rebuilding server [%s]", id)
	}
}

func (c *CommandRebuild) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *CommandRebuild) WaitFor(raw interface{}) {
	orig := raw.(map[string]interface{})
	id := orig["id"].(string)

	err := util.WaitFor(900, func() (bool, error) {
		var m map[string]map[string]interface{}
		err := servers.Get(c.ServiceClient, id).ExtractInto(&m)
		if err != nil {
			return false, err
		}

		switch m["server"]["status"].(string) {
		case "ACTIVE":
			m["server"]["adminPass"] = orig["adminPass"].(string)
			openstack.GC.DoneChan <- m["server"]
			return true, nil
		default:
			if !c.Quiet {
				openstack.GC.UpdateChan <- m["server"]["progress"].(float64)
			}
			return false, nil
		}
	})

	if err != nil {
		openstack.GC.DoneChan <- err
	}
}

func (c *CommandRebuild) InitProgress() {
	c.ProgressInfo = openstack.NewProgressInfo(2)
	c.ProgressInfo.RunningMsg = "Rebuilding"
	c.ProgressInfo.DoneMsg = "Rebuilt"
	c.ProgressCommand.InitProgress()
}

func (c *CommandRebuild) BarID(raw interface{}) string {
	orig := raw.(map[string]interface{})
	return orig["id"].(string)
}

func (c *CommandRebuild) ShowBar(id string) {
	s := new(openstack.ProgressStatusStart)
	s.Name = id
	c.StartChan <- s

	for {
		select {
		case r := <-openstack.GC.DoneChan:
			s := new(openstack.ProgressStatusComplete)
			s.Name = id
			c.ProgressInfo.CompleteChan <- s
			openstack.GC.ProgressDoneChan <- r
			return
		case r := <-openstack.GC.UpdateChan:
			s := new(openstack.ProgressStatusUpdate)
			s.Name = id
			s.Increment = int(r.(float64))
			c.ProgressInfo.UpdateChan <- s
		}
	}
}
