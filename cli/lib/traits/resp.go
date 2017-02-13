package traits

import cli "gopkg.in/urfave/cli.v1"

type Fieldsable struct {
	fields []string
}

func (c *Fieldsable) SetFields(f []string) {
	c.fields = f
}

func (c *Fieldsable) Fields() []string {
	return c.fields
}

func (c *Fieldsable) FieldsFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "fields",
			Usage: "[optional] Only return these comma-separated case-insensitive fields.", //+
			//fmt.Sprintf("\n\tChoices: %s", strings.Join(c.fields, ", ")),
		},
	}
}
