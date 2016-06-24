package openstack

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

// CommandFlags returns the flags for a given command. It takes as a parameter
// a function for returning flags specific to that command, and then appends those
// flags with flags that are valid for all commands.
func CommandFlags(flags []cli.Flag, keys []string) []cli.Flag {
	if len(keys) > 0 {
		fields := make([]string, len(keys))
		for i, key := range keys {
			fields[i] = strings.Join(strings.Split(strings.ToLower(key), " "), "-")
		}
		flagFields := cli.StringFlag{
			Name:  "fields",
			Usage: fmt.Sprintf("[optional] Only return these comma-separated case-insensitive fields.\n\tChoices: %s", strings.Join(fields, ", ")),
		}
		flags = append(flags, flagFields)
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
