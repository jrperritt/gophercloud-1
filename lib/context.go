package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"

	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
)

// Provider should be implemented by one object per cloud provider
type Provider interface {
	Name() string
	SetGlobalOptions(*cli.Context) error
	InitProvider(GlobalOptionser) error
	Authenticate() error
	GetResultsChannel() chan Resulter
	RunCommand(Commander, chan Resulter) error
	HandleResults(chan Resulter) error
	ProcessResult(Resulter) error
	PrintResult(Resulter) error
	CleanUp() error
	ErrExit1(error)

	NewGlobalOptions() GlobalOptionser
	NewResult() Resulter
}

// Cloud is a global variable representing the CLI's global context. It should
// be set (read: overridden) by each cloud provider
var Cloud Provider

// Run executes all the methods of a Provider for each command
func Run(cliContext *cli.Context, command Commander) {

	if Cloud == nil {
		panic("You must set the CloudProvider variable")
	}

	globalOptions := Cloud.NewGlobalOptions()

	err := Cloud.SetGlobalOptions(cliContext)
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.InitProvider(globalOptions)
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.Authenticate()
	if err != nil {
		Cloud.ErrExit1(err)
	}

	resultsChannel := Cloud.GetResultsChannel()

	err = Cloud.RunCommand(command, resultsChannel)
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.HandleResults(resultsChannel)
	if err != nil {
		Cloud.ErrExit1(err)
	}

	err = Cloud.CleanUp()
	if err != nil {
		Cloud.ErrExit1(err)
	}
}

// Context represents a convenience object that implements Provider. cloud
// providers using this library can embed Context in their own objects that
// implement Provider to make use of its methods
type Context struct {
	cliContext    *cli.Context
	globalOptions GlobalOptionser
	auth          Authenticater
	//command       Commander
	resource      Resourcer
	output        Outputter
	serviceClient *gophercloud.ServiceClient
	logger        *logrus.Logger
}

// SetGlobalOptions satisfies the Provider.SetGlobalOptions method
func (c *Context) SetGlobalOptions(cliContext *cli.Context) (GlobalOptionser, error) {

	if c.globalOptions == nil {
		return nil, nil
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
		return nil, multiErr
	}

	err = c.globalOptions.Set()
	if err != nil {
		return nil, err
	}

	return c.globalOptions, nil
}

// InitProvider satisfies the Provider.InitProvider method
func (c *Context) InitProvider(_ GlobalOptionser) error {

	//c.cliContext = cliContext
	//c.command = cmd
	return nil
}

// Authenticate satisfies the Provider.Authenticate method
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

// RunCommand satisfies the Provider.RunCommand method
func (c *Context) RunCommand(command Commander, resultsChannel chan Resulter) error {
	err := command.HandleFlags()
	if err != nil {
		return err
	}

	// can the command accept input on STDIN?
	if pipeableCommand, ok := command.(PipeHandler); ok {
		// should we expect something on STDIN?
		if c.cliContext.IsSet("stdin") {
			stdinField := c.cliContext.String("stdin")
			// if so, does the given field accept pipeable input?
			if util.Contains(pipeableCommand.StdinFields(), stdinField) {
				wg := sync.WaitGroup{}
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					item := scanner.Text()
					wg.Add(1)
					go func() {
						err := pipeableCommand.HandlePipe(c.resource, item)
						if err != nil {
							c.resource.Result().SetError(fmt.Errorf("Error handling pipeable command on %s: %s\n", item, err))
							command.ResultsChan() <- c.resource.Result()
						} else {
							pipeableCommand.Execute(c.resource)
							command.ResultsChan() <- c.resource.Result()
						}
						wg.Done()
					}()
				}
				if scanner.Err() != nil {
					c.ErrExit1(scanner.Err())
				}
				wg.Wait()
				close(command.ResultsChan())
				// else, does the given command and field accept streaming input?
			} else if streamPipeableCommand, ok := pipeableCommand.(StreamPipeHandler); ok && util.Contains(streamPipeableCommand.StreamFields(), stdinField) {
				go func() {
					err := streamPipeableCommand.HandleStreamPipe(c.resource)
					if err != nil {
						c.resource.Result().SetError(fmt.Errorf("Error handling streamable, pipeable command: %s\n", err))
					} else {
						streamPipeableCommand.Execute(c.resource)
					}
					command.ResultsChan() <- c.resource.Result()
					close(command.ResultsChan())
				}()
			} else {
				// the value provided to the `stdin` flag is not valid
				c.ErrExit1(fmt.Errorf("Unknown STDIN field: %s\n", stdinField))
			}
			// since no `stdin` flag was provided, treat as a singular execution
		} else {
			go func() {
				err := pipeableCommand.HandleSingle(c.resource)
				if err != nil {
					c.ErrExit1(err)
				}
				command.Execute(c.resource)
				command.ResultsChan() <- c.resource.Result()
				close(command.ResultsChan())
			}()
		}
		// the command is a single execution (as opposed to reading from a pipe)
	} else {
		go func() {
			command.Execute(c.resource)
			command.ResultsChan() <- c.resource.Result()
			close(command.ResultsChan())
		}()
	}

	return nil
}

// HandleResults satisfies the Provider.HandleResults method
func (c Context) HandleResults(resultsChannel chan Resulter) error {
	for result := range resultsChannel {
		err := c.ProcessResult(result)
		if err != nil {
			return err
		}
		err = c.PrintResult(result)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProcessResult satisfies the Provider.ProcessResult method
func (c Context) ProcessResult(result Resulter) error {

	if result.GetError() != nil {
		c.cliContext.App.Writer = os.Stderr
		return nil
	}

	if result.GetValue() == nil {
		for _, t := range result.Types() {
			if reflect.TypeOf(result.GetValue()).AssignableTo(reflect.TypeOf(t)) {
				v, err := t.HandleEmpty()
				if err != nil {
					return err
				}
				result.SetValue(v)
			}
		}
	}

	return nil
}

// PrintResult satisfies the Provider.PrintResult method
func (c Context) PrintResult(result Resulter) error {
	//c.output.
	w := c.cliContext.App.Writer
	keys := resource.Keys
	noHeader := false
	if ctx.GlobalOptions.NoHeader {
		noHeader = true
	}

	// limit the returned fields if any were given in the `fields` flag
	c.LimitFields(resource)

	switch ctx.GlobalOptions.output {
	case "json":
		if jsoner, ok := command.(PreJSONer); ok {
			err = jsoner.PreJSON(resource)
		}
	default:
		if tabler, ok := command.(PreTabler); ok {
			err = tabler.PreTable(resource)
		}
	}
	if err != nil {
		resource.Keys = []string{"error"}
		resource.Result = map[string]interface{}{"error": err.Error()}
	}

	switch r := resource.Result.(type) {
	case map[string]interface{}:
		m = onlyNonNil(r)
		switch ctx.GlobalOptions.output {
		case "json":
			MetadataJSON(w, r, keys)
		default:
			MetadataTable(w, r, keys)
		}
	case []map[string]interface{}:
		for i, m := range r {
			r[i] = onlyNonNil(m)
		}
		switch ctx.GlobalOptions.output {
		case "json":
			ListJSON(w, r, keys)
		default:
			ListTable(w, r, keys, noHeader)
		}
	case io.Reader:
		if _, ok := resource.Result.(io.ReadCloser); ok {
			defer resource.Result.(io.ReadCloser).Close()
		}
		_, err := io.Copy(w, r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying (io.Reader) result: %s\n", err)
		}
	default:
		switch ctx.GlobalOptions.output {
		case "json":
			DefaultJSON(w, resource.Result)
		default:
			fmt.Fprintf(w, "%v\n", resource.Result)
		}
	}
}

// CleanUp satisfies the Provider.CleanUp method
func (c Context) CleanUp() error {
	if authFromCacher, ok := c.auth.(AuthFromCacher); ok {
		err = authFromCacher.StoreCredentials()
		if err != nil {
			return err
		}
	}
	return nil
}

// ErrExit1 satisfies the Provider.ErrExit1 method
func (c Context) ErrExit1(err error) {
	var result Resulter
	result.SetError(err)
	err = c.CleanUp()
	if err != nil {

	}
	os.Exit(1)
}
