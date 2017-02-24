package openstack

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/internal/cli/lib"
	"github.com/gophercloud/gophercloud/internal/cli/util"

	"gopkg.in/ini.v1"
	cli "gopkg.in/urfave/cli.v1"
)

type GlobalOption struct {
	name  string
	value interface{}
	from  string
}

type GlobalOptions struct {
	authOptions *gophercloud.AuthOptions
	region      string
	urlType     gophercloud.Availability
	profile     string
	nocache     bool
	loglevel    uint
	have        map[string]GlobalOption
	want        []GlobalOption
	fields      []string
}

// globalopts sets the global context's global options
func globalopts(ctx *cli.Context) (gopts *GlobalOptions, err error) {
	gopts = new(GlobalOptions)

	gopts.want = []GlobalOption{
		{name: "username"},
		{name: "user-id"},
		{name: "password"},
		{name: "auth-tenant-id"},
		{name: "auth-token"},
		{name: "auth-url"},
		{name: "region"},
		{name: "profile"},
		{name: "log"},
	}
	gopts.have = make(map[string]GlobalOption)

	gopts.setdefaults()

	gopts.parseclopts(ctx)
	gopts.parseiniopts(ctx)
	gopts.parsevaropts()

	err = gopts.validate()
	if err != nil {
		return gopts, err
	}

	gopts.set()

	return gopts, nil
}

func (gopts *GlobalOptions) setdefaults() {
	gopts.urlType = gophercloud.AvailabilityPublic
}

func (gopts *GlobalOptions) validate() error {
	// error if the user didn't provide an auth URL
	if _, ok := gopts.have["auth-url"]; !ok || gopts.have["auth-url"].value == "" {
		return fmt.Errorf("You must provide an authentication endpoint")
	}
	return nil
}

func (gopts *GlobalOptions) set() error {

	gopts.authOptions = new(gophercloud.AuthOptions)
	var err error
	for name, opt := range gopts.have {
		switch name {
		case "username":
			gopts.authOptions.Username = opt.value.(string)
		case "user-id":
			gopts.authOptions.UserID = opt.value.(string)
		case "password":
			gopts.authOptions.Password = opt.value.(string)
		case "auth-tenant-id":
			gopts.authOptions.TenantID = opt.value.(string)
		case "auth-token":
			gopts.authOptions.TokenID = opt.value.(string)
		case "auth-url":
			gopts.authOptions.IdentityEndpoint = opt.value.(string)
		case "region":
			gopts.region = opt.value.(string)
		case "profile":
			gopts.profile = opt.value.(string)
		case "no-cache":
			switch t := opt.value.(type) {
			case string:
				gopts.nocache, err = strconv.ParseBool(t)
			case bool:
				gopts.nocache = t
			}
		case "log":
			switch opt.value.(string) {
			case "dev":
				gopts.loglevel = uint(lib.Dev)
			case "debug":
				gopts.loglevel = uint(lib.Debug)
			case "warn":
			case "info":
			}
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// parseclopts parses global flags
func (gopts *GlobalOptions) parseclopts(ctx *cli.Context) error {
	tmp := make([]GlobalOption, 0)

	for _, opt := range gopts.want {
		if ctx.GlobalIsSet(opt.name) {
			gopts.have[opt.name] = GlobalOption{value: ctx.GlobalString(opt.name), from: "commandline"}
			continue
		}
		tmp = append(tmp, opt)
	}
	gopts.want = tmp

	return nil
}

// parseiniopts parses and stores options from a profile in a config
// file
func (gopts *GlobalOptions) parseiniopts(ctx *cli.Context) error {
	profile := ctx.String("profile")
	section, err := ProfileSection(profile)
	if err != nil {
		return err
	}

	if section == nil {
		return nil
	}

	tmp := make([]GlobalOption, 0)
	for _, opt := range gopts.want {
		if v := section.Key(opt.name).String(); v != "" {
			gopts.have[opt.name] = GlobalOption{value: v, from: fmt.Sprintf("config file (profile: %s)", section.Name())}
			continue
		}
		tmp = append(tmp, opt)
	}
	gopts.want = tmp
	return nil
}

// parsevaropts parses global options stores in environment variables
func (gopts *GlobalOptions) parsevaropts() error {
	vars := map[string]string{
		"username":       "OS_USERNAME",
		"user-id":        "OS_USERID",
		"auth-tenant-id": "OS_TENANTID",
		"password":       "OS_PASSWORD",
		"auth-url":       "OS_AUTH_URL",
		"region":         "OS_REGION_NAME",
	}

	tmp := make([]GlobalOption, 0)
	for _, opt := range gopts.want {
		if v := os.Getenv(strings.ToUpper(vars[opt.name])); v != "" {
			gopts.have[opt.name] = GlobalOption{value: v, from: "envvar"}
			continue
		}
		tmp = append(tmp, opt)
	}
	gopts.want = tmp
	return nil
}

func ProfileSection(profile string) (*ini.Section, error) {
	dir, err := util.StackDir()
	if err != nil {
		return nil, nil
	}
	f := path.Join(dir, "config")
	cfg, err := ini.Load(f)
	if err != nil {
		return nil, nil
	}
	cfg.BlockMode = false
	section, err := cfg.GetSection(profile)
	if err != nil && profile != "" {
		return nil, fmt.Errorf("Invalid config file profile: %s\n", profile)
	}
	return section, nil
}
