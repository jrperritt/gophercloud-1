package openstack

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
)

func exec(cmd interfaces.Commander, out chan interface{}) {
	wg := new(sync.WaitGroup)
	if stdin := cmd.Context().String("stdin"); stdin != "" {
		if p, ok := cmd.(interfaces.PipeCommander); ok && util.Contains(p.PipeFieldOptions(), stdin) {
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
					handleitem(cmd, out, item)
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
			case interfaces.PipeCommander:
				item, err := pc.HandleSingle()
				if err != nil {
					out <- err
					return
				}
				handleitem(cmd, out, item)
			default:
				cmd.Execute(nil, out)
			}
		}()
	}

	go func() {
		wg.Wait()
		if proger, ok := cmd.(interfaces.Progresser); ok {
			close(proger.ProgStartCh())
		}
		close(out)
	}()
}

func handleitem(cmd interfaces.Commander, out chan interface{}, item interface{}) {
	if proger, ok := cmd.(interfaces.Progresser); ok {
		if pi, ok := item.(interfaces.ProgressItemer); ok {
			proger.ProgStartCh() <- pi
		}
	}
	cmd.Execute(item, out)
}
