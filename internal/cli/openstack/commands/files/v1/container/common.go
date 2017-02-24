package container

import (
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
)

func handleEmpty(c ContainerV1Command, container string) error {
	allPages, err := objects.List(c.ServiceClient(), container, nil).AllPages()
	if err != nil {
		return err
	}
	names, err := objects.ExtractNames(allPages)
	if err != nil {
		return err
	}

	for _, name := range names {
		_, err := objects.Delete(c.ServiceClient(), container, name, nil).Extract()
		if err != nil {
			return err
		}
	}

	header, err := containers.Get(c.ServiceClient(), c.name).Extract()
	if err != nil {
		return err
	}

	if header.ObjectCount != 0 {
		handleEmpty(c, container)
	}

	return nil
}
