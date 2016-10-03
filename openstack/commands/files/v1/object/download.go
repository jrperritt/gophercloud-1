package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/openstack/commands"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandDownload struct {
	ObjectV1Command
	commands.Progressable
	file string
}

var (
	cDownload = new(commandDownload)
	//	_         openstack.Progresser     = cDownload
	_ openstack.CustomWriterer = cDownload

	flagsDownload = openstack.CommandFlags(cDownload)
)

var download = cli.Command{
	Name:         "Download",
	Usage:        util.Usage(commandPrefix, "download", "--container <containerName> --name <objectName>"),
	Description:  "Downloads an object",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cDownload) },
	Flags:        flagsDownload,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsDownload) },
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

func (c *commandDownload) Execute(_ interface{}, out chan interface{}) {
	res := objects.Download(c.ServiceClient, c.container, c.name, nil)
	switch res.Err {
	case nil:
		out <- res.Body
	default:
		out <- res.Err
	}
}

func (c *commandDownload) InitProgress() {
	c.ProgressInfo = openstack.NewProgressInfo(1)
	c.Progressable.InitProgress()
}

func (c *commandDownload) ShowBar(raw interface{}) {
	orig := raw.(map[string]interface{})
	id := orig["id"].(string)

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
			s.Msg = r.(string)
			c.ProgressInfo.UpdateChan <- s
		}
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
