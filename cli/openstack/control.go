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

	lib.Log.SetLevel(gopts.loglevel)

	ao := &authopts{
		cmd:     cmd,
		region:  gopts.region,
		gao:     gopts.authOptions,
		nocache: gopts.nocache,
		urltype: gopts.urlType,
	}
	sc, err := auth(ao)
	if err != nil {
		return ErrExit1{err}
	}

	if !gopts.nocache {
		defer cachecreds(ao, sc)
	}

	cmd.SetServiceClient(sc)

	if t, ok := cmd.(interfaces.Tabler); ok {
		lib.Log.Debugln("cmd implements Tabler")
		t.SetTable(cmd.Context().IsSet("table"))
		lib.Log.Debugln("p.ShouldTable() : ", t.ShouldTable())
		t.SetHeader(cmd.Context().IsSet("no-header"))
		lib.Log.Debugln("p.ShouldHeader() : ", t.ShouldHeader())
	}

	if f, ok := cmd.(interfaces.Fieldser); ok {
		lib.Log.Debugln("cmd implements Fieldser")
		if len(f.Fields()) > 0 {
			f.SetFields(strings.Split(cmd.Context().String("fields"), ","))
		}
	}

	if p, ok := cmd.(interfaces.Progresser); ok {
		lib.Log.Debugln("cmd implements Progresser")
		p.SetProgress(cmd.Context().IsSet("quiet"))
		lib.Log.Debugln("p.ShouldProgress() : ", p.ShouldProgress())
	}

	if w, ok := cmd.(interfaces.Waiter); ok {
		lib.Log.Debugln("cmd implements Waiter")
		w.SetWait(cmd.Context().IsSet("wait"))
		lib.Log.Debugln("w.ShouldWait() : ", w.ShouldWait())
	}

	if flagser, ok := cmd.(interfaces.Flagser); ok {
		err = flagser.HandleFlags()
		if err != nil {
			return ErrExit1{err}
		}
	}

	progch := make(chan interface{})
	outch := make(chan interface{})

	go exec(cmd, progch)
	go prog(cmd, progch, outch)
	err = outres(cmd, outch)

	if err != nil {
		return ErrExit1{err}
	}

	return nil
}

func prog(cmd interfaces.Commander, inch, outch chan interface{}) {
	if p, okp := cmd.(interfaces.Progresser); okp {
		holdch := make(chan interface{})
		wg := new(sync.WaitGroup)
		var once sync.Once

		for res := range inch {
			switch t := res.(type) {
			case error:
				holdch <- t
			default:
				if pi, okpi := t.(interfaces.ProgressItemer); okpi {
					id := pi.ID()
					if w, okw := cmd.(interfaces.Waiter); okw {
						go w.WaitFor(pi, nil)
					}
					var b interfaces.ProgressBarrer
					if p.ShouldProgress() {
						once.Do(func() {
							lib.Log.Devln("initializing progress...")
							p.InitProgress()
							lib.Log.Devln("adding summary bar...")
							p.AddSummaryBar()
						})
						lib.Log.Devln("creating bar...")
						b = p.CreateBar(pi)
						lib.Log.Devln("created bar. starting bar...")
						p.StartBar()
					}

					wg.Add(1)
					go func() {
						defer wg.Done()
						for {
							select {
							case up := <-pi.UpCh():
								if p.ShouldProgress() {
									s := new(traits.ProgressStatusUpdate)
									s.SetBarID(id)
									s.SetChange(up)
									b.Update(s)
								}
							case res := <-pi.EndCh():
								switch t := res.(type) {
								case error:
									if p.ShouldProgress() {
										s := new(traits.ProgressStatusError)
										s.SetBarID(id)
										s.SetErr(t)
										//b.Error(s)
										p.ErrorBar()
									}
									holdch <- t
									return
								default:
									if p.ShouldProgress() {
										s := new(traits.ProgressStatusComplete)
										s.SetBarID(id)
										b.Complete(s)
										p.CompleteBar()
									}
									holdch <- t
									return
								}
							}
						}
					}()

				}
			}
		}

		go func() {
			defer close(holdch)
			wg.Wait()
		}()

		go func() {
			defer close(outch)
			for r := range holdch {
				outch <- r
			}
		}()

		return
	}

	defer close(outch)
	for r := range inch {
		outch <- r
	}
}
