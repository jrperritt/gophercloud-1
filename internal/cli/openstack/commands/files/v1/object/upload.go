package object

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandUpload struct {
	ObjectV1Command
	traits.Pipeable
	traits.Progressable
	opts       objects.CreateOptsBuilder
	pipedField string
	proger     interfaces.Progresser
	//multiprogressable
}

//type multiprogressable struct {
//	proger interfaces.Progresser
//	traits.Progressable
//}

var (
	cUpload                            = new(commandUpload)
	_       interfaces.PipeCommander   = cUpload
	_       interfaces.StreamCommander = cUpload
	_       interfaces.BytesProgresser = cUpload

	_ interfaces.ReadBytesProgressItemer = new(uploaddata)

	flagsUpload = openstack.CommandFlags(cUpload)
)

var upload = cli.Command{
	Name:         "upload",
	Usage:        util.Usage(commandPrefix, "upload", "--container <containerName> --name <objectName>"),
	Description:  "Uploads an object",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpload) },
	Flags:        flagsUpload,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpload) },
}

func (c *commandUpload) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "container",
			Usage: "[required] The name of the container to upload the object into.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided with value of 'file'] The name the object should have in the Cloud Files container.",
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
		cli.StringFlag{
			Name:  "content-type",
			Usage: "[optional] The Content-Type header.",
		},
		cli.IntFlag{
			Name:  "content-length",
			Usage: "[optional] The Content-Length header.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
	}
}

type createopts struct {
	objects.CreateOpts
}

func (opts *createopts) ToObjectCreateParams() (io.Reader, map[string]string, string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return nil, nil, "", err
	}
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		return nil, nil, "", err
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	return opts.Content, h, q.String(), nil
}

func (c *commandUpload) HandleFlags() error {
	opts := new(createopts)
	opts.ContentLength = int64(c.Context().Int("content-length"))
	opts.ContentType = c.Context().String("content-type")

	if c.Context().IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	c.opts = opts
	c.pipedField = c.Context().String("stdin")

	return nil
}

func (c *commandUpload) HandlePipe(item string) (interface{}, error) {
	d := newuploaddata()
	switch c.pipedField {
	case "container":
		d.container = item
		switch c.Context().IsSet("file") {
		case true:
			f, err := os.Open(c.Context().String("file"))
			if err != nil {
				return nil, err
			}
			d.SetReader(f)
			switch c.Context().IsSet("name") {
			case true:
				d.object = c.Context().String("name")
			case false:
				d.object = f.Name()
			}
		case false:
			switch c.Context().IsSet("content") {
			case true:
				err := c.CheckFlagsSet([]string{"name"})
				if err != nil {
					return nil, err
				}
				d.object = c.Context().String("name")
				d.SetReader(strings.NewReader(c.Context().String("content")))
			case false:
				return nil, fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
			}
		}
	case "file":
		err := c.CheckFlagsSet([]string{"container"})
		if err != nil {
			return nil, err
		}
		d.container = c.Context().String("container")
		f, err := os.Open(item)
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		d.SetReader(f)
		d.object = f.Name()
		d.SetSize(fi.Size())
	case "content":
		err := c.CheckFlagsSet([]string{"container", "name"})
		if err != nil {
			return nil, err
		}
		d.SetReader(os.Stdin)
		d.container = c.Context().String("container")
		d.object = c.Context().String("name")
	}

	return d, nil
}

func (c *commandUpload) HandleSingle() (interface{}, error) {
	err := c.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return nil, err
	}

	d := newuploaddata()
	d.object = c.Context().String("name")
	d.container = c.Context().String("container")

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
		err = fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
	}

	return d, err
}

func (c *commandUpload) Execute(item interface{}, _ chan interface{}) {
	d := item.(*uploaddata)
	opts := *c.opts.(*createopts)

	defer func() {
		if closeable, ok := d.Reader().(io.ReadCloser); ok {
			closeable.Close()
		}
	}()

	opts.Content = d

	header, err := objects.Create(c.ServiceClient(), d.container, d.object, &opts).Extract()
	if err != nil {
		d.EndCh() <- err
		return
	}

	d.checksum = fmt.Sprintf("%x", d.hash.Sum(nil))
	if header.ETag != d.checksum {
		d.EndCh() <- fmt.Errorf("Different checksums: expected %s, got back %s\n", d.checksum, header.ETag)
		return
	}

	d.EndCh() <- fmt.Sprintf("Successfully uploaded object [%s] to container [%s]", d.object, d.container)
}

func (c *commandUpload) PipeFieldOptions() []string {
	return []string{"file", "container"}
}

func (c *commandUpload) HandleStream() (interface{}, error) {
	err := c.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return nil, err
	}

	d := newuploaddata()
	d.SetReader(os.Stdin)
	d.container = c.Context().String("container")
	d.object = c.Context().String("name")
	d.SetID(d.object)
	return d, nil
}

func (c *commandUpload) StreamFieldOptions() []string {
	return []string{"content"}
}

func (c *commandUpload) CreateBar(pi interfaces.ProgressItemer) interfaces.ProgressBarrer {
	if c.pipedField == "content" {
		proger := new(traits.BytesStreamProgressable)
		c.Progressable.Lock()
		proger.Progressable = c.Progressable
		c.Progressable.Unlock()
		return proger.CreateBar(pi)
	}
	proger := new(traits.BytesProgressable)
	c.Progressable.Lock()
	proger.Progressable = c.Progressable
	c.Progressable.Unlock()
	return proger.CreateBar(pi)
}

func newuploaddata() *uploaddata {
	d := new(uploaddata)
	d.hash = md5.New()
	d.ProgressItem.Init()
	return d
}

type uploaddata struct {
	traits.ProgressItemBytesRead
	container string
	object    string
	opts      createopts
	hash      hash.Hash
	checksum  string
}

func (d *uploaddata) Read(p []byte) (n int, err error) {
	n, err = d.Reader().Read(p)
	if err != nil {
		return
	}
	_, err = io.CopyN(d.hash, bytes.NewReader(p), int64(n))
	if err != nil {
		return
	}

	d.UpCh() <- n
	return
}
