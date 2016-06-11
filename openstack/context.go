package openstack

import (
	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud"
)

// Context satisfies the Provider interface
type Context struct {
	logger *logrus.Logger
}

// Name satisfies the Provider.Name method
func (c Context) Name() string {
	return "stack"
}

// NewGlobalOptionser satisfies the Provider.NewGlobalOptionser method
func (c Context) NewGlobalOptionser() lib.GlobalOptionser {
	return new(GlobalOptions)
}

// NewAuthenticater satisfies the Provider.NewAuthenticater method
func (c Context) NewAuthenticater(globalOptionser lib.GlobalOptionser, serviceType string) lib.Authenticater {
	globalOptions := globalOptionser.(*GlobalOptions)

	return auth{
		authOptions: &gophercloud.AuthOptions{
			Username:         globalOptions.username,
			Password:         globalOptions.password,
			TenantID:         globalOptions.authTenantID,
			TokenID:          globalOptions.authToken,
			IdentityEndpoint: globalOptions.authURL,
		},
		logger:      globalOptions.logger,
		noCache:     globalOptions.noCache,
		serviceType: serviceType,
		region:      globalOptions.region,
		profile:     globalOptions.profile,
	}
}

// NewResultOutputter satisfies the Provider.NewResultOutputter method
func (c Context) NewResultOutputter(globalOptionser lib.GlobalOptionser) lib.Outputter {
	globalOptions := globalOptionser.(*GlobalOptions)

	return output{
		//fields:   globalOptions.fields,
		noHeader: globalOptions.noHeader,
		format:   globalOptions.outputFormat,
		logger:   globalOptions.logger,
	}
}
