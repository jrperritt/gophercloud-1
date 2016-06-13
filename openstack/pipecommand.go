package openstack

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
)

type PipeCommand struct {
	Command
	stdinField string
}

func (c *PipeCommand) RunCommand(resultsChannel chan lib.Resulter) error {
	r := new(Resource)
	if stdinField := c.String("stdin"); stdinField != "" {
		r.SetStdinField(stdinField)
		// if so, does the given field accept pipeable input?
		if util.Contains(c.PipeFieldOptions(), stdinField) {
			wg := sync.WaitGroup{}
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				r := new(Resource)
				r.SetStdinValue(scanner.Text())
				wg.Add(1)
				go func() {
					resultsChannel <- c.Execute(r)
					wg.Done()
				}()
			}
			if scanner.Err() != nil {
				return scanner.Err()
			}
			wg.Wait()
			close(resultsChannel)
			return nil
		}
		// does the given command and field accept streaming input?
		if streamPipeableCommand, ok := interface{}(c).(lib.StreamPipeCommander); ok && util.Contains(streamPipeableCommand.StreamFieldOptions(), stdinField) {
			r.SetStdinValue(os.Stdin)
			resultsChannel <- c.Execute(r)
			return nil
		}
		// the value provided to the `stdin` flag is not valid
		return fmt.Errorf("Unknown STDIN field: %s\n", stdinField)
	}

	err := c.HandleSingle(new(Resource))
	if err != nil {
		return err
	}

	resultsChannel <- c.Execute(new(Resource))
	return nil
}

func (c *PipeCommand) HandlePipe(_ lib.Resourcer, _ string) error {
	return nil
}

/*
func (c *PipeCommand) PipeField() string {
	return c.stdinField
}

func (c *PipeCommand) SetPipeField(f string) {
	c.stdinField = f
}
*/

func (c *PipeCommand) HandleSingle(_ lib.Resourcer) error {
	return nil
}

func (c *PipeCommand) PipeFieldOptions() []string {
	return nil
}
