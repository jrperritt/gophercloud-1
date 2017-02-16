package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/lib/interfaces"
	"github.com/gophercloud/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

type authopts struct {
	cmd     interfaces.Commander
	region  string
	gao     *gophercloud.AuthOptions
	nocache bool
	urltype gophercloud.Availability
}

// auth authenticates a user against an endpoint
func auth(ao *authopts) (sc *gophercloud.ServiceClient, err error) {
	ao.gao.AllowReauth = true

	if !ao.nocache {
		sc, err = AuthFromCache(ao)
		if err != nil {
			return nil, err
		}
	}

	if sc == nil {
		sc, err = AuthFromScratch(ao)
		if err != nil {
			return nil, err
		}
	}

	sc.ReauthFunc = func() error {
		return openstack.AuthenticateV3(sc.ProviderClient, ao.gao, gophercloud.EndpointOpts{})
	}

	return sc, err
}

func AuthFromScratch(ao *authopts) (sc *gophercloud.ServiceClient, err error) {
	serviceType := ao.cmd.ServiceType()
	serviceVersion := ao.cmd.ServiceVersion()

	lib.Log.Debugln("Authenticating from scratch.\n")

	lib.Log.Debugf("auth options: %+v\n", *ao)

	pc, err := openstack.NewClient(ao.gao.IdentityEndpoint)
	if err != nil {
		return nil, err
	}
	pc.HTTPClient = newHTTPClient()
	pc.UserAgent.Prepend(util.UserAgent)

	err = openstack.Authenticate(pc, *ao.gao)
	if err != nil {
		return nil, err
	}

	sc, err = ao.cmd.ServiceClientFunc()(pc, gophercloud.EndpointOpts{
		Region:       ao.region,
		Availability: ao.urltype,
	})
	if err != nil {
		return sc, err
	}

	if sc == nil {
		return nil, fmt.Errorf("Unable to create service client: Unknown service type and version: %s %s", serviceType, serviceVersion)
	}

	lib.Log.Debugf("Created %s service client: %+v", serviceType, sc)

	return sc, nil
}

func AuthFromCache(ao *authopts) (sc *gophercloud.ServiceClient, err error) {
	lib.Log.Debugln("Authenticating from cache")

	cache := GetCache()
	cacheKey := GetCacheKey(ao)
	lib.Log.Debugf("Looking in the cache for cache key: %s", cacheKey)
	creds, err := cache.GetCacheValue(cacheKey)

	if err == nil && creds != nil {
		lib.Log.Debugf("Using token from cache: %s", creds.TokenID)
		pc, err := openstack.NewClient(ao.gao.IdentityEndpoint)
		if err == nil {
			pc.TokenID = creds.GetToken()
			pc.HTTPClient = newHTTPClient()
			sc = &gophercloud.ServiceClient{
				ProviderClient: pc,
				Endpoint:       creds.ServiceEndpoint,
			}
			sc.UserAgent.Prepend(util.UserAgent)

			ok, err := tokens.Validate(sc, pc.TokenID)
			if err == nil && ok {
				return sc, nil
			}
		}
	}
	return AuthFromScratch(ao)
}

// GetCache retreives the cache
func GetCache() *Cache {
	return &Cache{items: map[string]CacheItem{}}
}

// GetCacheKey retreives a cache key
func GetCacheKey(ao *authopts) string {
	var usernameOrTenantID string
	switch {
	case ao.gao.Username != "":
		usernameOrTenantID = ao.gao.Username
	case ao.gao.UserID != "":
		usernameOrTenantID = ao.gao.UserID
	default:
		lib.Log.Debugf("Username nor User ID set in auth: %+v", ao)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s", usernameOrTenantID, ao.gao.IdentityEndpoint, ao.region, ao.cmd.ServiceType(), ao.urltype)
}

func cachecreds(ao *authopts, sc *gophercloud.ServiceClient) error {
	ep := sc.Endpoint
	if rb := sc.ResourceBase; rb != "" {
		ep = rb
	}
	newCacheValue := &CacheItem{
		TokenID:         sc.TokenID,
		ServiceEndpoint: ep,
	}
	cacheKey := GetCacheKey(ao)
	lib.Log.Debugf("Setting cache key [%s] to: %s", cacheKey, newCacheValue)
	return GetCache().SetCacheValue(cacheKey, newCacheValue)
}
