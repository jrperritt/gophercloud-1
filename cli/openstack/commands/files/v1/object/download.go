package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandDownload struct {
	ObjectV1Command
	traits.BytesProgressable
	file string
}

var (
	cDownload                           = new(commandDownload)
	_         interfaces.Progresser     = cDownload
	_         interfaces.CustomWriterer = cDownload

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
	c.container = c.Context().String("container")
	c.name = c.Context().String("object")

	if c.Context().IsSet("wait") {
		err := c.CheckFlagsSet([]string{"file"})
		if err != nil {
			return err
		}
		c.file = c.Context().String("file")
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

type filecontent struct {
	writer      io.Writer
	bytessentch chan (interface{})
}

func (b *filecontent) Write(p []byte) (n int, err error) {
	n, err = b.writer.Write(p)
	if err != nil {
		return
	}
	b.bytessentch <- n
	return
}

func (c *commandDownload) CustomWriter() (io.Writer, error) {
	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("\nA file named %s already exists. Overwrite? (y/n): ", c.file)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			switch strings.ToLower(choice) {
			case "y", "yes":
				goto done
			case "n", "no":
				return nil, err
			default:
				continue
			}
		}
	}

done:
	fc := new(filecontent)
	fc.writer = f
	fc.bytessentch = c.ProgUpdateChIn()
	return fc, nil
}
