package openstack

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

var (
	_ lib.AuthFromCacher = new(auth)
)

type auth struct {
	logger        *logrus.Logger
	noCache       bool
	serviceType   string
	serviceClient *gophercloud.ServiceClient
	authOptions   *gophercloud.AuthOptions
	region        string
	urlType       gophercloud.Availability
	profile       string
}

// Authenticate satisfies the Provider.Authenticate method
func (a *auth) Authenticate() (*gophercloud.ServiceClient, error) {
	var client *gophercloud.ServiceClient
	var err error

	a.authOptions.AllowReauth = true

	if authFromCacher, ok := interface{}(a).(lib.AuthFromCacher); ok {
		client, err = authFromCacher.AuthFromCache()
		if err != nil {
			return nil, err
		}
	}

	if client == nil {
		client, err = a.AuthFromScratch()
		if err != nil {
			return nil, err
		}
	}

	client.HTTPClient.Transport.(*LogRoundTripper).Logger = a.logger
	a.serviceClient = client
	a.logger.Debug("Returning from Authenticate...")
	return client, nil
}

func (a *auth) AuthFromScratch() (*gophercloud.ServiceClient, error) {
	a.logger.Info("Authenticating from scratch.\n")
	a.urlType = gophercloud.AvailabilityPublic

	a.logger.Debugf("auth options: %+v\n", *a.authOptions)

	pc, err := openstack.NewClient(a.authOptions.IdentityEndpoint)
	if err != nil {
		return nil, err
	}
	pc.HTTPClient = newHTTPClient(a.logger)

	err = openstack.Authenticate(pc, *a.authOptions)
	if err != nil {
		return nil, err
	}

	//a.logger.Debugf("provider client: %+v\n", pc)
	//a.logger.Debugf("a: %+v\n", a)

	var sc *gophercloud.ServiceClient
	switch a.serviceType {
	case "compute":
		sc, err = openstack.NewComputeV2(pc, gophercloud.EndpointOpts{
			Region:       a.region,
			Availability: a.urlType,
		})
		break
	case "files":
		sc, err = openstack.NewObjectStorageV1(pc, gophercloud.EndpointOpts{
			Region:       a.region,
			Availability: a.urlType,
		})
		break
	case "block-storage":
		sc, err = openstack.NewBlockStorageV1(pc, gophercloud.EndpointOpts{
			Region:       a.region,
			Availability: a.urlType,
		})
		break
	case "networking":
		sc, err = openstack.NewNetworkV2(pc, gophercloud.EndpointOpts{
			Region:       a.region,
			Availability: a.urlType,
		})
		break
	case "orchestration":
		sc, err = openstack.NewOrchestrationV1(pc, gophercloud.EndpointOpts{
			Region:       a.region,
			Availability: a.urlType,
		})
		break
	}
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, fmt.Errorf("Unable to create service client: Unknown service type: %s\n", a.serviceType)
	}

	a.logger.Debugf("Created %s service client: %+v", a.serviceType, sc)
	sc.UserAgent.Prepend(util.UserAgent)
	return sc, nil
}

func (a *auth) AuthFromCache() (*gophercloud.ServiceClient, error) {
	a.logger.Info("Authenticating from cache")

	switch a.urlType {
	case "":
		a.logger.Info("Using public endpoint")
		a.urlType = gophercloud.AvailabilityPublic
	}

	if a.noCache {
		return a.AuthFromScratch()
	}

	cache := a.GetCache()
	cacheKey := a.GetCacheKey()
	a.logger.Infof("Looking in the cache for cache key: %s", cacheKey)
	// get the value from the cache
	credser, err := cache.GetCacheValue(cacheKey)

	if err == nil && credser != nil {
		creds := credser.(*CacheItem)
		// we successfully retrieved a value from the cache
		a.logger.Infof("Using token from cache: %s", creds.TokenID)
		pc, err := openstack.NewClient(a.authOptions.IdentityEndpoint)
		if err == nil {
			pc.TokenID = creds.GetToken()
			pc.ReauthFunc = func() error {
				return openstack.AuthenticateV3(pc, a.authOptions, gophercloud.EndpointOpts{Availability: a.urlType})
			}
			pc.UserAgent.Prepend(util.UserAgent)
			pc.HTTPClient = newHTTPClient(a.logger)
			return &gophercloud.ServiceClient{
				ProviderClient: pc,
				Endpoint:       creds.ServiceEndpoint,
			}, nil
		}
	}
	// if there was an error accessing the cache or there was nothing in the cache,
	// authenticate from scratch
	return a.AuthFromScratch()
}

func (a *auth) GetCache() lib.Cacher {
	return &Cache{items: map[string]CacheItem{}}
}

// CacheKey returns the cache key formed from the user's authentication credentials.
func (a *auth) GetCacheKey() string {
	var usernameOrTenantID string
	switch {
	case a.authOptions.Username != "":
		usernameOrTenantID = a.authOptions.Username
	case a.authOptions.TenantID != "":
		usernameOrTenantID = a.authOptions.TenantID
	default:
		return ""
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s", usernameOrTenantID, a.authOptions.IdentityEndpoint, a.region, a.serviceType, a.urlType)
}

// StoreCredentials caches the users auth credentials if available and the `no-cache`
// flag was not provided.
func (a *auth) StoreCredentials() error {
	// if serviceClient is nil, the HTTP request for the command didn't get sent.
	// don't set cache if the `no-cache` flag is provided
	if a.noCache {
		return nil
	}

	newCacheValue := &CacheItem{
		TokenID:         a.serviceClient.TokenID,
		ServiceEndpoint: a.serviceClient.Endpoint,
	}
	// get the cache key
	cacheKey := a.GetCacheKey()

	a.logger.Debugf("Setting cache key [%s] to: %s", cacheKey, newCacheValue)

	// set the cache value to the current values
	return a.GetCache().SetCacheValue(cacheKey, newCacheValue)
}
