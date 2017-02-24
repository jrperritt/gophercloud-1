package traits

import cli "gopkg.in/urfave/cli.v1"

// Tableable is a trait the should be embedded in command types that offer tabular output
type Tableable struct {
	//Fieldsable
	table, noheader bool
}

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

func (c *Tableable) SetTable(b bool) {
	c.table = b
}

// ShouldTable returns whether or not to output in tablular format.
// Partially satisfies interfaces.Tabler interface
func (c *Tableable) ShouldTable() bool {
	return c.table
}

func (c *Tableable) SetHeader(b bool) {
	c.noheader = b
}

// ShouldHeader returns whether or not to output the header for tablular format.
// Partially satisfies interfaces.Tabler interface
func (c *Tableable) ShouldHeader() bool {
	return c.table && !c.noheader
}
