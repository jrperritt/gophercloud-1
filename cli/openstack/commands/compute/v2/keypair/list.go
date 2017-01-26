package keypair

import (
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"github.com/gophercloud/gophercloud/cli/openstack"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"gopkg.in/urfave/cli.v1"
)

type CommandList struct {
	KeypairV2Command
	traits.DataResp
	opts keypairs.CreateOptsBuilder
}

var (
	cList = new(CommandList)

	flagsList = openstack.CommandFlags(cList)
)

var list = cli.Command{
	Name:         "list",
	Usage:        util.Usage(commandPrefix, "list", ""),
	Description:  "Lists all keypairs",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cList) },
	Flags:        flagsList,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsList) },
}

func (c *CommandList) Flags() []cli.Flag {
	return nil
}

func (c *CommandList) DefaultTableFields() []string {
	return []string{"name", "fingerprint"}
}

func (c *CommandList) Execute(item interface{}, out chan interface{}) {
	p, err := keypairs.List(c.ServiceClient).AllPages()
	if err != nil {
		out <- err
		return
	}

	var m map[string][]map[string]interface{}
	err = (p.(keypairs.KeyPairPage)).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}
	kps := m["keypairs"]

	r := make([]map[string]interface{}, len(kps))
	for i, kp := range kps {
		r[i] = kp["keypair"].(map[string]interface{})
	}

	out <- r
}
