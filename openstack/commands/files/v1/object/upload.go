package object

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"gopkg.in/urfave/cli.v1"
)

type commandUpload struct {
	openstack.CommandUtil
	ObjectV1Command
	container string
	object    string
	stream    io.Reader
	opts      objects.CreateOptsBuilder
}

var (
	cUpload                         = new(commandUpload)
	_       lib.StreamPipeCommander = cUpload
	_       lib.Waiter              = cUpload

	flagsUpload = openstack.CommandFlags(cUpload)
)

var upload = cli.Command{
	Name:         "upload",
	Usage:        util.Usage(commandPrefix, "upload", "--container <containerName> --name <objectName>"),
	Description:  "Uploads an object",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpload) },
	Flags:        flagsUpload,
	BashComplete: func(_ *cli.Context) { openstack.BashComplete(flagsUpload) },
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

func (c *commandUpload) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"container"})
	if err != nil {
		return err
	}
	c.container = c.Context.String("container")

	c.Wait = c.Context.IsSet("wait")
	c.Quiet = c.Context.IsSet("quiet")

	opts := &objects.CreateOpts{
		ContentLength: int64(c.Context.Int("content-length")),
		ContentType:   c.Context.String("content-type"),
	}

	if c.Context.IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	c.opts = opts

	return nil
}

func (c *commandUpload) HandlePipe(item string) (interface{}, error) {
	f, err := os.Open(item)
	if err != nil {
		return nil, err
	}
	c.object = f.Name()
	return f, nil
}

func (c *commandUpload) HandleStreamPipe(stream io.Reader) (io.Reader, error) {
	err := c.CheckFlagsSet([]string{"name"})
	if err != nil {
		return nil, err
	}
	c.object = c.Context.String("name")
	return stream, nil
}

func (c *commandUpload) HandleSingle() (s interface{}, err error) {
	err = c.CheckFlagsSet([]string{"name"})
	if err != nil {
		return
	}
	c.object = c.Context.String("name")

	switch c.Context.IsSet("file") {
	case true:
		s, err = os.Open(c.Context.String("file"))
	case false:
		switch c.Context.IsSet("content") {
		case true:
			s = strings.NewReader(c.Context.String("content"))
		case false:
			err = fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
		}
	}

	return
}

func (c *commandUpload) Execute(in, out chan interface{}) {
	defer close(out)
	for item := range in {
		reader, ok := item.(io.Reader)
		if !ok {
			out <- fmt.Errorf("Expected an io.Reader but instead got %v", reflect.TypeOf(item))
		}
		c.opts.(*objects.CreateOpts).Content = reader
		defer func() {
			if closeable, ok := reader.(io.ReadCloser); ok {
				closeable.Close()
			}
		}()
		var m map[string]interface{}
		err := objects.Create(c.ServiceClient, c.container, c.object, c.opts).ExtractInto(&m)
		switch err {
		case nil:
			out <- fmt.Sprintf("Successfully uploaded object [%s] to container [%s]", c.object, c.container)
		default:
			out <- err
		}
	}
}

func (c *commandUpload) PipeFieldOptions() []string {
	return []string{"file"}
}

func (c *commandUpload) StreamFieldOptions() []string {
	return []string{"content"}
}

func (c *commandUpload) ExecuteAndWait(in, out chan interface{}) {
	openstack.ExecuteAndWait(c, in, out)
}
