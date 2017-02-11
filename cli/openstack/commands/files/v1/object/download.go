package object

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud/cli/lib"
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
	Name:         "download",
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
			Usage: "[optional; required if `quiet` is not provided] The name for the file to which \n" +
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
	c.name = c.Context().String("name")

	if c.ShouldProgress() {
		err := c.CheckFlagsSet([]string{"file"})
		if err != nil {
			return err
		}
		c.file = c.Context().String("file")
	}

	return nil
}

func (c *commandDownload) Execute(_ interface{}, out chan interface{}) {
	w, err := c.CustomWriter()
	if err != nil {
		out <- err
		return
	}
	res := objects.Download(c.ServiceClient(), c.container, c.name, nil)
	if res.Err != nil {
		lib.Log.Debugf("err from objects.Download: %s\n", res.Err)
	}
	dh, err := res.Extract()
	if err != nil {
		lib.Log.Debugf("err from res.Extract: %s\n", err)
		out <- err
		return
	}
	if c.ShouldProgress() {
		id := fmt.Sprintf("%s/%s", c.name, c.container)
		c.Sizes.Set(id, int(dh.ContentLength))
		out <- id

		go func() {
			_, err = io.Copy(w, res.Body)
			if err != nil {
				lib.Log.Debugf("Error copying (io.Reader) result: %s\n", err)
			}
		}()
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
	lib.Log.Debugf("wrote %d bytes to bytessentch", n)
	return
}

func (c *commandDownload) CustomWriter() (io.Writer, error) {
	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
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
		} else {
			return nil, err
		}
	}

done:
	fc := new(filecontent)
	fc.writer = f
	fc.bytessentch = c.ProgUpdateChIn()
	return fc, nil
}
