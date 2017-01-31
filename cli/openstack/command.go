package openstack

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
)

func runPipeCommand(cmd interfaces.Commander) {
	switch cmd.(type) {
	case interfaces.StreamPipeCommander:
		handleStreamPipeCommand(cmd)
	case interfaces.PipeCommander:
		handlePipeCommands(cmd)
	default:
	}
}

func handlePipeCommand(cmd interfaces.Commander, text string) {
	defer gctx.wgExecute.Done()
	gctx.Logger.Debugln("Running HandlePipe...")
	item, err := cmd.(interfaces.PipeCommander).HandlePipe(text)
	switch err {
	case nil:
		gctx.Logger.Debugln("Running Execute...")
		cmd.Execute(item, gctx.ExecuteResults)
	default:
		gctx.ExecuteResults <- err
	}
}

func handlePipeCommands(cmd interfaces.Commander) {
	switch util.Contains(cmd.(interfaces.PipeCommander).PipeFieldOptions(), cmd.Context().String("stdin")) {
	case true:
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			gctx.wgExecute.Add(1)
			text := scanner.Text()
			go handlePipeCommand(cmd, text)
		}
		if scanner.Err() != nil {
			gctx.ExecuteResults <- scanner.Err()
		}
	default:
		gctx.ExecuteResults <- fmt.Errorf("Unknown STDIN field: %s\n", cmd.Context().String("stdin"))
	}
}

func handleStreamPipeCommand(cmd interfaces.Commander) {
	switch util.Contains(cmd.(interfaces.StreamPipeCommander).StreamFieldOptions(), cmd.Context().String("stdin")) {
	case true:
		gctx.Logger.Debugln("Running HandleStreamPipe...")
		stream, err := cmd.(interfaces.StreamPipeCommander).HandleStreamPipe(os.Stdin)
		switch err {
		case nil:
			gctx.wgExecute.Add(1)
			go func() {
				defer gctx.wgExecute.Done()
				cmd.Execute(stream, gctx.ExecuteResults)
			}()
		default:
			gctx.ExecuteResults <- err
		}
	default:
		handlePipeCommands(cmd)
	}
}

func runSingleCommand(cmd interfaces.Commander) {
	switch cmd.(type) {
	case interfaces.PipeCommander, interfaces.StreamPipeCommander:
		gctx.Logger.Debugln("Running HandleSingle...")
		item, err := cmd.(interfaces.PipeCommander).HandleSingle()
		switch err {
		case nil:
			gctx.Logger.Debugln("Running Execute...")
			cmd.Execute(item, gctx.ExecuteResults)
		default:
			gctx.ExecuteResults <- err
		}
	default:
		cmd.Execute(nil, gctx.ExecuteResults)
	}
}

func handleProgress(cmd interfaces.Commander) {
	p := cmd.(interfaces.Progresser)
	go p.InitProgress()
	for item := range gctx.ExecuteResults {
		item := item
		gctx.wgProgress.Add(1)
		go func() {
			defer gctx.wgProgress.Done()
			switch e := item.(type) {
			case error:
				gctx.DoneChan <- e
			default:
				gctx.Logger.Debugf("running WaitFor for item: %v", item)
				if w, ok := p.(interfaces.Waiter); ok {
					go w.WaitFor(item)
				}

				id := p.BarID(item)
				p.ShowBar(id)
			}
			gctx.Logger.Debugf("done waiting on item: %v", item)
		}()
	}

	go func() {
		gctx.wgProgress.Wait()
		gctx.Logger.Debugln("closing gctx.DoneChan...")
		close(gctx.ProgressDoneChan)
	}()

	progressResults := make([]interface{}, 0)

	gctx.Logger.Debugln("Waiting for items on gctx.ProgressDoneChan...")
	for r := range gctx.ProgressDoneChan {
		progressResults = append(progressResults, r)
	}

	for _, r := range progressResults {
		gctx.ResultsRunCommand <- r
	}
}

func handleWait(cmd interfaces.Commander) {
	for item := range gctx.ExecuteResults {
		item := item
		gctx.wgProgress.Add(1)
		go func() {
			defer gctx.wgProgress.Done()
			switch e := item.(type) {
			case error:
				gctx.DoneChan <- e
			default:
				gctx.Logger.Debugf("running WaitFor for item: %v", item)
				cmd.(interfaces.Waiter).WaitFor(item)
			}
		}()
	}

	go func() {
		gctx.wgProgress.Wait()
		gctx.Logger.Debugln("closing gctx.DoneChan...")
		close(gctx.DoneChan)
	}()

	waitResults := make([]interface{}, 0)

	gctx.Logger.Debugln("Waiting for items on gctx.DoneChan...")
	for r := range gctx.DoneChan {
		waitResults = append(waitResults, r)
	}

	for _, r := range waitResults {
		gctx.ResultsRunCommand <- r
	}
}

func handleQuietNoWait(_ interfaces.Commander) {
	for r := range gctx.ExecuteResults {
		gctx.ResultsRunCommand <- r
	}
}

func runcmd(cmd interfaces.Commander) {
	defer close(gctx.ResultsRunCommand)
	switch cmd.Context().IsSet("stdin") {
	case true:
		gctx.Logger.Debugln("Running runPipeCommand...")
		runPipeCommand(cmd)
	default:
		gctx.Logger.Debugln("Running runSingleCommand...")
		gctx.wgExecute.Add(1)
		go func() {
			defer gctx.wgExecute.Done()
			runSingleCommand(cmd)
		}()
	}

	go func() {
		gctx.wgExecute.Wait()
		close(gctx.ExecuteResults)
	}()

	if p, ok := cmd.(interfaces.Progresser); ok && p.ShouldProgress() {
		handleProgress(cmd)
	} else if w, ok := cmd.(interfaces.Waiter); ok && w.ShouldWait() {
		handleWait(cmd)
	} else {
		handleQuietNoWait(cmd)
	}
}
