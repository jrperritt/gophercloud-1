package instance

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type commandRebuild struct {
	openstack.CommandUtil
	InstanceV2Command
	opts servers.RebuildOptsBuilder
	*openstack.Progress
}

var (
	cRebuild                   = new(commandRebuild)
	_        lib.PipeCommander = cRebuild
	_        lib.Progresser    = cRebuild
	_        lib.Waiter        = cRebuild

	flagsRebuild = openstack.CommandFlags(cRebuild)
)

var rebuild = cli.Command{
	Name:         "rebuild",
	Usage:        util.Usage(commandPrefix, "rebuild", "[--id <serverID> | --name <serverName> | --stdin id] [--image-id | --image-name]"),
	Description:  "Rebuilds a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cRebuild) },
	Flags:        flagsRebuild,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsRebuild) },
}

func (c *commandRebuild) Flags() []cli.Flag {
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

func (c *commandRebuild) Fields() []string {
	return []string{""}
}

func (c *commandRebuild) HandleFlags() error {
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

func (c *commandRebuild) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandRebuild) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}

func (c *commandRebuild) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandRebuild) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		id := item.(string)
		m := make(map[string]map[string]interface{})
		err := servers.Rebuild(c.ServiceClient, id, c.opts).ExtractInto(&m)
		if err != nil {
			out <- err
			return
		}
		switch c.Wait {
		case true:
			out <- id
		default:
			out <- fmt.Sprintf("Deleting server [%s]", id)
		}
	}
}

func (c *commandRebuild) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandRebuild) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Rebuilding"
	c.Progress.DoneMsg = "Rebuilt"
	c.ProgressChan = make(chan *openstack.ProgressStatus)
	go c.Progress.Listen(c.ProgressChan)
	if !c.Quiet {
		c.Progress.Start()
	}
}

func (c *commandRebuild) ShowProgress(raw interface{}, out chan interface{}) {
	orig := raw.(map[string]interface{})
	id := orig["id"].(string)

	c.ProgressChan <- &openstack.ProgressStatus{
		Name:      id,
		StartTime: time.Now(),
		Type:      "start",
	}

	err := util.WaitFor(900, func() (bool, error) {
		var m map[string]map[string]interface{}
		err := servers.Get(c.ServiceClient, id).ExtractInto(&m)
		if err != nil {
			return false, err
		}

		switch m["server"]["status"].(string) {
		case "ACTIVE":
			c.ProgressChan <- &openstack.ProgressStatus{
				Name: id,
				Type: "complete",
			}
			out <- m
			return true, nil
		default:
			c.ProgressChan <- &openstack.ProgressStatus{
				Name: id,
				Type: "update",
			}
			return false, nil
		}
	})

	if err != nil {
		c.ProgressChan <- &openstack.ProgressStatus{
			Name: id,
			Err:  err,
			Type: "error",
		}
		out <- err
	}
}
