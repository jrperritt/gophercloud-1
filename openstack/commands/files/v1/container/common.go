package container

import (
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
)

func handleEmpty(c ContainerV1Command) error {
	allPages, err := objects.List(c.ServiceClient, c.name, nil).AllPages()
	if err != nil {
		return err
	}
	names, err := objects.ExtractNames(allPages)
	if err != nil {
		return err
	}

	for _, name := range names {
		objects.Delete(c.ServiceClient, c.name, name, nil)
	}

	header, err := containers.Get(c.ServiceClient, c.name).Extract()
	if err != nil {
		return err
	}

	if header.ObjectCount != "0" {
		handleEmpty(c)
	}

	return nil
}
