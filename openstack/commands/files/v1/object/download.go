package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandDownload struct {
	openstack.CommandUtil
	ObjectV1Command
	*openstack.Progress
	container string
	name      string
	file      string
}

var (
	cDownload                    = new(commandDownload)
	_         lib.Waiter         = cDownload
	_         lib.Progresser     = cDownload
	_         lib.CustomWriterer = cDownload

	flagsDownload = openstack.CommandFlags(cDownload)
)

var download = cli.Command{
	Name:         "Download",
	Usage:        util.Usage(commandPrefix, "download", "--container <containerName> --name <objectName>"),
	Description:  "Downloads an object",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDownload) },
	Flags:        flagsDownload,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsDownload) },
}

func (c *commandDownload) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the object.",
		},
		cli.StringFlag{
			Name: "file",
			Usage: "[optional; required if `wait` is provided] The name for the file to which \n" +
				"\t the object should be saved",
		},
	}
}

func (c *commandDownload) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return err
	}
	c.container = c.Context.String("container")
	c.name = c.Context.String("object")

	if c.Context.IsSet("wait") {
		err := c.CheckFlagsSet([]string{"file"})
		if err != nil {
			return err
		}
		c.file = c.Context.String("file")
	}

	return nil
}

func (c *commandDownload) Execute(_, out chan interface{}) {
	defer close(out)
	res := objects.Download(c.ServiceClient, c.container, c.name, nil)
	switch res.Err {
	case nil:
		out <- res.Body
	default:
		out <- res.Err
	}
}

func (c *commandDownload) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}

func (c *commandDownload) InitProgress() {
	c.Progress = openstack.NewProgress(1)
	c.ProgressChan = make(chan *openstack.ProgressStatus)
	go c.Progress.Listen(c.ProgressChan)
	if !c.Quiet {
		c.Progress.Start()
	}
}

func (c *commandDownload) ShowProgress(raw interface{}, out chan interface{}) {
	orig := raw.(map[string]interface{})
	id := orig["id"].(string)

	c.ProgressChan <- &openstack.ProgressStatus{
		Name:      id,
		TotalSize: 100,
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
			m["server"]["adminPass"] = orig["adminPass"].(string)
			out <- m["server"]
			return true, nil
		default:
			c.ProgressChan <- &openstack.ProgressStatus{
				Name:      id,
				Increment: int(m["server"]["progress"].(float64)),
				Type:      "update",
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

func (c *commandDownload) CustomWriter() io.Writer {
	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("\nA file named %s already exists. Overwrite? (y/n): ", c.file)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			switch strings.ToLower(choice) {
			case "y", "yes":
				return f
			case "n", "no":
				os.Exit(0)
			default:
				continue
			}
		}
	}

	return f
}
