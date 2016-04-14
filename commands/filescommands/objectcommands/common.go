package objectcommands

import (
	"github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud"
	"github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud/rackspace/objectstorage/v1/containers"
)

func CheckContainerExists(sc *gophercloud.ServiceClient, containerName string) error {
	containerRaw := containers.Get(sc, containerName)
	return containerRaw.Err
}
