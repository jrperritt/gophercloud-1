package openstack

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"

	"gopkg.in/urfave/cli.v1"
)

// CommandFlags returns the flags for a given command. It takes as a parameter
// a function for returning flags specific to that command, and then appends those
// flags with flags that are valid for all commands.
func CommandFlags(cmd interfaces.Commander) (flags []cli.Flag) {
	if flagser, ok := cmd.(interfaces.Flagser); ok {
		flags = flagser.Flags()
	}

	if fieldser, ok := cmd.(interfaces.Fieldser); ok {
		flags = append(flags, fieldser.FieldsFlags()...)
	}

	if waiter, ok := cmd.(interfaces.Waiter); ok {
		flags = append(flags, waiter.WaitFlags()...)
	}

	if progresser, ok := cmd.(interfaces.Progresser); ok {
		flags = append(flags, progresser.ProgressFlags()...)
	}

	if tabler, ok := cmd.(interfaces.Tabler); ok {
		flags = append(flags, tabler.TableFlags()...)
	}

	if piper, ok := cmd.(interfaces.PipeCommander); ok {
		flags = append(flags, piper.PipeFlags()...)
	}

	return flags
}
