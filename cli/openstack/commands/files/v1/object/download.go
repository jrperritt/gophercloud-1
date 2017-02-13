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

type downloaddata struct {
	traits.ProgressItemBytesWrite
}

func newdownloaddata() *downloaddata {
	d := new(downloaddata)
	d.ProgressItem.Init()
	return d
}

func (d *downloaddata) Write(p []byte) (n int, err error) {
	n, err = d.Writer().Write(p)
	if err != nil {
		return
	}
	d.UpCh() <- n
	return
}

var (
	cDownload                            = new(commandDownload)
	_         interfaces.BytesProgresser = cDownload

	_ interfaces.WriteBytesProgressItemer = new(downloaddata)

	flagsDownload = openstack.CommandFlags(cDownload)
)

var download = cli.Command{
	Name:         "download",
	Usage:        util.Usage(commandPrefix, "download", "--container CONTAINER --name NAME"),
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
	}

	c.file = c.Context().String("file")

	if c.file != "" {
		_, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			if os.IsExist(err) {
				reader := bufio.NewReader(os.Stdin)
				for {
					fmt.Printf("A file named %s already exists. Overwrite? (y/n): ", c.file)
					choice, _ := reader.ReadString('\n')
					choice = strings.TrimSpace(choice)
					switch strings.ToLower(choice) {
					case "y", "yes":
						return nil
					case "n", "no":
						return fmt.Errorf("File (%s) already exists. Aborting...", c.file)
					default:
						continue
					}
				}
			}
		}
	}

	return nil
}

func (c *commandDownload) Execute(_ interface{}, out chan interface{}) {
	d := newdownloaddata()

	res := objects.Download(c.ServiceClient(), c.container, c.name, nil)
	if res.Err != nil {
		d.EndCh() <- res.Err
		return
	}

	dh, err := res.Extract()
	if err != nil {
		d.EndCh() <- err
		return
	}

	d.SetSize(dh.ContentLength)
	d.SetID(fmt.Sprintf("%s/%s", c.container, c.name))
	//	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_CREATE, 0666)
	f, err := os.Create(c.file)
	if err != nil {
		d.EndCh() <- err
		return
	}
	defer f.Close()
	d.SetWriter(f)

	c.ProgStartCh() <- d
	_, err = io.Copy(d, res.Body)
	if err != nil {
		d.EndCh() <- err
	}

	d.EndCh() <- fmt.Sprintf("Successfully downloaded %s to %s", d.ID(), c.file)
}
