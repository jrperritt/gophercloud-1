package openstack

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"gopkg.in/urfave/cli.v1"
)

type Commander interface {
	HandleFlags() error
	Execute(item interface{}, out chan interface{})
	Flags() []cli.Flag
	SetServiceClient(*gophercloud.ServiceClient) error
	SetContext(*cli.Context) error
	ServiceType() string
}

type PipeCommander interface {
	Commander
	Waiter
	HandleSingle() (interface{}, error)
	HandlePipe(string) (interface{}, error)
	PipeFieldOptions() []string
}

type StreamPipeCommander interface {
	PipeCommander
	HandleStreamPipe(io.Reader) (interface{}, error)
	StreamFieldOptions() []string
}

type Waiter interface {
	WaitFor(item interface{})
	ShouldWait() bool
	WaitFlags() []cli.Flag
}

type Fieldser interface {
	Fields() []string
}

type DefaultTableFieldser interface {
	DefaultTableFields() []string
}

func runPipeCommand() {
	switch GC.Command.(type) {
	case StreamPipeCommander:
		handleStreamPipeCommand()
	case PipeCommander:
		handlePipeCommands()
	default:
	}
}

func handlePipeCommand(text string) {
	defer GC.wgExecute.Done()
	GC.GlobalOptions.logger.Info("Running HandlePipe...")
	item, err := GC.Command.(PipeCommander).HandlePipe(text)
	switch err {
	case nil:
		GC.GlobalOptions.logger.Info("Running Execute...")
		GC.Command.Execute(item, GC.ExecuteResults)
	default:
		GC.ExecuteResults <- err
	}
}

func handlePipeCommands() {
	switch util.Contains(GC.Command.(PipeCommander).PipeFieldOptions(), GC.CommandContext.String("stdin")) {
	case true:
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			GC.wgExecute.Add(1)
			text := scanner.Text()
			go handlePipeCommand(text)
		}
		if scanner.Err() != nil {
			GC.ExecuteResults <- scanner.Err()
		}
	default:
		GC.ExecuteResults <- fmt.Errorf("Unknown STDIN field: %s\n", GC.CommandContext.String("stdin"))
	}
}

func handleStreamPipeCommand() {
	switch util.Contains(GC.Command.(StreamPipeCommander).StreamFieldOptions(), GC.CommandContext.String("stdin")) {
	case true:
		GC.GlobalOptions.logger.Info("Running HandleStreamPipe...")
		stream, err := GC.Command.(StreamPipeCommander).HandleStreamPipe(os.Stdin)
		switch err {
		case nil:
			GC.wgExecute.Add(1)
			go func() {
				defer GC.wgExecute.Done()
				GC.Command.Execute(stream, GC.ExecuteResults)
			}()
		default:
			GC.ExecuteResults <- err
		}
	default:
		handlePipeCommands()
	}
}

func runSingleCommand() {
	switch GC.Command.(type) {
	case PipeCommander, StreamPipeCommander:
		GC.GlobalOptions.logger.Info("Running HandleSingle...")
		item, err := GC.Command.(PipeCommander).HandleSingle()
		switch err {
		case nil:
			GC.GlobalOptions.logger.Info("Running Execute...")
			GC.Command.Execute(item, GC.ExecuteResults)
		default:
			GC.ExecuteResults <- err
		}
	default:
		GC.Command.Execute(nil, GC.ExecuteResults)
	}
}

func handleProgress() {
	p := GC.Command.(Progresser)
	go p.InitProgress()
	for item := range GC.ExecuteResults {
		item := item
		GC.wgProgress.Add(1)
		go func() {
			defer GC.wgProgress.Done()
			switch e := item.(type) {
			case error:
				GC.DoneChan <- e
			default:
				GC.GlobalOptions.logger.Infof("running WaitFor for item: %v", item)
				go p.WaitFor(item)
				id := p.BarID(item)
				p.ShowBar(id)
			}
			GC.GlobalOptions.logger.Info("done waiting on item: %v", item)
		}()
	}

	go func() {
		GC.wgProgress.Wait()
		GC.GlobalOptions.logger.Infoln("closing GC.DoneChan...")
		close(GC.ProgressDoneChan)
	}()

	progressResults := make([]interface{}, 0)

	GC.Logger.Info("Waiting for items on GC.ProgressDoneChan...")
	for r := range GC.ProgressDoneChan {
		progressResults = append(progressResults, r)
	}

	for _, r := range progressResults {
		GC.ResultsRunCommand <- r
	}
}

func handleWait() {
	for item := range GC.ExecuteResults {
		item := item
		GC.wgProgress.Add(1)
		go func() {
			defer GC.wgProgress.Done()
			switch e := item.(type) {
			case error:
				GC.DoneChan <- e
			default:
				GC.GlobalOptions.logger.Infof("running WaitFor for item: %v", item)
				GC.Command.(Waiter).WaitFor(item)
			}
		}()
	}

	go func() {
		GC.wgProgress.Wait()
		GC.GlobalOptions.logger.Infoln("closing GC.DoneChan...")
		close(GC.DoneChan)
	}()

	waitResults := make([]interface{}, 0)

	GC.Logger.Info("Waiting for items on GC.DoneChan...")
	for r := range GC.DoneChan {
		waitResults = append(waitResults, r)
	}

	for _, r := range waitResults {
		GC.ResultsRunCommand <- r
	}
}

func handleQuietNoWait() {
	for r := range GC.ExecuteResults {
		GC.ResultsRunCommand <- r
	}
}

func RunCommand() {
	defer close(GC.ResultsRunCommand)
	switch GC.CommandContext.IsSet("stdin") {
	case true:
		GC.GlobalOptions.logger.Info("Running runPipeCommand...")
		runPipeCommand()
	default:
		GC.GlobalOptions.logger.Info("Running runSingleCommand...")
		GC.wgExecute.Add(1)
		go func() {
			defer GC.wgExecute.Done()
			runSingleCommand()
		}()
	}

	go func() {
		GC.wgExecute.Wait()
		close(GC.ExecuteResults)
	}()

	if _, ok := GC.Command.(Progresser); ok && !GC.CommandContext.IsSet("quiet") {
		handleProgress()
	} else if GC.CommandContext.IsSet("wait") {
		handleWait()
	} else {
		handleQuietNoWait()
	}
}
