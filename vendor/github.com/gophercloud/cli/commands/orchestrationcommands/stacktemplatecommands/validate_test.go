package stacktemplatecommands

import (
	"flag"
	"testing"

	"github.com/gophercloud/cli/handler"
	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	th "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/testhelper"
)

func TestValidateContext(t *testing.T) {
	app := cli.NewApp()
	flagset := flag.NewFlagSet("flags", 1)
	c := cli.NewContext(app, flagset, nil)
	cmd := &commandValidate{
		Ctx: &handler.Context{
			CLIContext: c,
		},
	}
	expected := cmd.Ctx
	actual := cmd.Context()
	th.AssertDeepEquals(t, expected, actual)
}

func TestValidateKeys(t *testing.T) {
	cmd := &commandValidate{}
	expected := keysValidate
	actual := cmd.Keys()
	th.AssertDeepEquals(t, expected, actual)
}

func TestValidateServiceClientType(t *testing.T) {
	cmd := &commandValidate{}
	expected := serviceClientType
	actual := cmd.ServiceClientType()
	th.AssertEquals(t, expected, actual)
}
