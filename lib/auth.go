package lib

import "github.com/gophercloud/gophercloud"

// Authenticater is implemented by types that can authenticate a user
type Authenticater interface {
	Authenticate() (*gophercloud.ServiceClient, error)
	AuthFromScratch() (*gophercloud.ServiceClient, error)
	//SupportedServices() []string
}

// AuthFromCacher is implemented by types that can authenticate
// a user from a cache
type AuthFromCacher interface {
	AuthFromCache() (*gophercloud.ServiceClient, error)
	GetCache() Cacher
	GetCacheKey() string
	StoreCredentials() error
}
