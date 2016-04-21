package lib

import "github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"

type Authenticater interface {
	AuthFromScratch() (*gophercloud.ServiceClient, error)
	Credentials() error
}

type AuthFromCacher interface {
	AuthFromCache() (*gophercloud.ServiceClient, error)
	GetCache() Cacher
}
