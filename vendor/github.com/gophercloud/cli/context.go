package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/src/output"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
)

type Context struct {
	lib.Context
}

func (ctx Context) Name() string {
	return "stack"
}

// HandleCommand is the method that handles all commands. It accepts a Commander as
// a parameter, which all commands implement.
func (ctx Context) HandleCommand() error {
	// can the command accept input on STDIN?
	if pipeableCommand, ok := ctx.Commander.(PipeHandler); ok {
		// should we expect something on STDIN?
		if ctx.CLIContext.IsSet("stdin") {
			stdinField := ctx.CLIContext.String("stdin")
			// if so, does the given field accept pipeable input?
			if util.Contains(pipeableCommand.StdinField(), stdinField) {
				wg := sync.WaitGroup{}
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					item := scanner.Text()
					wg.Add(1)
					go func(resource Resource) {
						err := pipeableCommand.HandlePipe(&resource, item)
						if err != nil {
							resource.Err = fmt.Errorf("Error handling pipeable command on %s: %s\n", item, err)
							ctx.Results <- &resource
						} else {
							pipeableCommand.Execute(&resource)
							ctx.Results <- &resource
						}
						wg.Done()
					}(*resource)
				}
				if scanner.Err() != nil {
					resource.Err = scanner.Err()
					errExit1(command, resource)
				}
				wg.Wait()
				close(ctx.Results)
				// else, does the given command and field accept streaming input?
			} else if streamPipeableCommand, ok := pipeableCommand.(StreamPipeHandler); ok && util.Contains(streamPipeableCommand.StreamFields(), stdinField) {
				go func() {
					err := streamPipeableCommand.HandleStreamPipe(&resource)
					if err != nil {
						resource.Err = fmt.Errorf("Error handling streamable, pipeable command: %s\n", err)
					} else {
						streamPipeableCommand.Execute(&resource)
					}
					ctx.Results <- &resource
					close(ctx.Results)
				}()
			} else {
				// the value provided to the `stdin` flag is not valid
				resource.Err = fmt.Errorf("Unknown STDIN field: %s\n", stdinField)
				errExit1(command, resource)
			}
			// since no `stdin` flag was provided, treat as a singular execution
		} else {
			go func(resource Resource) {
				err := pipeableCommand.HandleSingle(&resource)
				if err != nil {
					resource.Err = err
					errExit1(command, &resource)
				}
				commander.Execute(&resource)
				ctx.Results <- &resource
				close(ctx.Results)
			}(*resource)
		}
		// the command is a single execution (as opposed to reading from a pipe)
	} else {
		go func(resource Resource) {
			commander.Execute(&resource)
			ctx.Results <- &resource
			close(ctx.Results)
		}(*resource)
	}

	//commander.Execute()

	resultsChan := commander.ResultsChan()
	for resource := range resultsChan {
		ctx.HandleResult()
	}
	return nil
}

func (ctx Context) ProcessResult(resource *Resource) error {
	// if an error was encountered during `handleExecution`, return it instead of
	// the `resource.Result`.
	if resource.Err != nil {
		ctx.CLIContext.App.Writer = os.Stderr
		ctx.Fields = []string{"error"}
		resource.Result = map[string]interface{}{"error": resource.Err.Error()}
		return nil
	}

	if resource.Result == nil {
		switch resource.Result.(type) {
		case []map[string]interface{}:
			resource.Result = fmt.Sprintf("No results found\n")
		default:
			resource.Result = fmt.Sprintf("No result found.\n")
		}
		return nil
	}

	// limit the returned fields if any were given in the `fields` flag
	ctx.LimitFields(resource)

	var err error
	// apply any output-specific transformations on the result
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
}

/*
var sampleMapStringInterface map[string]interface{} = nil
var sampleSliceMapStringInterface []map[string]interface{} = nil
var sampleIOReader io.Reader = nil

var outputMapping = map[string]map[reflect.Type]func(){
  "json": map[reflect.Type]func(){
    reflect.TypeOf(sampleMapStringInterface):
  },
}
*/

func (ctx Context) PrintResult(resource *Resource) error {
	w := ctx.CLIContext.App.Writer
	keys := resource.Keys
	noHeader := false
	if ctx.GlobalOptions.NoHeader {
		noHeader = true
	}
	switch r := resource.Result.(type) {
	case map[string]interface{}:
		m = onlyNonNil(r)
		switch ctx.GlobalOptions.output {
		case "json":
			output.MetadataJSON(w, r, keys)
		default:
			output.MetadataTable(w, r, keys)
		}
	case []map[string]interface{}:
		for i, m := range r {
			r[i] = onlyNonNil(m)
		}
		switch ctx.GlobalOptions.output {
		case "json":
			output.ListJSON(w, r, keys)
		default:
			output.ListTable(w, r, keys, noHeader)
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
			output.DefaultJSON(w, resource.Result)
		default:
			fmt.Fprintf(w, "%v\n", resource.Result)
		}
	}
}

// StoreCredentials caches the users auth credentials if available and the `no-cache`
// flag was not provided.
func (ctx Context) StoreCredentials() error {
	// if serviceClient is nil, the HTTP request for the command didn't get sent.
	// don't set cache if the `no-cache` flag is provided
	if ctx.ServiceClient != nil && !ctx.GlobalOptions.noCache {
		newCacheValue := &CacheItem{
			TokenID:         ctx.ServiceClient.TokenID,
			ServiceEndpoint: ctx.ServiceClient.Endpoint,
		}
		// get auth credentials
		credsResult, err := Credentials(ctx.CLIContext, nil)
		if err != nil {
			if ctx.logger != nil {
				ctx.logger.Infof("Error storing credentials in cache: %s\n", err)
			}
			return
		}
		ao := credsResult.AuthOpts
		region := credsResult.Region
		urlType := gophercloud.AvailabilityPublic
		// initialize the cache
		cache, err := InitCache()
		if err != nil {
			return err
		}
		// get the cache key
		cacheKey := cache.GetCacheKey()

		// set the cache value to the current values
		_ = cache.SetCacheValue(cacheKey, newCacheValue)
	}

	return nil
}

// ErrExit1 tells the CLI to print the error and exit.
func (ctx Context) ErrExit1(err error) {
	ctx.ProcessResult()
	ctx.PrintResult(&Resource{Err: err})
	os.Exit(1)
}

// LimitFields returns only the fields the user specified in the `fields` flag. If
// the flag wasn't provided, all fields are returned.
func (ctx Context) LimitFields(resource *Resource) {
	if ctx.CLIContext.IsSet("fields") {
		fields := strings.Split(strings.ToLower(ctx.CLIContext.String("fields")), ",")
		newKeys := []string{}
		for _, key := range resource.Keys {
			if util.Contains(fields, strings.ToLower(key)) {
				newKeys = append(newKeys, key)
			}
		}
		resource.Keys = newKeys
	}
}
