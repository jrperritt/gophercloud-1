package lib

import (
	"fmt"

	"github.com/gophercloud/cli/vendor/github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
)

type Provider interface {
	Name() string
	Init(*cli.Context, Commander) error
	HandleGlobalOptions() error
	Authenticate() error
	RunCommand() error
	HandleResults() error
	CleanUp() error
	ErrExit1(error)
}

var Cloud Provider

func Run(cliContext *cli.Context, cmd Commander) {

	if Cloud == nil {
		panic("You must set the CloudProvider variable")
	}

	err := Cloud.Init(cliContext, cmd)
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.HandleGlobalOptions()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.Authenticate()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.RunCommand()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.HandleResults()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.CleanUp()
	if err != nil {
		Cloud.ErrExit1(err)
	}
}

type Context struct {
	cliContext    *cli.Context
	globalOptions GlobalOptionser
	auth          Authenticater
	command       Commander
	result        Resulter
	serviceClient *gophercloud.ServiceClient
	logger        *logrus.Logger
}

func (c *Context) Init(cliContext *cli.Context, cmd Commander) error {
	// there shouldn't be any arguments after the last subcommand.
	if lenArgs := len(cliContext.Args()); lenArgs != 0 {
		return fmt.Errorf("Expected %d args but got %d\nUsage: %s\n", 0, lenArgs, cliContext.Command.Usage)
	}
	c.cliContext = cliContext
	c.command = cmd
	return nil
}

func (c *Context) HandleGlobalOptions() error {

	if c.globalOptions == nil {
		return nil
	}

	err := c.globalOptions.InitGlobalOptions()
	if err != nil {
		return err
	}

	// we may get multiple errors while trying to handle the global options
	// so we'll try to return all of them at once, instead of returning just one,
	// only return a different one after that one's been rectified.
	multiErr := make(MultiError, 0)

	// for each source where a user could provide a global option,
	// parse the options from that source. sources will be parsed in the order
	// in which they appear in the Sources method
	for _, source := range c.globalOptions.Sources() {
		if parseOptions := c.globalOptions.MethodsMap()[source]; parseOptions != nil {
			err := parseOptions()
			if err != nil {
				multiErr = append(multiErr, err)
			}
		}
	}

	// after the global options have been parsed, run each global option's
	// validation function, if it exists
	err = c.globalOptions.Validate()
	if err != nil {
		multiErr = append(multiErr, err)
	}

	if len(multiErr) > 0 {
		return multiErr
	}

	err = c.globalOptions.Set()
	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates the user and acquires a client with which to make
// requests
func (c *Context) Authenticate() error {
	var client *gophercloud.ServiceClient
	var err error

	if authFromCacher, ok := c.auth.(AuthFromCacher); ok {
		client, err = authFromCacher.AuthFromCache()
		if err != nil {
			return err
		}
	}

	if client == nil {
		client, err = c.auth.AuthFromScratch()
		if err != nil {
			return err
		}
	}

	//client.HTTPClient.Transport.(*LogRoundTripper).Logger = c.logger
	c.serviceClient = client
	return nil
}

func (c *Context) RunCommand() error {
	err := c.command.HandleFlags()
	if err != nil {
		return err
	}
	return nil
}
