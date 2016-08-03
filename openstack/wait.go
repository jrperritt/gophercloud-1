package openstack

import (
	"sync"

	"github.com/gophercloud/cli/lib"
)

func ExecuteAndWait(c lib.Commander, in, out chan interface{}) {
	defer close(out)
	var once sync.Once
	var wg sync.WaitGroup
	ch1 := make(chan interface{})

	ch3 := make(chan interface{})

	go c.Execute(in, ch1)

	go func() {
		for item := range ch1 {
			ch2 := make(chan interface{})
			item := item
			go func() {
				ch2 <- item
				close(ch2)
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				if progresser, ok := c.(lib.Progresser); ok {
					once.Do(progresser.InitProgress)
					progresser.ShowProgress(ch2, ch3)
				} else {
					for item := range ch2 {
						ch3 <- item
					}
				}
			}()
		}
		go func() {
			wg.Wait()
			close(ch3)
		}()
	}()

	msgs := make([]interface{}, 0)

	for raw := range ch3 {
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
