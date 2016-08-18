package openstack

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

// Context satisfies the Provider interface
type Context struct {
	outChannel chan interface{}
	commander  lib.Commander
}

// Name satisfies the Provider.Name method
func (c *Context) Name() string {
	return "openstack"
}

// NewGlobalOptionser satisfies the Provider.NewGlobalOptionser method
func (c *Context) NewGlobalOptionser(context lib.Contexter) lib.GlobalOptionser {
	g := new(GlobalOptions)
	g.cliContext = context.(*cli.Context)
	return g
}

// NewAuthenticater satisfies the Provider.NewAuthenticater method
func (c *Context) NewAuthenticater(globalOptionser lib.GlobalOptionser, serviceType string) lib.Authenticater {
	globalOptions := globalOptionser.(*GlobalOptions)

	return &auth{
		AuthOptions: &gophercloud.AuthOptions{
			Username:         globalOptions.username,
			UserID:           globalOptions.userID,
			Password:         globalOptions.password,
			TenantID:         globalOptions.authTenantID,
			TokenID:          globalOptions.authToken,
			IdentityEndpoint: globalOptions.authURL,
		},
		logger:      globalOptions.logger,
		noCache:     globalOptions.noCache,
		serviceType: serviceType,
		region:      globalOptions.region,
		profile:     globalOptions.profile,
	}
}

func (c *Context) InputChannel() chan interface{} {
	return make(chan interface{})
}

func (c *Context) FillInputChannel(commander lib.Commander, in chan interface{}) {
	defer close(in)
	ctx := commander.Ctx().(*cli.Context)
	switch t := commander.(type) {
	case lib.PipeCommander:
		switch ctx.IsSet("stdin") {
		case true:
			stdin := ctx.String("stdin")
			switch util.Contains(t.PipeFieldOptions(), stdin) {
			case true:
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					item, err := t.HandlePipe(scanner.Text())
					switch err {
					case nil:
						in <- item
					default:
						c.outChannel <- err
					}
				}
				if scanner.Err() != nil {
					c.outChannel <- scanner.Err()
				}
			default:
				c.outChannel <- fmt.Errorf("Unknown STDIN field: %s\n", stdin)
			}
		default:
			item, err := t.HandleSingle()
			switch err {
			case nil:
				in <- item
			default:
				c.outChannel <- err
			}
		}
	case lib.StreamPipeCommander:
		switch ctx.IsSet("stdin") {
		case true:
			stdin := ctx.String("stdin")
			switch util.Contains(t.StreamFieldOptions(), stdin) {
			case true:
				stream, err := t.HandleStreamPipe(os.Stdin)
				switch err {
				case nil:
					in <- stream
				default:
					c.outChannel <- err
				}
			default:
				c.outChannel <- fmt.Errorf("Unknown STDIN field: %s\n", stdin)
			}
		default:
			item, err := t.HandleSingle()
			switch err {
			case nil:
				in <- item
			default:
				c.outChannel <- err
			}
		}
	default:
		in <- 0
	}
}

func (c *Context) ResultsChannel() chan interface{} {
	c.outChannel = make(chan interface{})
	return c.outChannel
}

// NewResultOutputter satisfies the Provider.NewResultOutputter method
func (c *Context) NewResultOutputter(globalOptionser lib.GlobalOptionser, commander lib.Commander) lib.Outputter {
	globalOptions := globalOptionser.(*GlobalOptions)

	return &output{
		fields:    globalOptions.fields,
		noHeader:  globalOptions.noHeader,
		format:    globalOptions.outputFormat,
		logger:    globalOptions.logger,
		commander: commander,
	}
}

func (c *Context) ErrExit1(err error) {
	fmt.Println(err)
	os.Exit(1)
	//panic(err)
}
