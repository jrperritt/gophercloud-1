package stacktemplatecommands

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/gophercloud/cli/commandoptions"
	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/github.com/fatih/structs"
	osStackTemplates "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacktemplates"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacktemplates"
	"github.com/gophercloud/cli/util"
)

var validate = cli.Command{
	Name:        "validate",
	Usage:       util.Usage(commandPrefix, "validate", "[--template-file <templateFile> | --template-url <templateURL>]"),
	Description: "Validate a specified template",
	Action:      actionValidate,
	Flags:       commandoptions.CommandFlags(flagsValidate, keysValidate),
	BashComplete: func(c *cli.Context) {
		commandoptions.CompleteFlags(commandoptions.CommandFlags(flagsValidate, keysValidate))
	},
}

func flagsValidate() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "template-file",
			Usage: "[optional; required if `template-url` isn't provided] The path to template file.",
		},
		cli.StringFlag{
			Name:  "template-url",
			Usage: "[optional; required if `template-file` isn't provided] The url to template.",
		},
	}
}

type paramsValidate struct {
	opts *osStackTemplates.ValidateOpts
}

var keysValidate = []string{"Description", "Parameters", "ParameterGroups"}

type commandValidate handler.Command

func actionValidate(c *cli.Context) {
	command := &commandValidate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	handler.Handle(command)
}

func (command *commandValidate) Context() *handler.Context {
	return command.Ctx
}

func (command *commandValidate) Keys() []string {
	return keysValidate
}

func (command *commandValidate) ServiceClientType() string {
	return serviceClientType
}

func (command *commandValidate) HandleFlags(resource *handler.Resource) error {
	c := command.Ctx.CLIContext
	opts := osStackTemplates.ValidateOpts{}

	// check if either template url or template file is set
	if c.IsSet("template-file") {
		abs, err := filepath.Abs(c.String("template-file"))
		if err != nil {
			return err
		}
		template, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}
		opts.Template = string(template)
	} else if c.IsSet("template-url") {
		opts.TemplateURL = c.String("template-url")
	} else {
		return errors.New("Neither template-file nor template-url specified")
	}

	resource.Params = &paramsValidate{
		opts: &opts,
	}
	return nil
}

func (command *commandValidate) Execute(resource *handler.Resource) {
	params := resource.Params.(*paramsValidate).opts

	result, err := stacktemplates.Validate(command.Ctx.ServiceClient, params).Extract()
	if err != nil {
		resource.Err = err
		return
	}
	resource.Result = structs.Map(result)
}

func (command *commandValidate) PreCSV(resource *handler.Resource) error {
	resource.FlattenMap("Parameters")
	resource.FlattenMap("ParameterGroups")
	return nil
}

func (command *commandValidate) PreTable(resource *handler.Resource) error {
	return command.PreCSV(resource)
}
