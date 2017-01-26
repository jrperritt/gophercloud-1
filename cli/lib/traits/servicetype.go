package traits

import (
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/openstack"
)

type Computeable struct{}

func (_ Computeable) ServiceType() string {
	return "compute"
}

type ComputeV2able struct {
	Computeable
}

func (_ ComputeV2able) ServiceVersion() string {
	return "v2"
}

func (_ ComputeV2able) ServiceClientFunc() interfaces.ServiceClientFunc {
	return openstack.NewComputeV2
}

type Filesable struct{}

func (_ Filesable) ServiceType() string {
	return "files"
}

type FilesV1able struct {
	Filesable
}

func (_ FilesV1able) ServiceVersion() string {
	return "v1"
}

func (_ FilesV1able) ServiceClientFunc() interfaces.ServiceClientFunc {
	return openstack.NewObjectStorageV1
}

type BlockStorageable struct{}

func (_ BlockStorageable) ServiceType() string {
	return "block-storage"
}

type BlockStorageV2able struct {
	BlockStorageable
}

func (_ BlockStorageV2able) ServiceVersion() string {
	return "v2"
}

func (_ BlockStorageV2able) ServiceClientFunc() interfaces.ServiceClientFunc {
	return openstack.NewBlockStorageV2
}

type Networkingable struct{}

func (_ Networkingable) ServiceType() string {
	return "networking"
}

type NetworkingV2able struct {
	Networkingable
}

func (_ NetworkingV2able) ServiceVersion() string {
	return "v2"
}

func (_ NetworkingV2able) ServiceClientFunc() interfaces.ServiceClientFunc {
	return openstack.NewNetworkV2
}
