package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gophercloud/cli/util"
	"github.com/gophercloud/cli/vendor/github.com/Sirupsen/logrus"

	"github.com/gophercloud/cli/vendor/github.com/codegangsta/cli"
	"github.com/gophercloud/cli/vendor/gopkg.in/ini.v1"
)

type InputParam struct {
	Name     string
	Value    interface{}
	From     string
	Validate func() error
}

type GlobalOptions struct {
	Username     string
	Password     string
	AuthTenantID string
	AuthToken    string
	AuthURL      string
	Region       string
	Profile      string
	Output       string
	NoCache      bool
	NoHeader     bool
	Logger       *logrus.Logger
	cliContext   *cli.Context
	have         map[string]InputParam
	want         []string
}

func (o *GlobalOptions) Init() error {
	o.want = o.Options()
	o.have = o.Defaults()
	return nil
}

func (o GlobalOptions) Sources() []string {
	return []string{
		"commandline",
		"configfile",
		"envvar",
	}
}

func (o GlobalOptions) GetGlobalOptions() []GetGlobalOptioner {
	p := []InputParam{
		InputParam{Name: "username"},
		InputParam{Name: "password"},
		InputParam{Name: "auth-tenant-id"},
		InputParam{Name: "auth-token"},
		InputParam{Name: "auth-url"},
		InputParam{Name: "region"},
		InputParam{Name: "profile"},
	}
	return append(p, o.Defaults())
}

func (o GlobalOptions) Defaults() []InputParam {
	return []InputParam{
		InputParam{"output", "table", "default", o.ValidateOutputInputParam},
		InputParam{"no-cache", false, "default"},
		InputParam{"no-header", false, "default"},
		InputParam{"log", "", "default", o.ValidateLogInputParam},
	}
}

func (o GlobalOptions) MethodsMap() map[string]func() error {
	return map[string]func() error{
		"commandline": o.ParseCommandLineOptions,
		"configfile":  o.ParseConfigFileOptions,
		"envvar":      o.ParseEnvVarOptions,
	}
}

func (o GlobalOptions) Validate() []error {
	errs := make([]error)
	for _, inputParam := range o.Defaults() {
		if inputParam.Validate != nil {
			err := inputParam.Validate()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// error if the user didn't provide an auth URL
	if _, ok := have["auth-url"]; !ok || have["auth-url"].Value == "" {
		return nil, fmt.Errorf("You must provide an authentication endpoint")
	}

}

func (o GlobalOptions) Set() error {
	var level logrus.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	default:
		level = ""
	}
	o.Logger = &logrus.Logger{
		Out:       o.cliContext.App.Writer,
		Formatter: &logrus.TextFormatter{},
		Level:     level,
	}

	/*
		haveString := ""
		for k, v := range have {
			haveString += fmt.Sprintf("%s: %s (from %s)\n", k, v.Value, v.From)
		}
		ctx.logger.Infof("Global Options:\n%s\n", haveString)
	*/
}

func (o GlobalOptions) ParseCommandLineOptions() error {
	for i, opt := range o.Want {
		if o.cliContext.IsSet(opt) {
			o.Have[opt] = InputParam{Value: o.cliContext.String(opt), From: "commandline"}
			o.Want = append(o.Want[:i], o.Want[i+1:]...)
		}
	}
	return nil
}

func (o GlobalOptions) ParseConfigFileOptions() error {
	profile := o.cliContext.String("profile")
	section, err := ProfileSection(profile)
	if err != nil {
		return err
	}

	if section == nil {
		return nil
	}

	for i, opt := range o.Want {
		if v := section.Key(opt).String(); v != "" {
			o.Have[opt] = InputParam{Value: v, From: fmt.Sprintf("config file (profile: %s)", section.Name())}
			o.Want = append(o.Want[:i], o.Want[i+1:]...)
		}
	}

	return nil
}

func (o GlobalOptions) ParseEnvVarOptions() error {
	vars := map[string]string{
		"username": "OS_USERNAME",
		"password": "OS_PASSWORD",
		"auth-url": "OS_AUTH_URL",
		"region":   "OS_REGION_NAME",
	}
	for i, opt := range o.Want {
		if v := os.Getenv(strings.ToUpper(vars[opt])); v != "" {
			o.Have[opt] = InputParam{Value: v, From: "envvar"}
			o.Want = append(o.Want[:i], o.Want[i+1:]...)
		}
	}
	return nil
}

func (o GlobalOptions) ValidateOutputInputParam() error {
	switch o.Have["output"] {
	case "json", "table", "":
		return nil
	default:
		return fmt.Errorf("Invalid value for `output` flag: '%s'. Options are: json, table.", outputFormat)
	}
}

func (o GlobalOptions) ValidateLogInputParam() error {
	switch o.Have["log"] {
	case "debug", "info", "":
		return nil
	default:
		return fmt.Errorf("Invalid value for `log` flag: %s. Valid options are: debug, info", logLevel)
	}
}

func (o Logger) Set() error {

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
