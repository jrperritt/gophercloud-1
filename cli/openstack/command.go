package openstack

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud/cli/lib"
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
	p := cmd.(interfaces.PipeCommander)
	defer p.WG().Done()
	lib.Log.Debugln("Running HandlePipe...")
	item, err := p.HandlePipe(text)
	switch err {
	case nil:
		lib.Log.Debugln("Running Execute...")
		p.Execute(item, cmd.ExecDoneCh())
	default:
		cmd.ExecDoneCh() <- err
	}
}

func handlePipeCommands(cmd interfaces.Commander) {
	p := cmd.(interfaces.PipeCommander)
	switch util.Contains(p.PipeFieldOptions(), cmd.Context().String("stdin")) {
	case true:
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			p.WG().Add(1)
			text := scanner.Text()
			go handlePipeCommand(cmd, text)
		}
		if scanner.Err() != nil {
			cmd.ExecDoneCh() <- scanner.Err()
		}
	default:
		cmd.ExecDoneCh() <- fmt.Errorf("Unknown STDIN field: %s\n", cmd.Context().String("stdin"))
	}
}

func handleStreamPipeCommand(cmd interfaces.Commander) {
	p := cmd.(interfaces.StreamPipeCommander)
	switch util.Contains(p.StreamFieldOptions(), cmd.Context().String("stdin")) {
	case true:
		lib.Log.Debugln("Running HandleStreamPipe...")
		stream, err := p.HandleStreamPipe(os.Stdin)
		switch err {
		case nil:
			p.WG().Add(1)
			go func() {
				defer p.WG().Done()
				cmd.Execute(stream, cmd.ExecDoneCh())
			}()
		default:
			cmd.ExecDoneCh() <- err
		}
	default:
		handlePipeCommands(cmd)
	}
}

func runSingleCommand(cmd interfaces.Commander) {
	switch cmd.(type) {
	case interfaces.PipeCommander, interfaces.StreamPipeCommander:
		lib.Log.Debugln("Running HandleSingle...")
		item, err := cmd.(interfaces.PipeCommander).HandleSingle()
		switch err {
		case nil:
			lib.Log.Debugln("Running Execute...")
			cmd.Execute(item, cmd.ExecDoneCh())
		default:
			cmd.ExecDoneCh() <- err
		}
	default:
		cmd.Execute(nil, cmd.ExecDoneCh())
	}
}

func handleProgress(cmd interfaces.Commander) {
	p := cmd.(interfaces.Progresser)
	donech := make(chan interface{})
	go p.InitProgress(donech)
	for item := range cmd.ExecDoneCh() {
		item := item
		p.WG().Add(1)
		go func() {
			defer p.WG().Done()
			switch e := item.(type) {
			case error:
				cmd.AllDoneCh() <- e
			default:
				if w, ok := p.(interfaces.Waiter); ok {
					lib.Log.Debugf("running Waiter.WaitFor for item: %v", item)
					go w.WaitFor(item)
				}
				id := p.BarID(item)
				p.ShowBar(id)
			}
			lib.Log.Debugf("done waiting on item: %v", item)
		}()
	}

	go func() {
		p.WG().Wait()
		lib.Log.Debugln("closing Progresser.ProgDoneChOut()...")
		close(p.ProgDoneChOut())
	}()

	progressResults := make([]interface{}, 0)

	lib.Log.Debugln("Waiting for items on Progresser.ProgDoneChOut()...")
	for r := range p.ProgDoneChOut() {
		progressResults = append(progressResults, r)
	}

	for _, r := range progressResults {
		cmd.AllDoneCh() <- r
	}
}

func handleWait(cmd interfaces.Commander) {
	w := cmd.(interfaces.Waiter)
	for item := range cmd.ExecDoneCh() {
		item := item
		w.WG().Add(1)
		go func() {
			defer w.WG().Done()
			switch e := item.(type) {
			case error:
				cmd.AllDoneCh() <- e
			default:
				lib.Log.Debugf("running WaitFor for item: %v", item)
				w.WaitFor(item)
			}
		}()
	}

	go func() {
		w.WG().Wait()
		lib.Log.Debugln("closing w.WaitDoneCh()...")
		close(w.WaitDoneCh())
	}()

	waitResults := make([]interface{}, 0)

	lib.Log.Debugln("Waiting for items on w.WaitDoneCh()...")
	for r := range w.WaitDoneCh() {
		waitResults = append(waitResults, r)
	}

	for _, r := range waitResults {
		cmd.AllDoneCh() <- r
	}
}

func handleQuietNoWait(cmd interfaces.Commander) {
	for r := range cmd.ExecDoneCh() {
		cmd.AllDoneCh() <- r
	}
}

func runcmd(cmd interfaces.Commander) {
	defer close(cmd.AllDoneCh())
	switch cmd.Context().IsSet("stdin") {
	case true:
		lib.Log.Debugln("Running runPipeCommand...")
		runPipeCommand(cmd)
	default:
		lib.Log.Debugln("Running runSingleCommand...")
		cmd.WG().Add(1)
		go func() {
			defer cmd.WG().Done()
			runSingleCommand(cmd)
		}()
	}

	go func() {
		cmd.WG().Wait()
		close(cmd.ExecDoneCh())
	}()

	if p, ok := cmd.(interfaces.Progresser); ok && p.ShouldProgress() {
		handleProgress(cmd)
	} else if w, ok := cmd.(interfaces.Waiter); ok && w.ShouldWait() {
		handleWait(cmd)
	} else {
		handleQuietNoWait(cmd)
	}
}
