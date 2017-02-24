package util

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

// Complete returns the possible flags for bash completion.
func CompleteFlags(flags []cli.Flag) {
	for _, flag := range flags {
		fmt.Println("--" + flag.GetName())
	}
}
