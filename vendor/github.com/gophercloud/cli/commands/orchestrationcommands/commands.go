package orchestrationcommands

import (
	"github.com/gophercloud/cli/commands/orchestrationcommands/buildinfocommands"
	"github.com/gophercloud/cli/commands/orchestrationcommands/stackcommands"
	"github.com/gophercloud/cli/commands/orchestrationcommands/stackeventcommands"
	"github.com/gophercloud/cli/commands/orchestrationcommands/stackresourcecommands"
	"github.com/gophercloud/cli/commands/orchestrationcommands/stacktemplatecommands"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
)

var serviceClientType = "orchestration"

// Get returns all the commands allowed for an `orchestration` request.
func Get() []cli.Command {
	return []cli.Command{
		{
			Name:        "build-info",
			Usage:       "Build information.",
			Subcommands: buildinfocommands.Get(),
		},
		{
			Name:        "stack",
			Usage:       "Stack management.",
			Subcommands: stackcommands.Get(),
		},
		{
			Name:        "event",
			Usage:       "Stack event queries.",
			Subcommands: stackeventcommands.Get(),
		},
		{
			Name:        "resource",
			Usage:       "Stack resource queries.",
			Subcommands: stackresourcecommands.Get(),
		},
		{
			Name:        "template",
			Usage:       "Stack template queries.",
			Subcommands: stacktemplatecommands.Get(),
		},
	}
}
