package volumeattachment

import (
	"fmt"

	"github.com/gophercloud/gophercloud/internal/cli/lib"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"gopkg.in/urfave/cli.v1"
)

func serverIDorName(ctx *cli.Context, sc *gophercloud.ServiceClient) (string, error) {
	switch {
	case ctx.IsSet("server-id"):
		if ctx.IsSet("server-name") {
			return "", fmt.Errorf("Only one of either --server-id or --server-name may be provided.")
		}
		return ctx.String("server-id"), nil
	case ctx.IsSet("server-name"):
		name := ctx.String("server-name")
		id, err := servers.IDFromName(sc, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	default:
		return "", lib.ErrMissingFlag{Msg: "One of either --server-id or --server-name must be provided."}
	}
}

func volumeIDorName(ctx *cli.Context, sc *gophercloud.ServiceClient) (string, error) {
	switch {
	case ctx.IsSet("volume-id"):
		if ctx.IsSet("volume-name") {
			return "", fmt.Errorf("Only one of either --volume-id or --volume-name may be provided.")
		}
		return ctx.String("volume-id"), nil
	case ctx.IsSet("volume-name"):
		name := ctx.String("volume-name")
		id, err := servers.IDFromName(sc, name)
		if err != nil {
			return "", fmt.Errorf("Error converting name [%s] to ID: %s", name, err)
		}
		return id, nil
	default:
		return "", lib.ErrMissingFlag{Msg: "One of either --volume-id or --volume-name must be provided."}
	}
}
