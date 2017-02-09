package openstack

import (
	"os"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"gopkg.in/urfave/cli.v1"
)

// Action is the common function all commands run
func Action(ctx *cli.Context, cmd interfaces.Commander) error {
	lib.InitLog(os.Stdout)
	/*
		dir, err := util.RackDir()
		if err != nil {
			lib.Log.
		}
		logout, err := os.Open()
	*/

	cmd.SetContext(ctx)

	var err error

	gopts, err := globalopts(ctx)
	if err != nil {
		return ErrExit1{err}
	}

	lib.Log.SetDebug(gopts.debug)

	ao := &authopts{
		cmd:     cmd,
		region:  gopts.region,
		gao:     gopts.authOptions,
		nocache: gopts.noCache,
		urltype: gopts.urlType,
	}
	sc, err := auth(ao)
	if err != nil {
		return ErrExit1{err}
	}

	if !gopts.noCache {
		defer cachecreds(ao, sc)
	}

	cmd.SetServiceClient(sc)

	lib.Log.Debugln("Running HandleInterfaceFlags...")
	err = interfaces.HandleInterfaceFlags(cmd)
	if err != nil {
		return ErrExit1{err}
	}

	lib.Log.Debugln("Running HandleFlags...")
	err = cmd.HandleFlags()
	if err != nil {
		return ErrExit1{err}
	}

	lib.Log.Debugln("Running RunCommand...")
	execchout := make(chan interface{})
	go exec(cmd, execchout)

	if w, ok := cmd.(interfaces.Waiter); ok && w.ShouldWait() {
		waitchout := make(chan interface{})
		lib.Log.Debugln("going to wait")
		go wait(w, execchout, waitchout)
		if p, ok := w.(interfaces.Progresser); ok && p.ShouldProgress() {
			// wait and prog
			progchout := make(chan interface{})
			lib.Log.Debugln("going to prog")
			go prog(p, waitchout, progchout)
			lib.Log.Debugln("outres from prog")
			err = outres(cmd, progchout)
		} else {
			// wait and no prog
			lib.Log.Debugln("outres from wait")
			err = outres(cmd, waitchout)
		}
	} else if p, ok := cmd.(interfaces.Progresser); ok && p.ShouldProgress() {
		// no wait and prog
		progchout := make(chan interface{})
		lib.Log.Debugln("going to prog")
		go prog(p, execchout, progchout)
		lib.Log.Debugln("outres from prog")
		err = outres(cmd, progchout)
	} else {
		// no wait and no prog
		lib.Log.Debugln("outres from exec")
		err = outres(cmd, execchout)
	}

	if err != nil {
		return ErrExit1{err}
	}

	return nil
}

func wait(w interfaces.Waiter, in, out chan interface{}) {
	defer close(out)
	wg := new(sync.WaitGroup)
	waitchmid := make(chan interface{})
	for item := range in {
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch e := item.(type) {
			case error:
				out <- e
			default:
				lib.Log.Debugf("running WaitFor for item: %v", item)
				w.WaitFor(item, waitchmid)
			}
		}()
	}

	go func() {
		wg.Wait()
		lib.Log.Debugln("closing w.WaitDoneCh()...")
		close(waitchmid)
	}()

	waitResults := make([]interface{}, 0)

	lib.Log.Debugln("Waiting for items on waitchmid...")
	for r := range waitchmid {
		waitResults = append(waitResults, r)
	}

	for _, r := range waitResults {
		out <- r
	}
}

func prog(p interfaces.Progresser, in, out chan interface{}) {
	defer close(out)
	progchmid := make(chan interface{})
	p.InitProgress(progchmid)
	wg := new(sync.WaitGroup)
	for item := range in {
		item := item
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch e := item.(type) {
			case error:
				out <- e
			default:
				id := p.BarID(item)
				lib.Log.Debugln("running p.ShowBar...")
				p.ShowBar(id)
			}
			lib.Log.Debugf("done waiting on item: %v", item)
		}()
	}

	go func() {
		wg.Wait()
		lib.Log.Debugln("closing progchmid...")
		close(progchmid)
	}()

	progressResults := make([]interface{}, 0)

	lib.Log.Debugln("Waiting for items on progchmid...")
	for r := range progchmid {
		progressResults = append(progressResults, r)
	}

	for _, r := range progressResults {
		out <- r
	}
}
