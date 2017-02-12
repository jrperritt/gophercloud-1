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
	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandUpload struct {
	ObjectV1Command
	traits.BytesProgressable
	opts       objects.CreateOptsBuilder
	pipedField string
}

type pipedata struct {
	*traits.ProgressItemBytes
	container   string
	object      string
	size        int64
	reader      io.Reader
	bytessentch chan interface{}
	hash        hash.Hash
	checksum    string
}

func newpipedata() *pipedata {
	pd := new(pipedata)
	pd.hash = md5.New()
	pd.bytessentch = make(chan interface{})
	endch := make(chan interface{})
	pd.ProgressItemBytes = new(traits.ProgressItemBytes)
	pd.SetEndCh(endch)
	return pd
}

func (b *pipedata) Read(p []byte) (n int, err error) {
	n, err = b.reader.Read(p)
	if err != nil {
		return
	}
	_, err = io.CopyN(b.hash, bytes.NewReader(p), int64(n))
	if err != nil {
		return
	}

	b.bytessentch <- n
	return
}

func (b *pipedata) UpCh() chan interface{} {
	return b.bytessentch
}

func (b *pipedata) ID() string {
	return b.object
}

func (b *pipedata) Size() int64 {
	return b.size
}

var (
	cUpload                          = new(commandUpload)
	_       interfaces.PipeCommander = cUpload
	_       interfaces.Progresser    = cUpload

	_ interfaces.ProgressItemer = new(pipedata)

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
			Usage: "[optional; required if `file` or `content` isn't provided] The field being piped to STDIN, if any. Valid values are: file, content.",
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

type creatopts struct {
	objects.CreateOpts
}

func (opts *creatopts) ToObjectCreateParams() (io.Reader, map[string]string, string, error) {
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
	opts := new(creatopts)
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
	pd := newpipedata()
	switch c.pipedField {
	case "container":
		pd.container = item
		switch c.Context().IsSet("file") {
		case true:
			f, err := os.Open(c.Context().String("file"))
			if err != nil {
				return nil, err
			}
			pd.reader = f
			switch c.Context().IsSet("name") {
			case true:
				pd.object = c.Context().String("name")
			case false:
				pd.object = f.Name()
			}
		case false:
			switch c.Context().IsSet("content") {
			case true:
				err := c.CheckFlagsSet([]string{"name"})
				if err != nil {
					return nil, err
				}
				pd.object = c.Context().String("name")
				pd.reader = strings.NewReader(c.Context().String("content"))
			case false:
				return nil, fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
			}
		}
	case "file":
		err := c.CheckFlagsSet([]string{"container"})
		if err != nil {
			return nil, err
		}
		pd.container = c.Context().String("container")
		f, err := os.Open(item)
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		pd.reader = f
		pd.object = f.Name()
		pd.size = fi.Size()
	case "content":
		err := c.CheckFlagsSet([]string{"container", "name"})
		if err != nil {
			return nil, err
		}
		pd.reader = os.Stdin
		pd.container = c.Context().String("container")
		pd.object = c.Context().String("name")
	}

	return pd, nil
}

func (c *commandUpload) HandleSingle() (interface{}, error) {
	err := c.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return nil, err
	}

	pd := newpipedata()
	pd.object = c.Context().String("name")
	pd.container = c.Context().String("container")

	if ok := c.Context().IsSet("file"); ok {
		f, err := os.Open(c.Context().String("file"))
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		pd.size = fi.Size()
		pd.reader = f
	} else if ok := c.Context().IsSet("content"); ok {
		r := strings.NewReader(c.Context().String("content"))
		pd.size = int64(r.Len())
		pd.reader = r
	} else {
		err = fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
	}

	return pd, err
}

func (c *commandUpload) Execute(item interface{}, out chan interface{}) {
	pd := item.(*pipedata)

	defer func() {
		if closeable, ok := pd.reader.(io.ReadCloser); ok {
			closeable.Close()
		}
	}()

	c.opts.(*creatopts).Content = pd

	lib.Log.Debugln("running objects.Create...")
	header, err := objects.Create(c.ServiceClient(), pd.container, pd.object, c.opts).Extract()
	close(pd.UpCh())
	if err != nil {
		lib.Log.Debugf("err from objects.Create: %+v", err)
		pd.EndCh() <- err
		return
	}

	pd.checksum = fmt.Sprintf("%x", pd.hash.Sum(nil))
	if header.ETag != pd.checksum {
		lib.Log.Debugf("Different checksums: expected %s, got back %s\n", pd.checksum, header.ETag)
		pd.EndCh() <- fmt.Errorf("Different checksums: expected %s, got back %s\n", pd.checksum, header.ETag)
		return
	}

	lib.Log.Debugln("successfully uploaded object")
	pd.EndCh() <- fmt.Sprintf("Successfully uploaded object [%s] to container [%s]", pd.object, pd.container)
}

func (c *commandUpload) PipeFieldOptions() []string {
	return []string{"file", "container", "content"}
}
