package openstack

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
)

// Command is the type that commands have.
type Command struct {
	// cli.Context is the context that the `cli` library uses. Used to
	// access flags.
	*cli.Context
	// ServiceClient is the service client used to authenticate the user
	// and carry out the requests while processing the command.
	serviceClient *gophercloud.ServiceClient
	// ServiceClientType is the type of service client used (e.g. compute).
	serviceClientType string
	// results is a channel into which commands send results. It allows for streaming
	// output.
	//results chan *gophercloudCLILib.Resulter
	// flags are the command-specific flags
	flags []cli.Flag
	// fields are the fields available to output. These may be limited by the `fields`
	// flag.
	fields []string
	// logger is used to log information acquired while processing the command.
	logger *logrus.Logger
}

func NewCommand(c lib.Commander, flags []cli.Flag, serviceClientType string) cli.Command {
	c.SetFlags(flags)
	c.SetServiceClientType(serviceClientType)
	c.SetFields(util.BuildFields(c.ReturnType()))
	return cli.Command{
		Name:         c.Name(),
		Usage:        c.Usage(),
		Description:  c.Description(),
		Action:       c.Action,
		Flags:        c.Flags(),
		BashComplete: c.BashComplete,
	}
}

func (c *Command) Name() string {
	return ""
}

func (c *Command) Usage() string {
	return ""
}

func (c *Command) Description() string {
	return ""
}

func (c *Command) ReturnType() reflect.Type {
	return reflect.TypeOf("")
}

func (c *Command) HandleFlags() error {
	return nil
}

func (c *Command) Execute(_ lib.Resourcer) lib.Resulter {
	return nil
}

func (c *Command) Action(cliContext *cli.Context) {
	c.Context = cliContext
	lib.Run(cliContext, c)
}

func (c Command) Fields() []string {
	return c.fields
}

func (c *Command) SetFields(fields []string) {
	c.fields = fields
}

func (c Command) Flags() []cli.Flag {
	return CommandFlags(c.flags, c.Fields())
}

func (c *Command) SetFlags(flags []cli.Flag) {
	c.flags = flags
}

func (c Command) BashComplete(cliContext *cli.Context) {
	CompleteFlags(c.Flags())
}

func (c *Command) SetServiceClient(serviceClient *gophercloud.ServiceClient) {
	c.serviceClient = serviceClient
}

func (c Command) ServiceClient() *gophercloud.ServiceClient {
	return c.serviceClient
}

func (c *Command) SetServiceClientType(serviceClientType string) {
	c.serviceClientType = serviceClientType
}

func (c Command) ServiceClientType() string {
	return c.serviceClientType
}

func (c Command) RunCommand(resultsChannel chan lib.Resulter) error {
	result := c.Execute(new(Resource))
	resultsChannel <- result
	close(resultsChannel)
	return nil
}

// IDOrName is a function for retrieving a resource's unique identifier based on
// whether an `id` or a `name` flag was provided
func (c Command) IDOrName(idFromNameFunc func(*gophercloud.ServiceClient, string) (string, error)) (string, error) {
	if c.IsSet("id") {
		if c.IsSet("name") {
			return "", fmt.Errorf("Only one of either --id or --name may be provided.")
		}
		return c.String("id"), nil
	} else if c.IsSet("name") {
		name := c.String("name")
		id, err := idFromNameFunc(c.serviceClient, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	} else {
		return "", lib.ErrMissingFlag{Msg: "One of either --id or --name must be provided."}
	}
}

// CheckFlagsSet checks that the given flag names are set for the command.
func (c Command) CheckFlagsSet(flagNames []string) error {
	for _, flagName := range flagNames {
		if !c.IsSet(flagName) {
			return lib.ErrMissingFlag{Msg: fmt.Sprintf("--%s is required.", flagName)}
		}
	}
	return nil
}

// CheckKVFlag is a function used for verifying the format of a key-value flag.
func (c Command) ValidateKVFlag(flagName string) (map[string]string, error) {
	kv := make(map[string]string)
	kvStrings := strings.Split(c.String(flagName), ",")
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
func (c Command) ValidateStructFlag(flagValues []string) ([]map[string]interface{}, error) {
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
