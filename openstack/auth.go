package openstack

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

type auth struct {
	logger        *logrus.Logger
	noCache       bool
	serviceType   string
	serviceClient *gophercloud.ServiceClient
	AuthOptions   *gophercloud.AuthOptions
	region        string
	urlType       gophercloud.Availability
	profile       string
}

func Authenticate() error {
	GC.GlobalOptions.authOptions.AllowReauth = true

	if !GC.GlobalOptions.noCache {
		err := AuthFromCache()
		if err != nil {
			return err
		}
	}

	if GC.ServiceClient == nil {
		err := AuthFromScratch()
		if err != nil {
			return err
		}
	}

	GC.ServiceClient.HTTPClient.Transport.(*LogRoundTripper).Logger = GC.GlobalOptions.logger
	return nil
}

func AuthFromScratch() error {
	serviceType := GC.Command.ServiceType()
	ao := GC.GlobalOptions.authOptions
	l := GC.GlobalOptions.logger

	l.Info("Authenticating from scratch.\n")

	l.Debugf("auth options: %+v\n", *ao)

	pc, err := openstack.NewClient(ao.IdentityEndpoint)
	if err != nil {
		return err
	}
	pc.HTTPClient = newHTTPClient(l)

	err = openstack.Authenticate(pc, *ao)
	if err != nil {
		return err
	}

	switch serviceType {
	case "compute":
		GC.ServiceClient, err = openstack.NewComputeV2(pc, gophercloud.EndpointOpts{
			Region:       GC.GlobalOptions.region,
			Availability: GC.GlobalOptions.urlType,
		})
		break
	case "files":
		GC.ServiceClient, err = openstack.NewObjectStorageV1(pc, gophercloud.EndpointOpts{
			Region:       GC.GlobalOptions.region,
			Availability: GC.GlobalOptions.urlType,
		})
		break
	case "block-storage":
		GC.ServiceClient, err = openstack.NewBlockStorageV1(pc, gophercloud.EndpointOpts{
			Region:       GC.GlobalOptions.region,
			Availability: GC.GlobalOptions.urlType,
		})
		break
	case "networking":
		GC.ServiceClient, err = openstack.NewNetworkV2(pc, gophercloud.EndpointOpts{
			Region:       GC.GlobalOptions.region,
			Availability: GC.GlobalOptions.urlType,
		})
		break
	case "orchestration":
		GC.ServiceClient, err = openstack.NewOrchestrationV1(pc, gophercloud.EndpointOpts{
			Region:       GC.GlobalOptions.region,
			Availability: GC.GlobalOptions.urlType,
		})
		break
	}
	if err != nil {
		return err
	}
	if GC.ServiceClient == nil {
		return fmt.Errorf("Unable to create service client: Unknown service type: %s\n", serviceType)
	}

	l.Debugf("Created %s service client: %+v", serviceType, GC.ServiceClient)
	GC.ServiceClient.UserAgent.Prepend(util.UserAgent)
	return nil
}

func AuthFromCache() error {
	ao := GC.GlobalOptions.authOptions
	l := GC.GlobalOptions.logger

	l.Info("Authenticating from cache")

	cache := GetCache()
	cacheKey := GetCacheKey()
	l.Infof("Looking in the cache for cache key: %s", cacheKey)
	credser, err := cache.GetCacheValue(cacheKey)

	if err == nil && credser != nil {
		creds := credser.(*CacheItem)
		l.Infof("Using token from cache: %s", creds.TokenID)
		pc, err := openstack.NewClient(ao.IdentityEndpoint)
		if err == nil {
			pc.UserAgent.Prepend(util.UserAgent)
			pc.TokenID = creds.GetToken()
			pc.HTTPClient = newHTTPClient(l)
			pc.ReauthFunc = func() error {
				return openstack.AuthenticateV3(pc, ao, gophercloud.EndpointOpts{})
			}
			GC.ServiceClient = &gophercloud.ServiceClient{
				ProviderClient: pc,
				Endpoint:       creds.ServiceEndpoint,
			}
			return nil
		}
	}
	return AuthFromScratch()
}

func GetCache() lib.Cacher {
	return &Cache{items: map[string]CacheItem{}}
}

func GetCacheKey() string {
	serviceType := GC.Command.ServiceType()
	ao := GC.GlobalOptions.authOptions
	l := GC.GlobalOptions.logger

	var usernameOrTenantID string
	switch {
	case ao.Username != "":
		usernameOrTenantID = ao.Username
	case ao.UserID != "":
		usernameOrTenantID = ao.UserID
	default:
		l.Debugf("Username nor User ID set in auth: %+v", ao)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s", usernameOrTenantID, ao.IdentityEndpoint, GC.GlobalOptions.region, serviceType, GC.GlobalOptions.urlType)
}

func StoreCredentials() error {
	newCacheValue := &CacheItem{
		TokenID:         GC.ServiceClient.TokenID,
		ServiceEndpoint: GC.ServiceClient.Endpoint,
	}
	cacheKey := GetCacheKey()
	GC.GlobalOptions.logger.Debugf("Setting cache key [%s] to: %s", cacheKey, newCacheValue)
	return GetCache().SetCacheValue(cacheKey, newCacheValue)
}
