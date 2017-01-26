package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"
	"text/tabwriter"

	"gopkg.in/urfave/cli.v1"
)

var appHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} <command> <subcommand> <action> [FLAGS]
   {{if .Version}}
VERSION:
   {{.Version}}
   {{end}}{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
	 {{end}}{{end}}{{if .Flags}}
FLAGS:
   {{range .Flags}}{{flag .}}
   {{end}}{{end}}
`

var commandHelpTemplate = `NAME: {{.Name}} - {{.Usage}}{{if .Description}}
DESCRIPTION: {{.Description}}{{end}}{{if .Flags}}
COMMAND FLAGS:
   {{range .Flags}}{{flag .}}
   {{end}}{{end}}
`

var subcommandHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}}{{if eq (len (split .Name " ")) 2}} <subcommand>{{end}} <action> [FLAGS]
{{if eq (len (split .Name " ")) 2}}SUBCOMMANDS{{else}}ACTIONS{{end}}:
   {{range .Commands}}{{join .Names ", "}}
   {{end}}
`

func printHelp(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"split": strings.Split,
		"join":  strings.Join,
		"flag":  flag,
	}

	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func flag(cliflag cli.Flag) string {
	if cliflag.GetName() == "generate-bash-completion" {
		return ""
	}
	var flagString string
	switch flagType := cliflag.(type) {
	case cli.StringFlag:
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), flagType.Usage)
	case cli.IntFlag:
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), flagType.Usage)
	case cli.BoolFlag:
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), flagType.Usage)
	case cli.StringSliceFlag:
		flagString = fmt.Sprintf("%s\t%s", fmt.Sprintf("--%s", flagType.Name), flagType.Usage)
	}
	return flagString
}

/*
func commandInfo(c cli.Command) string {
	if c, ok := interface{}(c).(lib.CommandInfoer); ok {
		return c.CommandInfo()
	}
	return ""
}
*/
