package keypair

import (
	"io/ioutil"

	"github.com/gophercloud/cli/lib/traits"
	"github.com/gophercloud/cli/openstack"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"gopkg.in/urfave/cli.v1"
)

type CommandUpload struct {
	KeypairV2Command
	traits.DataResp
	opts keypairs.CreateOptsBuilder
}

var (
	cUpload = new(CommandUpload)

	flagsUpload = openstack.CommandFlags(cUpload)
)

var upload = cli.Command{
	Name:         "upload",
	Usage:        util.Usage(commandPrefix, "upload", "--name <NAME> --file <FILEPATH>"),
	Description:  "Uploads a keypair",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cUpload) },
	Flags:        flagsUpload,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsUpload) },
}

func (c *CommandUpload) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[required] The name of the keypair",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "[required] The name of the file containing the public key.",
		},
	}
}

func (c *CommandUpload) HandleFlags() error {
	err := c.CheckFlagsSet([]string{"name", "file"})
	if err != nil {
		return err
	}

	s := c.Context.String("file")
	pk, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}

	c.opts = &keypairs.CreateOpts{
		Name:      c.Context.String("name"),
		PublicKey: string(pk),
	}

	return nil
}

func (c *CommandUpload) Execute(item interface{}, out chan interface{}) {
	var m map[string]map[string]interface{}
	err := keypairs.Create(c.ServiceClient, c.opts).ExtractInto(&m)
	switch err {
	case nil:
		out <- m["keypair"]
	default:
		out <- err
	}
}
