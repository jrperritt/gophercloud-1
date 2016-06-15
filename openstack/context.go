package openstack

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud"
)

// Context satisfies the Provider interface
type Context struct {
	//logger *logrus.Logger
}

// Name satisfies the Provider.Name method
func (c *Context) Name() string {
	return "stack"
}

// NewGlobalOptionser satisfies the Provider.NewGlobalOptionser method
func (c *Context) NewGlobalOptionser(ctx *cli.Context) lib.GlobalOptionser {
	g := new(GlobalOptions)
	g.cliContext = ctx
	return g
}

// NewAuthenticater satisfies the Provider.NewAuthenticater method
func (c *Context) NewAuthenticater(globalOptionser lib.GlobalOptionser, serviceType string) lib.Authenticater {
	globalOptions := globalOptionser.(*GlobalOptions)

	fmt.Printf("auth-tenant-id: %s\n", globalOptions.authTenantID)

	return auth{
		authOptions: &gophercloud.AuthOptions{
			Username: globalOptions.username,
			UserID:   globalOptions.userID,
			Password: globalOptions.password,
			//ProjectID:        globalOptions.projectID,
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

func (c *Context) ResultsChannel() chan lib.Resulter {
	ch := make(chan lib.Resulter)
	return ch
}

// NewResultOutputter satisfies the Provider.NewResultOutputter method
func (c *Context) NewResultOutputter(globalOptionser lib.GlobalOptionser) lib.Outputter {
	globalOptions := globalOptionser.(*GlobalOptions)

	return output{
		//fields:   globalOptions.fields,
		noHeader: globalOptions.noHeader,
		format:   globalOptions.outputFormat,
		logger:   globalOptions.logger,
	}
}

func (c *Context) ErrExit1(err error) {
	panic(err)
}
