package traits

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib"
	"gopkg.in/urfave/cli.v1"
)

type Commandable struct {
	ctx    *cli.Context
	sc     *gophercloud.ServiceClient
	execch chan interface{}
	outch  chan interface{}
}

func (c *Commandable) SetContext(ctx *cli.Context) error {
	c.ctx = ctx
	return nil
}

func (c *Commandable) Context() *cli.Context {
	return c.ctx
}

func (c *Commandable) SetServiceClient(sc *gophercloud.ServiceClient) {
	c.sc = sc
}

func (c *Commandable) ServiceClient() *gophercloud.ServiceClient {
	return c.sc
}

func (c *Commandable) HandleFlags() error {
	return nil
}

// IDOrName is a function for retrieving a resource's unique identifier based on
// whether an `id` or a `name` flag was provided
func (c *Commandable) IDOrName(idFromNameFunc func(*gophercloud.ServiceClient, string) (string, error)) (string, error) {
	switch c.ctx.IsSet("id") {
	case true:
		switch c.ctx.IsSet("name") {
		case true:
			return "", fmt.Errorf("Only one of either --id or --name may be provided.")
		case false:
			return c.ctx.String("id"), nil
		}
	case false:
		switch c.ctx.IsSet("name") {
		case true:
			name := c.ctx.String("name")
			id, err := idFromNameFunc(c.sc, name)
			return id, err
		}
	}
	return "", lib.ErrMissingFlag{Msg: "One of either --id or --name must be provided."}
}

// CheckFlagsSet checks that the given flag names are set for the command.
func (c *Commandable) CheckFlagsSet(flagNames []string) error {
	for _, flagName := range flagNames {
		if !c.ctx.IsSet(flagName) {
			return lib.ErrMissingFlag{Msg: fmt.Sprintf("--%s is required.", flagName)}
		}
	}
	return nil
}

// ValidateKVFlag is a function used for verifying the format of a key-value flag.
func (c *Commandable) ValidateKVFlag(flagName string) (map[string]string, error) {
	kv := make(map[string]string)
	kvStrings := strings.Split(c.ctx.String(flagName), ",")
	for _, kvString := range kvStrings {
		temp := strings.Split(kvString, "=")
		if len(temp) != 2 {
			return nil, lib.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s for --%s.\n", kvString, flagName)}
		}
		kv[temp[0]] = temp[1]
	}
	return kv, nil
}

// ValidateStructFlag is a function used for verifying the format of a struct flag.
func (c *Commandable) ValidateStructFlag(flagValues []string) ([]map[string]interface{}, error) {
	valSliceMap := make([]map[string]interface{}, len(flagValues))
	for i, flagValue := range flagValues {
		kvStrings := strings.Split(flagValue, ",")
		m := make(map[string]interface{})
		for _, kvString := range kvStrings {
			temp := strings.Split(kvString, "=")
			if len(temp) != 2 {
				return nil, lib.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s.\n", kvString)}
			}
			m[temp[0]] = temp[1]
		}
		valSliceMap[i] = m
	}
	return valSliceMap, nil
}
