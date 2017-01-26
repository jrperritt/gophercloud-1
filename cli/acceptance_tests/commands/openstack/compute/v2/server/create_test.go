package server

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/cli/openstack/commands/compute/v2/server"
	"github.com/gophercloud/gophercloud/cli/util"
)

var base = fmt.Sprintf("%s create", server.CommandPrefix)

func TestServer(t *testing.T) {
	testsCreateSingle := []struct {
		name  string
		flags string
	}{
		{
			"Create-Single-NoWait-NoProgress",
			"-name=s1 -image-name=cirros -flavor-id=1 -quiet",
		},
		{
			"Create-Single-Wait-NoProgress",
			"-name=s2 -image-name=cirros -flavor-id=1 -quiet -wait",
		},
		{
			"Create-Single-Progress",
			"-name=s3 -image-name=cirros -flavor-id=1",
		},
	}
	for _, testCreateSingle := range testsCreateSingle {
		cmd := exec.Command(util.Name, strings.Split(base+testCreateSingle.flags, " ")...)
		err := cmd.Run()
		if err != nil {
			t.Fatal(fmt.Sprintf("%s: %s", testCreateSingle.name, err))
		}
	}
}
