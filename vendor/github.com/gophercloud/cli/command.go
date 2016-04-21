package main

import (
	"fmt"
	"strings"

	openstackCLILib "github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/vendor/github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
)

// Command is the type that commands have.
type Command struct {
	// CLIContext is the context that the `cli` library uses. Used to
	// access flags.
	CLIContext *cli.Context
	// ServiceClient is the service client used to authenticate the user
	// and carry out the requests while processing the command.
	ServiceClient *gophercloud.ServiceClient
	// ServiceClientType is the type of service client used (e.g. compute).
	ServiceClientType string
	// Results is a channel into which commands send results. It allows for streaming
	// output.
	Results chan *Resource
	// Keys are the fields available to output. These may be limited by the `fields`
	// flag.
	Keys []string
	// logger is used to log information acquired while processing the command.
	logger *logrus.Logger
}

// IDOrName is a function for retrieving a resource's unique identifier based on
// whether an `id` or a `name` flag was provided
func (ctx Command) IDOrName(idFromNameFunc func(*gophercloud.ServiceClient, string) (string, error)) (string, error) {
	if ctx.CLIContext.IsSet("id") {
		if ctx.CLIContext.IsSet("name") {
			return "", fmt.Errorf("Only one of either --id or --name may be provided.")
		}
		return ctx.CLIContext.String("id"), nil
	} else if ctx.CLIContext.IsSet("name") {
		name := ctx.CLIContext.String("name")
		id, err := idFromNameFunc(ctx.ServiceClient, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	} else {
		return "", openstackCLILib.ErrMissingFlag{Msg: "One of either --id or --name must be provided."}
	}
}

// CheckFlagsSet checks that the given flag names are set for the command.
func (ctx Command) CheckFlagsSet(flagNames []string) error {
	for _, flagName := range flagNames {
		if !ctx.CLIContext.IsSet(flagName) {
			return openstackCLILib.ErrMissingFlag{Msg: fmt.Sprintf("--%s is required.", flagName)}
		}
	}
	return nil
}

// CheckKVFlag is a function used for verifying the format of a key-value flag.
func (ctx Command) ValidateKVFlag(flagName string) (map[string]string, error) {
	kv := make(map[string]string)
	kvStrings := strings.Split(ctx.CLIContext.String(flagName), ",")
	for _, kvString := range kvStrings {
		temp := strings.Split(kvString, "=")
		if len(temp) != 2 {
			return nil, openstackCLILib.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s for --%s.\n", kvString, flagName)}
		}
		kv[temp[0]] = temp[1]
	}
	return kv, nil
}

// CheckStructFlag is a function used for verifying the format of a struct flag.
func (ctx Command) ValidateStructFlag(flagValues []string) ([]map[string]interface{}, error) {
	valSliceMap := make([]map[string]interface{}, len(flagValues))
	for i, flagValue := range flagValues {
		kvStrings := strings.Split(flagValue, ",")
		m := make(map[string]interface{})
		for _, kvString := range kvStrings {
			temp := strings.Split(kvString, "=")
			if len(temp) != 2 {
				return nil, openstackCLILib.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s.\n", kvString)}
			}
			m[temp[0]] = temp[1]
		}
		valSliceMap[i] = m
	}
	return valSliceMap, nil
}
