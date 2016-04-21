package objectcommands

import (
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
)

func CheckContainerExists(sc *gophercloud.ServiceClient, containerName string) error {
	containerRaw := containers.Get(sc, containerName)
	return containerRaw.Err
}
