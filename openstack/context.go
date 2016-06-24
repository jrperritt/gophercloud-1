package openstack

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud"
	"github.com/jrperritt/rack/util"
)

// Context satisfies the Provider interface
type Context struct {
	outChannel chan interface{}
}

// Name satisfies the Provider.Name method
func (c *Context) Name() string {
	return "openstack"
}

// NewGlobalOptionser satisfies the Provider.NewGlobalOptionser method
func (c *Context) NewGlobalOptionser(context lib.Context) lib.GlobalOptionser {
	g := new(GlobalOptions)
	g.cliContext = context.(*cli.Context)
	return g
}

// NewAuthenticater satisfies the Provider.NewAuthenticater method
func (c *Context) NewAuthenticater(globalOptionser lib.GlobalOptionser, serviceType string) lib.Authenticater {
	globalOptions := globalOptionser.(*GlobalOptions)

	return auth{
		authOptions: &gophercloud.AuthOptions{
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
	ctx := commander.Ctx()
	switch t := commander.(type) {
	case lib.PipeCommander:
		switch ctx.IsSet("stdin") {
		case true:
			stdin := ctx.String("stdin")
			switch util.Contains(t.PipeFieldOptions(), stdin) {
			case true:
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					in <- t.HandlePipe(scanner.Text())
				}
				if scanner.Err() != nil {
					c.outChannel <- scanner.Err()
				}
			default:
				c.outChannel <- fmt.Errorf("Unknown STDIN field: %s\n", stdin)
			}
		default:
			in <- t.HandleSingle()
		}
	case lib.StreamPipeCommander:
		switch ctx.IsSet("stdin") {
		case true:
			stdin := ctx.String("stdin")
			switch util.Contains(t.StreamFieldOptions(), stdin) {
			case true:
				in <- t.HandleStreamPipe(stdin)
			default:
				c.outChannel <- fmt.Errorf("Unknown STDIN field: %s\n", stdin)
			}
		default:
			in <- t.HandleSingle()
		}
	default:
		in <- 0
	}
	close(in)
}

func (c *Context) ResultsChannel() chan interface{} {
	ch := make(chan interface{})
	c.outChannel = ch
	return ch
}

// NewResultOutputter satisfies the Provider.NewResultOutputter method
func (c *Context) NewResultOutputter(globalOptionser lib.GlobalOptionser) lib.Outputter {
	globalOptions := globalOptionser.(*GlobalOptions)

	return output{
		//fields:   globalOptions.fields,
		noHeader: globalOptions.noHeader,
		format:   globalOptions.outputFormat,
		logger:   globalOptions.logger,
	}
}

func (c *Context) ErrExit1(err error) {
	panic(err)
}
