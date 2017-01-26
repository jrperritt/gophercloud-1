package volumeattachment

import (
	"github.com/gophercloud/gophercloud/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "compute volume-attachment"

type VolumeAttachmentV2Command struct {
	traits.Commandable
	traits.ComputeV2able
}

// Get returns all the commands allowed for a `compute volume-attachment` v2 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		create,
		get,
		remove,
	}
}
