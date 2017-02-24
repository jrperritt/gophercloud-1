package traits

import (
	"runtime"

	cli "gopkg.in/urfave/cli.v1"
)

type Pipeable struct {
	concurrency int
}

func (c *Pipeable) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *Pipeable) PipeFieldOptions() []string {
	return []string{"id"}
}

func (c *Pipeable) PipeFlags() []cli.Flag {
	return []cli.Flag{
		cli.IntFlag{
			Name:  "concurrency",
			Usage: "The number of concurrent operations to perform. Defaults to the number of CPUs",
		},
	}
}

func (c *Pipeable) SetConcurrency(i int) {
	if i == 0 {
		i = runtime.NumCPU()
	}
	c.concurrency = i
}

func (c *Pipeable) Concurrency() int {
	return c.concurrency
}
