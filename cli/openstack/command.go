package openstack

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
)

func handlePipeCommands(cmd interfaces.Commander, out chan interface{}) {
	p := cmd.(interfaces.PipeCommander)
	switch util.Contains(p.PipeFieldOptions(), cmd.Context().String("stdin")) {
	case true:
		scanner := bufio.NewScanner(os.Stdin)
		wg := new(sync.WaitGroup)
		for scanner.Scan() {
			wg.Add(1)
			text := scanner.Text()
			go func() {
				p := cmd.(interfaces.PipeCommander)
				defer wg.Done()
				lib.Log.Debugln("Running HandlePipe...")
				item, err := p.HandlePipe(text)
				switch err {
				case nil:
					lib.Log.Debugln("Running Execute...")
					p.Execute(item, out)
				default:
					out <- err
				}
			}()
		}
		if scanner.Err() != nil {
			out <- scanner.Err()
		}
	default:
		out <- fmt.Errorf("Unknown STDIN field: %s\n", cmd.Context().String("stdin"))
	}
}

func exec(cmd interfaces.Commander, out chan interface{}) {
	wg := new(sync.WaitGroup)
	if cmd.Context().IsSet("stdin") {
		lib.Log.Debugln("Running runPipeCommand...")
		switch cmd.(type) {
		case interfaces.StreamPipeCommander:
			p := cmd.(interfaces.StreamPipeCommander)
			if util.Contains(p.StreamFieldOptions(), cmd.Context().String("stdin")) {
				lib.Log.Debugln("Running HandleStreamPipe...")
				stream, err := p.HandleStreamPipe(os.Stdin)
				if err != nil {
					out <- err
					return
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					cmd.Execute(stream, out)
				}()
			} else {
				handlePipeCommands(cmd, out)
			}
		case interfaces.PipeCommander:
			handlePipeCommands(cmd, out)
		default:
		}
	} else {
		lib.Log.Debugln("Running runSingleCommand...")
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch cmd.(type) {
			case interfaces.PipeCommander, interfaces.StreamPipeCommander:
				lib.Log.Debugln("Running HandleSingle...")
				item, err := cmd.(interfaces.PipeCommander).HandleSingle()
				switch err {
				case nil:
					lib.Log.Debugln("Running Execute...")
					cmd.Execute(item, out)
				default:
					out <- err
				}
			default:
				cmd.Execute(nil, out)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}
