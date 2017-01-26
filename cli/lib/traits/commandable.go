package traits

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib"
	"gopkg.in/urfave/cli.v1"
)

type Commandable struct {
	Context       *cli.Context
	ServiceClient *gophercloud.ServiceClient
}

func (c *Commandable) SetContext(ctx *cli.Context) error {
	c.Context = ctx
	return nil
}

func (c *Commandable) SetServiceClient(sc *gophercloud.ServiceClient) error {
	c.ServiceClient = sc
	return nil
}

func (c *Commandable) HandleFlags() error {
	return nil
}

// IDOrName is a function for retrieving a resource's unique identifier based on
// whether an `id` or a `name` flag was provided
func (c *Commandable) IDOrName(idFromNameFunc func(*gophercloud.ServiceClient, string) (string, error)) (string, error) {
	switch c.Context.IsSet("id") {
	case true:
		switch c.Context.IsSet("name") {
		case true:
			return "", fmt.Errorf("Only one of either --id or --name may be provided.")
		case false:
			return c.Context.String("id"), nil
		}
	case false:
		switch c.Context.IsSet("name") {
		case true:
			name := c.Context.String("name")
			id, err := idFromNameFunc(c.ServiceClient, name)
			return id, err
		}
	}
	return "", lib.ErrMissingFlag{Msg: "One of either --id or --name must be provided."}
}

// CheckFlagsSet checks that the given flag names are set for the command.
func (c *Commandable) CheckFlagsSet(flagNames []string) error {
	for _, flagName := range flagNames {
		if !c.Context.IsSet(flagName) {
			return lib.ErrMissingFlag{Msg: fmt.Sprintf("--%s is required.", flagName)}
		}
	}
	return nil
}

// CheckKVFlag is a function used for verifying the format of a key-value flag.
func (c *Commandable) ValidateKVFlag(flagName string) (map[string]string, error) {
	kv := make(map[string]string)
	kvStrings := strings.Split(c.Context.String(flagName), ",")
	for _, kvString := range kvStrings {
		temp := strings.Split(kvString, "=")
		if len(temp) != 2 {
			return nil, lib.ErrFlagFormatting{Msg: fmt.Sprintf("Expected key1=value1,key2=value2 format but got %s for --%s.\n", kvString, flagName)}
		}
		kv[temp[0]] = temp[1]
	}
	return kv, nil
}

// CheckStructFlag is a function used for verifying the format of a struct flag.
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

type DataResp struct{}

func (c *DataResp) Fields() []string {
	return []string{""}
}

type MsgResp struct{}
