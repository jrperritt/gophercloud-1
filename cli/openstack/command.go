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
		if p, ok := cmd.(interfaces.PipeCommander); ok && util.Contains(p.PipeFieldOptions(), stdin) {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				wg.Add(1)
				text := scanner.Text()
				go func() {
					defer wg.Done()
					item, err := p.HandlePipe(text)
					switch err {
					case nil:
						if proger, ok := cmd.(interfaces.Progresser); ok && proger.ShouldProgress() {
							if pi, ok := item.(interfaces.ProgressItemer); ok {
								lib.Log.Debugf("Sending on ProgStartCh: %+v", pi)
								proger.ProgStartCh() <- pi
							}
						}
						p.Execute(item, out)
					default:
						out <- err
					}
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
				switch err {
				case nil:
					if proger, ok := cmd.(interfaces.Progresser); ok && proger.ShouldProgress() {
						if pi, ok := item.(interfaces.ProgressItemer); ok {
							lib.Log.Debugf("Sending on ProgStartCh: %+v", pi)
							proger.ProgStartCh() <- pi
						}
					}
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
		if proger, ok := cmd.(interfaces.Progresser); ok && proger.ShouldProgress() {
			lib.Log.Debugln("closing proger.ProgStartCh")
			close(proger.ProgStartCh())
		}
		close(out)
	}()
}
