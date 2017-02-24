package main

import (
	"fmt"

	"github.com/gophercloud/gophercloud/internal/cli/util"
	"gopkg.in/urfave/cli.v1"
)

// GlobalFlags returns the flags that can be in any command, such as
// output flags and authentication flags.
var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "username",
		Usage: "The username with which to authenticate.",
	},
	cli.StringFlag{
		Name:  "user-id",
		Usage: "The user ID with which to authenticate.",
	},
	cli.StringFlag{
		Name:  "password",
		Usage: "The password with which to authenticate.",
	},
	cli.StringFlag{
		Name:  "tenant-id",
		Usage: "The tenant ID of the user to authenticate as. May only be provided as a command-line flag.",
	},
	cli.StringFlag{
		Name:  "auth-token",
		Usage: "The authentication token of the user to authenticate as. This must be used with the `tenant-id` flag.",
	},
	cli.StringFlag{
		Name:  "auth-url",
		Usage: "The endpoint to which authenticate.",
	},
	cli.StringFlag{
		Name:  "region",
		Usage: "The region to which authenticate.",
	},
	cli.StringFlag{
		Name:  "profile",
		Usage: "The config file profile from which to load flags.",
	},
	cli.BoolFlag{
		Name:  "no-cache",
		Usage: "Do not get or set authentication credentials in the cache.",
	},
	cli.StringFlag{
		Name: "log",
		Usage: "Print variout levels of information from the command." +
			"\nOptions are debug, info, warn",
	},
}

func globalComplete(ctx *cli.Context) {
	for _, cmd := range ctx.App.Commands {
		fmt.Println(cmd.Name)
	}
	util.CompleteFlags(globalFlags)
}
