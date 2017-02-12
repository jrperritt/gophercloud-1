package openstack

import (
	"os"
	"strings"
	"sync"

	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/lib/traits"
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

	if t, ok := cmd.(interfaces.Tabler); ok {
		lib.Log.Debugln("cmd implements Tabler")
		t.SetTable(cmd.Context().IsSet("table"))
		t.SetHeader(cmd.Context().IsSet("no-header"))
	}

	if f, ok := cmd.(interfaces.Fieldser); ok {
		lib.Log.Debugln("cmd implements Fieldser")
		f.SetFields(strings.Split(cmd.Context().String("fields"), ","))
	}

	if p, ok := cmd.(interfaces.Progresser); ok {
		lib.Log.Debugln("cmd implements Progresser")
		p.SetProgress(cmd.Context().IsSet("quiet"))
		lib.Log.Debugln("p.ShouldProgress() : ", p.ShouldProgress())
	}

	err = cmd.HandleFlags()
	if err != nil {
		return ErrExit1{err}
	}

	execchout := make(chan interface{})
	if p, ok := cmd.(interfaces.Progresser); ok && p.ShouldProgress() {
		progchout := make(chan interface{})
		p.InitProgress()
		exec(cmd, execchout)
		go prog(p, progchout)
		err = outres(cmd, progchout)
	} else {
		go exec(cmd, execchout)
		err = outres(cmd, execchout)
	}

	if err != nil {
		return ErrExit1{err}
	}

	return nil
}

func prog(p interfaces.Progresser, outch chan interface{}) {
	defer close(outch)
	waitch := make(chan interface{})
	wg := new(sync.WaitGroup)

	for pi := range p.ProgStartCh() {
		pi := pi
		id := pi.ID()
		b := p.CreateBar(pi)
		p.StartBar()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for up := range pi.UpCh() {
				s := new(traits.ProgressStatusUpdate)
				s.SetBarID(id)
				s.SetChange(up)
				b.Update(s)
			}

			switch t := (<-pi.EndCh()).(type) {
			case error:
				s := new(traits.ProgressStatusError)
				s.SetBarID(id)
				s.SetErr(t)
				//b.Error(s)
				p.ErrorBar()
				waitch <- t
			default:
				s := new(traits.ProgressStatusComplete)
				s.SetBarID(id)
				b.Complete(s)
				p.CompleteBar()
				waitch <- t
			}
		}()
	}

	go func() {
		wg.Wait()
		close(waitch)
	}()

	progressResults := make([]interface{}, 0)

	for r := range waitch {
		progressResults = append(progressResults, r)
	}

	for _, r := range progressResults {
		outch <- r
	}
}
