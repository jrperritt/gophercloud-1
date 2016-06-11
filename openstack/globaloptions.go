package openstack

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gophercloud/cli/lib"
	"github.com/gophercloud/cli/util"

	"github.com/codegangsta/cli"
	"gopkg.in/ini.v1"
)

type GlobalOption struct {
	name     string
	value    interface{}
	from     string
	validate func() error
}

func (p GlobalOption) Name() string {
	return p.name
}

func (p GlobalOption) Value() interface{} {
	return p.value
}

func (p GlobalOption) From() string {
	return p.from
}

type GlobalOptions struct {
	username     string
	password     string
	authTenantID string
	authToken    string
	authURL      string
	region       string
	profile      string
	outputFormat string
	noCache      bool
	noHeader     bool
	logLevel     string
	logger       *logrus.Logger
	cliContext   *cli.Context
	have         map[string]GlobalOption
	want         []GlobalOption
}

// ParseGlobalOptions satisfies the Provider.ParseGlobalOptions method
func (o *GlobalOptions) ParseGlobalOptions() error {
	// we may get multiple errors while trying to handle the global options
	// so we'll try to return all of them at once, instead of returning just one,
	// only return a different one after that one's been rectified.
	multiErr := make(lib.MultiError, 0)

	// for each source where a user could provide a global option,
	// parse the options from that source. sources will be parsed in the order
	// in which they appear in the Sources method
	for _, source := range o.Sources() {
		if parseOptions := o.MethodsMap()[source]; parseOptions != nil {
			err := parseOptions()
			if err != nil {
				multiErr = append(multiErr, err)
			}
		}
	}

	// after the global options have been parsed, run each global option's
	// validation function, if it exists
	err := o.Validate()
	if err != nil {
		multiErr = append(multiErr, err)
	}

	if len(multiErr) > 0 {
		return multiErr
	}

	err = o.Set()
	if err != nil {
		return err
	}

	return nil
}

func (o *GlobalOptions) Init() error {
	o.want = []GlobalOption{
		{name: "username"},
		{name: "password"},
		{name: "auth-tenant-id"},
		{name: "auth-token"},
		{name: "auth-url"},
		{name: "region"},
		{name: "profile"},
	}

	for _, d := range o.Defaults() {
		o.want = append(o.want, d.(GlobalOption))
	}

	return nil
}

func (o GlobalOptions) Sources() []string {
	return []string{
		"commandline",
		"configfile",
		"envvar",
	}
}

func (o GlobalOptions) Defaults() []lib.GlobalOptioner {
	return []lib.GlobalOptioner{
		GlobalOption{"output", "table", "default", o.ValidateOutputInputParam},
		GlobalOption{"no-cache", false, "default", nil},
		GlobalOption{"no-header", false, "default", nil},
		GlobalOption{"log", "", "default", o.ValidateLogInputParam},
	}
}

func (o GlobalOptions) MethodsMap() map[string]func() error {
	return map[string]func() error{
		"commandline": o.ParseCommandLineOptions,
		"configfile":  o.ParseConfigFileOptions,
		"envvar":      o.ParseEnvVarOptions,
	}
}

func (o GlobalOptions) Validate() error {
	errs := make(lib.MultiError, 0)
	ds := o.Defaults()
	for _, d := range ds {
		inputParam := d.(GlobalOption)
		if inputParam.validate != nil {
			err := inputParam.validate()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// error if the user didn't provide an auth URL
	if _, ok := o.have["auth-url"]; !ok || o.have["auth-url"].value == "" {
		return fmt.Errorf("You must provide an authentication endpoint")
	}

	return nil
}

func (o GlobalOptions) Set() error {
	var level logrus.Level
	switch strings.ToLower(o.logLevel) {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	default:
		level = 0
	}
	o.logger = &logrus.Logger{
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

	return nil
}

func (o GlobalOptions) ParseCommandLineOptions() error {
	for i, opt := range o.want {
		if o.cliContext.IsSet(opt.name) {
			o.have[opt.name] = GlobalOption{value: o.cliContext.String(opt.name), from: "commandline"}
			o.want = append(o.want[:i], o.want[i+1:]...)
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

	for i, opt := range o.want {
		if v := section.Key(opt.name).String(); v != "" {
			o.have[opt.name] = GlobalOption{value: v, from: fmt.Sprintf("config file (profile: %s)", section.Name())}
			o.want = append(o.want[:i], o.want[i+1:]...)
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
	for i, opt := range o.want {
		if v := os.Getenv(strings.ToUpper(vars[opt.name])); v != "" {
			o.have[opt.name] = GlobalOption{value: v, from: "envvar"}
			o.want = append(o.want[:i], o.want[i+1:]...)
		}
	}
	return nil
}

func (o GlobalOptions) ValidateOutputInputParam() error {
	switch o.have["output"].value {
	case "json", "table", "":
		return nil
	default:
		return fmt.Errorf("Invalid value for `output` flag: '%s'. Options are: json, table.", o.outputFormat)
	}
}

func (o GlobalOptions) ValidateLogInputParam() error {
	switch o.have["log"].value {
	case "debug", "info", "":
		return nil
	default:
		return fmt.Errorf("Invalid value for `log` flag: %s. Valid options are: debug, info", o.logLevel)
	}
}

/*
func (o Logger) Set() error {

}
*/

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
