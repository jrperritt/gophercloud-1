package openstack

import (
	"sync"

	"github.com/gophercloud/cli/lib"
)

func ExecuteAndWait(waiter lib.Waiter, in, out chan interface{}) {
	defer close(out)

	var wg sync.WaitGroup

	chExec := make(chan interface{})
	chRes := make(chan interface{})

	go waiter.Execute(in, chExec)

	go func() {
		progresser, ok := waiter.(lib.Progresser)
		switch ok && progresser.ShouldProgress() {
		case true:
			progresser.InitProgress()
			for item := range chExec {
				item := item
				switch e := item.(type) {
				case error:
					chRes <- e
					continue
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					progresser.ShowProgress(item, chRes)
				}()
			}
			go func() {
				wg.Wait()
				close(chRes)
			}()
		case false:
			for item := range chExec {
				chRes <- item
			}
			close(chRes)
		}
	}()

	msgs := make([]interface{}, 0)

	for raw := range chRes {
		switch msg := raw.(type) {
		case error:
			out <- msg
		default:
			msgs = append(msgs, msg)
		}
	}

	for _, msg := range msgs {
		out <- msg
	}
}
