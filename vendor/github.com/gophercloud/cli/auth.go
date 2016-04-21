package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud"
	"github.com/gophercloud/cli/vendor/github.com/gophercloud/gophercloud/openstack"
)

type Auth struct {
	Logger        *logrus.Logger
	GlobalOptions *GlobalOptions
	AuthOptions   *gophercloud.AuthOptions
	Region        string
}

func (a *Auth) AuthFromScratch() (*gophercloud.ServiceClient, error) {
	a.Logger.Info("Not using cache; Authenticating from scratch.\n")

	ao := credsResult.AuthOpts
	region := credsResult.Region

	pc, err := openstack.AuthenticatedClient(*ao)
	if err != nil {
		switch err.(type) {
		case *tokens2.ErrNoPassword:
			return nil, errors.New("Please supply an API key.")
		}
		return nil, err
	}
	pc.HTTPClient = newHTTPClient()
	var sc *gophercloud.ServiceClient
	switch serviceType {
	case "compute":
		sc, err = openstack.NewComputeV2(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "object-store":
		sc, err = openstack.NewObjectStorageV1(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "blockstorage":
		sc, err = openstack.NewBlockStorageV1(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "network":
		sc, err = openstack.NewNetworkV2(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	case "orchestration":
		sc, err = openstack.NewOrchestrationV1(pc, gophercloud.EndpointOpts{
			Region:       region,
			Availability: urlType,
		})
		break
	}
	if err != nil {
		return nil, err
	}
	if sc == nil {
		return nil, fmt.Errorf("Unable to create service client: Unknown service type: %s\n", serviceType)
	}

	logger.Debugf("Created %s service client: %+v", serviceType, sc)
	sc.UserAgent.Prepend(util.UserAgent)
	return sc, nil
}

func (a Auth) AuthFromCache() (*gophercloud.ServiceClient, error) {
	err := a.Credentials()
	if err != nil {
		return nil, err
	}

	logMsg := "Using public endpoint"
	urlType := gophercloud.AvailabilityPublic
	a.Logger.Infoln(logMsg)

	if a.GlobalOptions.NoCache {
		return a.AuthFromScratch()
	}

	cache := a.GetCache()
	cacheKey := cache.GetKey()
	logger.Infof("Looking in the cache for cache key: %s\n", cacheKey)
	// get the value from the cache
	creds, err := cache.GetValue(cacheKey)
	// if there was an error accessing the cache or there was nothing in the cache,
	// authenticate from scratch
	if err == nil && creds != nil {
		// we successfully retrieved a value from the cache
		logger.Infof("Using token from cache: %s\n", creds.TokenID)
		pc, err := openstack.NewClient(ao.IdentityEndpoint)
		if err == nil {
			pc.TokenID = creds.TokenID
			pc.ReauthFunc = func() error {
				return openstack.AuthenticateV2(pc, ao)
			}
			pc.UserAgent.Prepend(util.UserAgent)
			pc.HTTPClient = newHTTPClient()
			return &gophercloud.ServiceClient{
				ProviderClient: pc,
				Endpoint:       creds.ServiceEndpoint,
			}, nil
		}
	} else {
		return authFromScratch(credsResult, serviceType, urlType, logger)
	}
}

func (a Auth) GetCache() Cacher {
	return &Cache{}
}

func (a Auth) Credentials() error {
	a.AuthOptions = &gophercloud.AuthOptions{
		AllowReauth:      true,
		IdentityEndpoint: a.GlobalOptions.AuthURL,
		Username:         a.GlobalOptions.Username,
		Password:         a.GlobalOptions.Password,
		TenantID:         a.GlobalOptions.AuthTenantID,
		TokenID:          a.GlobalOptions.AuthToken,
		//UserID: a.GlobalOptions.UserID,
	}

	a.Region = a.GlobalOptions.Region

	/*
		if logger != nil {
			haveString := ""
			for k, v := range have {
				haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.Value, v.From)
			}
			a.Logger.Infof("Authentication Credentials:\n%s\n", haveString)
		}
	*/

	return nil
}

var usernameAuthErrSlice = []string{"There are some required credentials that we couldn't find.",
	"Here's what we have:",
	"%s",
	"and here's what we're missing:",
	"%s",
	"",
	"You can set any of these credentials in the following ways:",
	"- Run `rack configure` to interactively create a configuration file,",
	"- Specify it in the command as a flag (--username, --api-key), or",
	"- Export it as an environment variable (RS_USERNAME, RS_API_KEY).",
	"",
}

var tenantIDAuthErrSlice = []string{"There are some required credentials that we couldn't find.",
	"Here's what we have:",
	"%s",
	"and here's what we're missing:",
	"%s",
	"",
	"You can set the missing credentials with command-line flags (--auth-token, --auth-tenant-id)",
	"",
}

// Err returns the custom error to print when authentication fails.
func Err(have map[string]Cred, want map[string]string, errMsg []string) error {
	haveString := ""
	for k, v := range have {
		haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.Value, v.From)
	}

	if len(want) > 0 {
		wantString := ""
		for k := range want {
			wantString += fmt.Sprintf("%s\n", k)
		}

		return fmt.Errorf(fmt.Sprintf(strings.Join(errMsg, "\n"), haveString, wantString))
	}

	return nil
}
