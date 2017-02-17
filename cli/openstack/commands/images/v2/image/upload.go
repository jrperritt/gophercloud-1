package image

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"gopkg.in/urfave/cli.v1"
)

type CommandUpload struct {
	ImageV2Command
	traits.Progressable
	pipedField string
	proger     interfaces.Progresser
	id         string
}

var (
	cUpload                            = new(CommandUpload)
	_       interfaces.StreamCommander = cUpload
	_       interfaces.BytesProgresser = cUpload

	_ interfaces.ReadBytesProgressItemer = new(uploaddata)

	flagsUpload = openstack.CommandFlags(cUpload)
)

var upload = cli.Command{
	Name:         "upload",
	Usage:        util.Usage(commandPrefix, "upload", "[--id ID | --name NAME]"),
	Description:  "Uploads image data to an already-existing image",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpload) },
	Flags:        flagsUpload,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpload) },
}

func (c *CommandUpload) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Usage: "[optional; required if `name` isn't provided] The ID of the image",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `id` isn't provided] The name of the image",
		},
		cli.StringFlag{
			Name:  "content",
			Usage: "[optional; required if `file` or `stdin` isn't provided] The string contents to upload.",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "[optional; required if `content` or `stdin` isn't provided] The file name containing the contents to upload.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `file` or `content` isn't provided] The field being piped to STDIN, if any. Valid values are: file, container, content.",
		},
	}
}

func (c *CommandUpload) HandleFlags() error {
	c.pipedField = c.Context().String("stdin")

	id, err := c.IDOrName(IDFromName)
	if err != nil {
		return err
	}

	c.id = id

	return nil
}

func (c *CommandUpload) HandleSingle() (interface{}, error) {
	d := newuploaddata()
	d.SetID(c.id)
	if ok := c.Context().IsSet("file"); ok {
		f, err := os.Open(c.Context().String("file"))
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		d.SetSize(fi.Size())
		d.SetReader(f)
	} else if ok := c.Context().IsSet("content"); ok {
		r := strings.NewReader(c.Context().String("content"))
		d.SetSize(int64(r.Len()))
		d.SetReader(r)
	} else {
		return nil, fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
	}

	return d, nil
}

func (c *CommandUpload) Execute(item interface{}, out chan interface{}) {
	d := item.(*uploaddata)
	err := imagedata.Upload(c.ServiceClient(), d.ID(), d.Reader()).ExtractErr()
	if err != nil {
		d.EndCh() <- err
		return
	}

	d.EndCh() <- fmt.Sprintf("Successfully uploaded image data to image [%s]", d.ID())
}

func (c *CommandUpload) HandleStream() (interface{}, error) {
	d := newuploaddata()
	d.SetID(c.id)
	d.SetReader(os.Stdin)
	return d, nil
}

func (c *CommandUpload) StreamFieldOptions() []string {
	return []string{"content"}
}

func (c *CommandUpload) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	if c.pipedField == "content" {
		proger := new(traits.BytesStreamProgressable)
		proger.Progressable = c.Progressable
		return proger.CreateBar(pi)
	}
	proger := new(traits.BytesProgressable)
	proger.Progressable = c.Progressable
	return proger.CreateBar(pi)
}

func newuploaddata() *uploaddata {
	d := new(uploaddata)
	d.ProgressItem.Init()
	return d
}

type uploaddata struct {
	traits.ProgressItemBytesRead
}

func (d *uploaddata) Read(p []byte) (n int, err error) {
	n, err = d.Reader().Read(p)
	if err != nil {
		return
	}

	d.UpCh() <- n
	return
}
