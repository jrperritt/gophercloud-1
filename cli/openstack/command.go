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

func exec(cmd interfaces.Commander, out chan interface{}) {
	wg := new(sync.WaitGroup)

	if stdin := cmd.Context().String("stdin"); stdin != "" {
		if sp, ok := cmd.(interfaces.StreamCommander); ok && util.Contains(sp.StreamFieldOptions(), stdin) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				item, err := sp.HandleStream()
				if err != nil {
					out <- err
					return
				}
				out <- item
				cmd.Execute(item, out)
			}()
		} else if p, ok := cmd.(interfaces.PipeCommander); ok && util.Contains(p.PipeFieldOptions(), stdin) {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				wg.Add(1)
				text := scanner.Text()
				go func() {
					defer wg.Done()
					item, err := p.HandlePipe(text)
					if err != nil {
						out <- err
						return
					}
					out <- item
					cmd.Execute(item, out)
				}()
			}
			if scanner.Err() != nil {
				out <- scanner.Err()
			}
		} else {
			out <- fmt.Errorf("Unknown STDIN field: %s\n", stdin)
		}
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch pc := cmd.(type) {
			case interfaces.StreamCommander:
				item, err := pc.HandleSingle()
				if err != nil {
					lib.Log.Debugf("Error from HandleSingle: %s", err)
					out <- err
					return
				}
				out <- item
				cmd.Execute(item, out)
			case interfaces.PipeCommander:
				item, err := pc.HandleSingle()
				if err != nil {
					lib.Log.Debugf("Error from HandleSingle: %s", err)
					out <- err
					return
				}
				out <- item
				cmd.Execute(item, out)
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
