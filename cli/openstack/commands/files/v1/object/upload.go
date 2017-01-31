package object

import (
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

type commandUpload struct {
	ObjectV1Command
	traits.Progressable
	opts       objects.CreateOptsBuilder
	pipedField string
}

type pipeData struct {
	container string
	object    string
	content   io.Reader
}

var (
	cUpload                                = new(commandUpload)
	_       interfaces.StreamPipeCommander = cUpload
	_       interfaces.Progresser          = cUpload

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

func (c *commandUpload) HandleFlags() error {
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

	if c.Context.IsSet("stdin") {
		c.pipedField = c.Context.String("stdin")
	}

	return nil
}

func (c *commandUpload) HandlePipe(item string) (interface{}, error) {
	pd := new(pipeData)
	switch c.pipedField {
	case "container":
		pd.container = item
		switch c.Context.IsSet("file") {
		case true:
			f, err := os.Open(c.Context.String("file"))
			if err != nil {
				return nil, err
			}
			pd.content = f
			switch c.Context.IsSet("name") {
			case true:
				pd.object = c.Context.String("name")
			case false:
				pd.object = f.Name()
			}
		case false:
			switch c.Context.IsSet("content") {
			case true:
				err := c.CheckFlagsSet([]string{"name"})
				if err != nil {
					return nil, err
				}
				pd.object = c.Context.String("name")
				pd.content = strings.NewReader(c.Context.String("content"))
			case false:
				return nil, fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
			}
		}
	case "file":
		err := c.CheckFlagsSet([]string{"container"})
		if err != nil {
			return nil, err
		}
		pd.container = c.Context.String("container")
		f, err := os.Open(item)
		if err != nil {
			return nil, err
		}
		pd.content = f
		pd.object = f.Name()
	}

	return pd, nil
}

func (c *commandUpload) HandleStreamPipe(stream io.Reader) (interface{}, error) {
	err := c.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return nil, err
	}
	pd := new(pipeData)
	pd.content = stream
	pd.container = c.Context.String("container")
	pd.object = c.Context.String("name")
	return pd, nil
}

func (c *commandUpload) HandleSingle() (interface{}, error) {
	err := c.CheckFlagsSet([]string{"container", "name"})
	if err != nil {
		return nil, err
	}

	pd := new(pipeData)
	pd.object = c.Context.String("name")
	pd.container = c.Context.String("container")

	switch c.Context.IsSet("file") {
	case true:
		pd.content, err = os.Open(c.Context.String("file"))
	case false:
		switch c.Context.IsSet("content") {
		case true:
			pd.content = strings.NewReader(c.Context.String("content"))
		case false:
			err = fmt.Errorf("One of `--file` and `--content` must be provided if not piping to STDIN")
		}
	}

	return pd, err
}

type bytescontent struct {
	reader        io.Reader
	bytessentchan chan (int)
}

func (b *bytescontent) Read(p []byte) (n int, err error) {
	n, err = b.reader.Read(p)
	if err != nil {
		return
	}
	b.bytessentchan <- n
	return
}

func (c *commandUpload) Execute(item interface{}, out chan interface{}) {
	pd := item.(*pipeData)

	reader, ok := pd.content.(io.Reader)
	if !ok {
		out <- fmt.Errorf("Expected an io.Reader but instead got %T", item)
	}

	defer func() {
		if closeable, ok := reader.(io.ReadCloser); ok {
			closeable.Close()
		}
	}()

	bc := new(bytescontent)
	bc.bytessentchan = make(chan int)
	bc.reader = reader
	c.opts.(*objects.CreateOpts).Content = bc

	var m map[string]interface{}
	err := objects.Create(c.ServiceClient, pd.container, pd.object, c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- fmt.Sprintf("Successfully uploaded object [%s] to container [%s]", pd.object, pd.container)
	default:
		out <- err
	}
}

func (c *commandUpload) PipeFieldOptions() []string {
	return []string{"file", "container"}
}

func (c *commandUpload) StreamFieldOptions() []string {
	return []string{"content"}
}

func (c *commandUpload) InitProgress() {
	c.ProgressInfo = openstack.NewProgressInfo(1)
	c.Progressable.InitProgress()
}

func (c *commandUpload) ShowBar(id string) {
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
			s.Increment = int(r.(float64))
			c.ProgressInfo.UpdateChan <- s
		}
	}
}
