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

	if w, ok := cmd.(interfaces.Waiter); ok {
		lib.Log.Debugln("cmd implements Waiter")
		w.SetWait(cmd.Context().IsSet("wait"))
		lib.Log.Debugln("w.ShouldWait() : ", w.ShouldWait())
	}

	err = cmd.HandleFlags()
	if err != nil {
		return ErrExit1{err}
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

func prog(cmd interfaces.Commander, in, out chan interface{}) {
	if p, ok := cmd.(interfaces.Progresser); ok {
		waitch := make(chan interface{})
		wg := new(sync.WaitGroup)
		var once sync.Once

		wg.Add(1)
		go func() {
			defer wg.Done()
			for res := range in {
				switch t := res.(type) {
				case error:
					waitch <- t
				default:
					if pi, ok := t.(interfaces.ProgressItemer); ok {
						id := pi.ID()
						var b interfaces.ProgressBarrer
						if p.ShouldProgress() {
							once.Do(func() {
								p.InitProgress()
								p.AddSummaryBar()
							})

							if p.ShouldProgress() {
								b = p.CreateBar(pi)
								p.StartBar()
							}
						}

						wg.Add(1)
						go func() {
							defer wg.Done()
							for {
								select {
								case up := <-pi.UpCh():
									s := new(traits.ProgressStatusUpdate)
									s.SetBarID(id)
									s.SetChange(up)
									if p.ShouldProgress() {
										b.Update(s)
									}
								case res := <-pi.EndCh():
									switch t := res.(type) {
									case error:
										s := new(traits.ProgressStatusError)
										s.SetBarID(id)
										s.SetErr(t)
										if p.ShouldProgress() {
											//b.Error(s)
											p.ErrorBar()
										}
										waitch <- t
										return
									default:
										s := new(traits.ProgressStatusComplete)
										s.SetBarID(id)
										if p.ShouldProgress() {
											b.Complete(s)
											p.CompleteBar()
										}
										waitch <- t
										return
									}
								}
							}
						}()

					}
				}
			}
		}()

		go func() {
			wg.Wait()
			close(waitch)
		}()

		go func() {
			defer close(out)
			for r := range waitch {
				out <- r
			}
		}()

		return
	}

	defer close(out)
	for r := range in {
		out <- r
	}
}
