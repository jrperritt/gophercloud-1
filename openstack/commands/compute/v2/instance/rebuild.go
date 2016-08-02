package instance

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type commandRebuild struct {
	openstack.CommandUtil
	InstanceV2Command
	wait bool
	opts servers.RebuildOptsBuilder
	*openstack.Progress
}

var (
	cRebuild                   = new(commandRebuild)
	_        lib.PipeCommander = cRebuild
	_        lib.Progresser    = cRebuild
)

var rebuild = cli.Command{
	Name:         "rebuild",
	Usage:        util.Usage(commandPrefix, "rebuild", "[--id <serverID> | --name <serverName> | --stdin id] [--image-id | --image-name]"),
	Description:  "Rebuilds a server",
	Action:       actionRebuild,
	Flags:        openstack.CommandFlags(flagsRebuild, []string{""}),
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsRebuild) },
}

func actionRebuild(ctx *cli.Context) {
	c := new(commandRebuild)
	c.Context = ctx
	lib.Run(ctx, c)
}

var flagsRebuild = []cli.Flag{
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
		Usage: "[optional] The server's admin password",
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
	cli.BoolFlag{
		Name:  "wait",
		Usage: "[optional] If provided, will wait to return until the server has been rebuilt.",
	},
}

func (c *commandRebuild) HandleFlags() error {
	c.wait = c.Context.IsSet("wait")

	opts := &servers.RebuildOpts{
		ImageID:       c.Context.String("image-id"),
		ImageName:     c.Context.String("image-name"),
		AdminPass:     c.Context.String("admin-pass"),
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

	c.opts = opts

	return nil
}

func (c *commandRebuild) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *commandRebuild) HandleSingle() (interface{}, error) {
	return c.IDOrName(servers.IDFromName)
}

func (c *commandRebuild) Execute(in, out chan interface{}) {
	defer close(out)

	var wg sync.WaitGroup
	var once sync.Once

	ch := make(chan interface{})

	for item := range in {
		wg.Add(1)
		item := item
		go func() {
			defer wg.Done()
			id := item.(string)
			m := make(map[string]map[string]interface{})
			err := servers.Rebuild(c.ServiceClient, id, c.opts).ExtractInto(&m)
			if err != nil {
				switch c.wait {
				case true:
					ch <- err
				case false:
					out <- err
				}
				return
			}

			switch c.wait {
			case true:
				once.Do(c.InitProgress)
				c.StartBar(&openstack.ProgressStatus{
					Name:      id,
					StartTime: time.Now(),
				})

				err := util.WaitFor(900, func() (bool, error) {
					var m map[string]map[string]interface{}
					err := servers.Get(c.ServiceClient, id).ExtractInto(&m)
					if err != nil {
						return false, err
					}

					switch m["server"]["status"].(string) {
					case "ACTIVE":
						c.CompleteBar(&openstack.ProgressStatus{
							Name: item.(string),
						})
						ch <- m
						return true, nil
					default:
						c.UpdateBar(&openstack.ProgressStatus{
							Name: item.(string),
						})
						return false, nil
					}
				})

				if err != nil {
					c.ErrorBar(&openstack.ProgressStatus{
						Name: item.(string),
						Err:  err,
					})
					ch <- err
				}
			default:
				out <- m
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	msgs := make([]map[string]interface{}, 0)

	for raw := range ch {
		switch msg := raw.(type) {
		case error:
			out <- msg
		case map[string]interface{}:
			msgs = append(msgs, msg)
		}
	}

	for _, msg := range msgs {
		out <- msg
	}
}

func (c *commandRebuild) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *commandRebuild) InitProgress() {
	c.Progress = openstack.NewProgress(2)
	c.Progress.RunningMsg = "Rebuilding"
	c.Progress.DoneMsg = "Rebuilt"
	c.Progress.Start()
}
