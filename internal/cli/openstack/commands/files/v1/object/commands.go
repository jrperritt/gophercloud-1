package object

import (
	"github.com/gophercloud/gophercloud/internal/cli/lib/traits"
	"gopkg.in/urfave/cli.v1"
)

var commandPrefix = "files object"

// ObjectV1Command should be embedded by operations that use the
// Object Storage v1 client
type ObjectV1Command struct {
	traits.Commandable
	traits.FilesV1able
	container string
	name      string
}

// Get returns all the commands allowed for a `files object` v1 request.
func Get() []cli.Command {
	return []cli.Command{
		list,
		get,
		upload,
		download,
		remove,
	}
}
