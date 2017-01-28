package traits

import cli "gopkg.in/urfave/cli.v1"

// Tableable is a trait the should be embedded in command types that offer tabular output
type Tableable struct{}

// TableFlags returns flags for commands that offer tabular output
// Partially satisfies interfaces.Tabler interface
func (c *Tableable) TableFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "table",
			Usage: "[optional] If provided, output will be in tabular format.",
		},
		cli.BoolFlag{
			Name:  "no-header",
			Usage: "[optional] Do not return a header for tabular output.",
		},
	}
}
