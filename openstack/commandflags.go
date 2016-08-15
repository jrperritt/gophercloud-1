package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/cli/lib"
	"gopkg.in/urfave/cli.v1"
)

// CommandFlags returns the flags for a given command. It takes as a parameter
// a function for returning flags specific to that command, and then appends those
// flags with flags that are valid for all commands.
func CommandFlags(c lib.Commander) []cli.Flag {
	flags := c.Flags()

	if fieldser, ok := c.(lib.Fieldser); ok {
		keys := fieldser.Fields()
		if len(keys) > 0 {
			usage := "[optional] Only return these comma-separated case-insensitive fields."
			if keys[0] != "" {
				usage = fmt.Sprintf(usage+"\n\tChoices: %s", strings.Join(keys, ", "))
			}

			flagFields := cli.StringFlag{
				Name:  "fields",
				Usage: usage,
			}

			flags = append(flags, flagFields)
		}
	}

	if waiter, ok := c.(lib.Waiter); ok {
		flags = append(flags, waiter.WaitFlags()...)
	}

	flags = append(flags, GlobalFlags()...)

	return flags
}

// CompleteFlags returns the possible flags for bash completion.
func CompleteFlags(flags []cli.Flag) {
	for _, flag := range flags {
		flagName := ""
		switch f := flag.(type) {
		case cli.StringFlag:
			flagName = f.Name
		case cli.IntFlag:
			flagName = f.Name
		case cli.BoolFlag:
			flagName = f.Name
		case cli.StringSliceFlag:
			flagName = f.Name
		default:
			continue
		}
		fmt.Println("--" + flagName)
	}
}
