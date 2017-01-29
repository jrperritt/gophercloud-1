package openstack

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/cli/lib"
	"github.com/gophercloud/gophercloud/cli/util"

	"gopkg.in/ini.v1"
)

type GlobalOption struct {
	name     string
	value    interface{}
	from     string
	validate func() error
}

type GlobalOptions struct {
	authOptions  *gophercloud.AuthOptions
	region       string
	urlType      gophercloud.Availability
	profile      string
	outputFormat string
	noCache      bool
	logger       *logger
	have         map[string]GlobalOption
	want         []GlobalOption
	fields       []string
}

func SetGlobalOptions() error {
	GC.GlobalOptions = new(GlobalOptions)
	GC.GlobalOptions.want = []GlobalOption{
		{name: "username"},
		{name: "user-id"},
		{name: "password"},
		{name: "auth-tenant-id"},
		{name: "auth-token"},
		{name: "auth-url"},
		{name: "region"},
		{name: "profile"},
	}
	GC.GlobalOptions.have = make(map[string]GlobalOption)

	SetGlobalOptionsDefaults()

	ParseCommandLineOptions()
	ParseConfigFileOptions()
	ParseEnvVarOptions()

	setGlobalOptions()

	return nil
}

func GlobalOptionsDefaults() []GlobalOption {
	return []GlobalOption{
		GlobalOption{"no-cache", false, "default", nil},
	}
}

func SetGlobalOptionsDefaults() {
	for _, opt := range GlobalOptionsDefaults() {
		GC.GlobalOptions.have[opt.name] = opt
		GC.GlobalOptions.want = append(GC.GlobalOptions.want, opt)
	}
	GC.GlobalOptions.urlType = gophercloud.AvailabilityPublic
}

func Validate() error {
	errs := make(lib.MultiError, 0)
	for _, d := range GlobalOptionsDefaults() {
		if d.validate != nil {
			err := d.validate()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// error if the user didn't provide an auth URL
	if _, ok := GC.GlobalOptions.have["auth-url"]; !ok || GC.GlobalOptions.have["auth-url"].value == "" {
		return fmt.Errorf("You must provide an authentication endpoint")
	}

	return nil
}

func setGlobalOptions() error {
	l := new(logger)
	l.Logger = *log.New(GC.CommandContext.App.Writer, "", log.LstdFlags)
	GC.GlobalOptions.logger = l

	GC.GlobalOptions.authOptions = new(gophercloud.AuthOptions)
	var err error
	for name, opt := range GC.GlobalOptions.have {
		switch name {
		case "username":
			GC.GlobalOptions.authOptions.Username = opt.value.(string)
		case "user-id":
			GC.GlobalOptions.authOptions.UserID = opt.value.(string)
		case "password":
			GC.GlobalOptions.authOptions.Password = opt.value.(string)
		case "auth-tenant-id":
			GC.GlobalOptions.authOptions.TenantID = opt.value.(string)
		case "auth-token":
			GC.GlobalOptions.authOptions.TokenID = opt.value.(string)
		case "auth-url":
			GC.GlobalOptions.authOptions.IdentityEndpoint = opt.value.(string)
		case "region":
			GC.GlobalOptions.region = opt.value.(string)
		case "profile":
			GC.GlobalOptions.profile = opt.value.(string)
		case "output":
			GC.GlobalOptions.outputFormat = opt.value.(string)
		case "no-cache":
			switch t := opt.value.(type) {
			case string:
				GC.GlobalOptions.noCache, err = strconv.ParseBool(t)
			case bool:
				GC.GlobalOptions.noCache = t
			}
		/*case "no-header":
		switch t := opt.value.(type) {
		case string:
			GC.GlobalOptions.noHeader, err = strconv.ParseBool(t)
		case bool:
			GC.GlobalOptions.noHeader = t
		}*/
		case "debug":
			GC.GlobalOptions.logger.debug = true
		}
	}

	if err != nil {
		return err
	}

	switch GC.CommandContext.IsSet("fields") {
	case true:
		GC.GlobalOptions.fields = strings.Split(GC.CommandContext.String("fields"), ",")
	}

	return nil
}

func ParseCommandLineOptions() error {
	tmp := make([]GlobalOption, 0)

	for _, opt := range GC.GlobalOptions.want {
		if GC.CommandContext.GlobalIsSet(opt.name) {
			GC.GlobalOptions.have[opt.name] = GlobalOption{value: GC.CommandContext.GlobalString(opt.name), from: "commandline"}
			continue
		}
		tmp = append(tmp, opt)
	}
	GC.GlobalOptions.want = tmp

	return nil
}

func ParseConfigFileOptions() error {
	profile := GC.CommandContext.String("profile")
	section, err := ProfileSection(profile)
	if err != nil {
		return err
	}

	if section == nil {
		return nil
	}

	tmp := make([]GlobalOption, 0)
	for _, opt := range GC.GlobalOptions.want {
		if v := section.Key(opt.name).String(); v != "" {
			GC.GlobalOptions.have[opt.name] = GlobalOption{value: v, from: fmt.Sprintf("config file (profile: %s)", section.Name())}
			continue
		}
		tmp = append(tmp, opt)
	}
	GC.GlobalOptions.want = tmp
	return nil
}

func ParseEnvVarOptions() error {
	vars := map[string]string{
		"username":       "OS_USERNAME",
		"user-id":        "OS_USERID",
		"auth-tenant-id": "OS_TENANTID",
		"password":       "OS_PASSWORD",
		"auth-url":       "OS_AUTH_URL",
		"region":         "OS_REGION_NAME",
	}

	tmp := make([]GlobalOption, 0)
	for _, opt := range GC.GlobalOptions.want {
		if v := os.Getenv(strings.ToUpper(vars[opt.name])); v != "" {
			GC.GlobalOptions.have[opt.name] = GlobalOption{value: v, from: "envvar"}
			continue
		}
		tmp = append(tmp, opt)
	}
	GC.GlobalOptions.want = tmp
	return nil
}

func ProfileSection(profile string) (*ini.Section, error) {
	dir, err := util.RackDir()
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
