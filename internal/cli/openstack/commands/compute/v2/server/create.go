package server

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/internal/cli/lib"
	"github.com/gophercloud/gophercloud/internal/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"github.com/gophercloud/gophercloud/internal/cli/openstack"
	"github.com/gophercloud/gophercloud/internal/cli/util"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

type CommandCreate struct {
	ServerV2Command
	traits.Waitable
	traits.Pipeable
	traits.PercentageProgressable
	traits.Fieldsable
	opts servers.CreateOptsBuilder
}

type createdata struct {
	traits.ProgressItemPct
	res map[string]interface{}
}

func newcreatedata() *createdata {
	d := new(createdata)
	d.ProgressItem.Init()
	return d
}

var (
	cCreate                              = new(CommandCreate)
	_       interfaces.PipeCommander     = cCreate
	_       interfaces.Waiter            = cCreate
	_       interfaces.PercentProgresser = cCreate

	flagsCreate = openstack.CommandFlags(cCreate)
)

var create = cli.Command{
	Name:         "create",
	Usage:        util.Usage(CommandPrefix, "create", "[--name <name> | --stdin name] \n\t [--image-id <imageID> | --image-name <imageName>] [--flavor-id <flavorID> | --flavor-name <flavorName]"),
	Description:  "Creates a server",
	Action:       func(ctx *cli.Context) error { return openstack.Action(ctx, cCreate) },
	Flags:        flagsCreate,
	BashComplete: func(_ *cli.Context) { util.CompleteFlags(flagsCreate) },
}

func (c *CommandCreate) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "[optional; required if `stdin` isn't provided] The name that the instance should have.",
		},
		cli.StringFlag{
			Name:  "stdin",
			Usage: "[optional; required if `name` isn't provided] The field being piped into STDIN. Valid values are: name",
		},
		cli.StringFlag{
			Name:  "image-id",
			Usage: "[optional; required if `image-name` or `block-device` is not provided] The image ID from which to create the server.",
		},
		cli.StringFlag{
			Name:  "image-name",
			Usage: "[optional; required if `image-id` or `block-device` is not provided] The name of the image from which to create the server.",
		},
		cli.StringFlag{
			Name:  "flavor-id",
			Usage: "[optional; required if `flavor-name` is not provided] The flavor ID that the server should have.",
		},
		cli.StringFlag{
			Name:  "flavor-name",
			Usage: "[optional; required if `flavor-id` is not provided] The name of the flavor that the server should have.",
		},
		cli.StringFlag{
			Name:  "security-groups",
			Usage: "[optional] A comma-separated string of names of the security groups to which this server should belong.",
		},
		cli.StringFlag{
			Name: "personality",
			Usage: "[optional] A comma-separated list of key=value pairs. The key is the\n" +
				"\tdestination to inject the file on the created server; the value is the its local location.\n" +
				"\tExample: --personality \"C:\\cloud-automation\\bootstrap.cmd=open_hatch.cmd\"",
		},
		cli.StringFlag{
			Name:  "user-data",
			Usage: "[optional] Configuration information or scripts to use after the server boots.",
		},
		cli.StringFlag{
			Name:  "networks",
			Usage: "[optional] A comma-separated string of IDs of the networks to attach to this server. If not provided, a public and private network will be attached.",
		},
		cli.StringFlag{
			Name:  "metadata",
			Usage: "[optional] A comma-separated string of key=value pairs.",
		},
		cli.StringFlag{
			Name:  "admin-pass",
			Usage: "[optional] The root password for the server. If not provided, one will be randomly generated and returned in the output.",
		},
		cli.StringFlag{
			Name:  "keypair",
			Usage: "[optional] The name of the already-existing SSH KeyPair to be injected into this server.",
		},
		cli.StringFlag{
			Name: "block-device",
			Usage: strings.Join([]string{"[optional] Used to boot from volume.",
				"\tIf provided, the instance will be created based upon the comma-separated key=value pairs provided to this flag.",
				"\tOptions:",
				"\t\tsource-type\t[required] The source type of the device. Options: volume, snapshot, image.",
				"\t\tsource-id\t[required] The ID of the source resource (volume, snapshot, or image) from which to create the instance.",
				"\t\tboot-index\t[optional] The boot index of the device. Default is 0.",
				"\t\tdelete-on-termination\t[optional] Whether or not to delete the attached volume when the server is delete. Default is false. Options: true, false.",
				"\t\tdestination-type\t[optional] The type that gets created. Options: volume, local.",
				"\t\tvolume-size\t[optional] The size of the volume to create (in gigabytes).",
				"\tExamle: --block-device source-type=image,source-id=bb02b1a3-bc77-4d17-ab5b-421d89850fca,volume-size=100,destination-type=volume,delete-on-termination=false",
			}, "\n"),
		},
	}
}

var flagsCreateExt = []cli.Flag{
	cli.StringFlag{
		Name:  "keypair",
		Usage: "[optional] The name of the already-existing SSH KeyPair to be injected into this server.",
	},
	cli.StringFlag{
		Name: "block-device",
		Usage: strings.Join([]string{"[optional] Used to boot from volume.",
			"\tIf provided, the instance will be created based upon the comma-separated key=value pairs provided to this flag.",
			"\tOptions:",
			"\t\tsource-type\t[required] The source type of the device. Options: volume, snapshot, image.",
			"\t\tsource-id\t[required] The ID of the source resource (volume, snapshot, or image) from which to create the instance.",
			"\t\tboot-index\t[optional] The boot index of the device. Default is 0.",
			"\t\tdelete-on-termination\t[optional] Whether or not to delete the attached volume when the server is delete. Default is false. Options: true, false.",
			"\t\tdestination-type\t[optional] The type that gets created. Options: volume, local.",
			"\t\tvolume-size\t[optional] The size of the volume to create (in gigabytes).",
			"\tExamle: --block-device source-type=image,source-id=bb02b1a3-bc77-4d17-ab5b-421d89850fca,volume-size=100,destination-type=volume,delete-on-termination=false",
		}, "\n"),
	},
}

func (c *CommandCreate) HandleFlags() error {
	opts := &servers.CreateOpts{
		ImageRef:      c.Context().String("image-id"),
		ImageName:     c.Context().String("image-name"),
		FlavorRef:     c.Context().String("flavor-id"),
		FlavorName:    c.Context().String("flavor-name"),
		AdminPass:     c.Context().String("admin-pass"),
		ServiceClient: c.ServiceClient(),
	}

	if c.Context().IsSet("security-groups") {
		opts.SecurityGroups = strings.Split(c.Context().String("security-groups"), ",")
	}

	if c.Context().IsSet("user-data") {
		abs, err := filepath.Abs(c.Context().String("user-data"))
		if err != nil {
			return err
		}
		userData, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}
		opts.UserData = userData
		opts.ConfigDrive = gophercloud.Enabled
	}

	if c.Context().IsSet("personality") {

		filesToInjectMap, err := c.ValidateKVFlag("personality")
		if err != nil {
			return err
		}

		if len(filesToInjectMap) > 5 {
			return fmt.Errorf("A maximum of 5 files may be provided for the `personality` flag")
		}

		filesToInject := make(servers.Personality, 0)
		for destinationPath, localPath := range filesToInjectMap {
			localAbsFilePath, err := filepath.Abs(localPath)
			if err != nil {
				return err
			}

			fileData, err := ioutil.ReadFile(localAbsFilePath)
			if err != nil {
				return err
			}

			if len(fileData)+len(destinationPath) > 1000 {
				return fmt.Errorf("The maximum length of a file-path-and-content pair for `personality` is 1000 bytes."+
					" Current pair size: path (%s): %d, content: %d", destinationPath, len(destinationPath), len(fileData))
			}

			filesToInject = append(filesToInject, &servers.File{
				Path:     destinationPath,
				Contents: fileData,
			})
		}
		opts.Personality = filesToInject
	}

	if c.Context().IsSet("networks") {
		netIDs := strings.Split(c.Context().String("networks"), ",")
		networks := make([]servers.Network, len(netIDs))
		for i, netID := range netIDs {
			networks[i] = servers.Network{
				UUID: netID,
			}
		}
		opts.Networks = networks
	}

	if c.Context().IsSet("metadata") {
		metadata, err := c.ValidateKVFlag("metadata")
		if err != nil {
			return err
		}
		opts.Metadata = metadata
	}

	// -------------- Extensions logic starts here -------------------------
	var optsExt servers.CreateOptsBuilder = opts

	if c.Context().IsSet("keypair") {
		optsExt = keypairs.CreateOptsExt{
			CreateOptsBuilder: opts,
			KeyName:           c.Context().String("keypair"),
		}
	}

	if c.Context().IsSet("block-device") {
		bfvMap, err := c.ValidateKVFlag("block-device")
		if err != nil {
			return err
		}

		sourceID, ok := bfvMap["source-id"]
		if !ok {
			return fmt.Errorf("The source-id key is required when using the --block-device flag.\n")
		}

		sourceTypeRaw, ok := bfvMap["source-type"]
		if !ok {
			return fmt.Errorf("The source-type key is required when using the --block-device flag.\n")
		}
		var sourceType bootfromvolume.SourceType
		switch sourceTypeRaw {
		case "volume", "image", "snapshot":
			sourceType = bootfromvolume.SourceType(sourceTypeRaw)
		default:
			return fmt.Errorf("Invalid value for source-type: %s. Options are: volume, image, snapshot.\n", sourceType)
		}

		bd := bootfromvolume.BlockDevice{
			SourceType: sourceType,
			UUID:       sourceID,
		}

		if volumeSizeRaw, ok := bfvMap["volume-size"]; ok {
			volumeSize, err := strconv.ParseInt(volumeSizeRaw, 10, 16)
			if err != nil {
				return fmt.Errorf("Invalid value for volume-size: %d. Value must be an integer.\n", volumeSize)
			}
			bd.VolumeSize = int(volumeSize)
		}

		if deleteOnTerminationRaw, ok := bfvMap["delete-on-termination"]; ok {
			deleteOnTermination, err := strconv.ParseBool(deleteOnTerminationRaw)
			if err != nil {
				return fmt.Errorf("Invalid value for delete-on-termination: %v. Options are: true, false.\n", deleteOnTermination)
			}
			bd.DeleteOnTermination = deleteOnTermination
		}

		if bootIndexRaw, ok := bfvMap["boot-index"]; ok {
			bootIndex, err := strconv.ParseInt(bootIndexRaw, 10, 8)
			if err != nil {
				return fmt.Errorf("Invalid value for boot-index: %d. Value must be an integer.\n", bootIndex)
			}
			bd.BootIndex = int(bootIndex)
		}

		if destinationType, ok := bfvMap["destination-type"]; ok {
			if destinationType != "volume" && destinationType != "local" {
				return fmt.Errorf("Invalid value for destination-type: %s. Options are: volume, local.\n", destinationType)
			}
			bd.DestinationType = bootfromvolume.DestinationType(destinationType)
		}

		optsExt = bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: optsExt,
			BlockDevice:       []bootfromvolume.BlockDevice{bd},
		}
	}

	c.opts = optsExt

	return nil
}

func (c *CommandCreate) HandleSingle() (interface{}, error) {
	return c.Context().String("name"), c.CheckFlagsSet([]string{"name"})
}

func (c *CommandCreate) Execute(item interface{}, out chan (interface{})) {
	var m map[string]map[string]interface{}
	opts := *c.opts.(*servers.CreateOpts)
	opts.Name = item.(string)
	err := servers.Create(c.ServiceClient(), opts).ExtractInto(&m)
	if err != nil {
		out <- err
		return
	}
	d := newcreatedata()
	d.res = m["server"]
	d.SetID(d.res["id"].(string))
	out <- d
}

func (c *CommandCreate) PipeFieldOptions() []string {
	return []string{"name"}
}

func (c *CommandCreate) WaitFor(raw interface{}, donech chan<- interface{}) {
	d := raw.(*createdata)

	err := gophercloud.WaitFor(900, func() (bool, error) {
		var m map[string]map[string]interface{}
		lib.Log.Debugf("running servers.Get for item: %s", d.ID())
		err := servers.Get(c.ServiceClient(), d.ID()).ExtractInto(&m)
		if err != nil {
			return false, err
		}

		switch m["server"]["status"].(string) {
		case "ACTIVE":
			m["server"]["adminPass"] = d.res["adminPass"].(string)
			d.EndCh() <- m["server"]
			return true, nil
		default:
			d.UpCh() <- m["server"]["progress"].(float64)
			return false, nil
		}
	})

	if err != nil {
		d.EndCh() <- err
	}
}
