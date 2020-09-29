package options

import (
	"flag"
	"github.com/x893675/gocron/pkg/config"
	genericoptions "github.com/x893675/gocron/pkg/server/options"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"strings"
)

type ServerRunOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
	*config.Config
	JwtSecret            string
	SkylineUrl           string
	SkylineAdminRoleName string
}

func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		Config:                  config.New(),
		JwtSecret:               "",
		SkylineUrl:              "",
		SkylineAdminRoleName:    "admin",
	}

	return s
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("generic")
	fs.StringVar(&s.JwtSecret, "jwt-secret", s.JwtSecret, "jwt secret for authenticate")
	fs.StringVar(&s.SkylineUrl, "skyline-url", s.SkylineUrl, "skyline url for authnz")
	fs.StringVar(&s.SkylineAdminRoleName, "skyline-admin-role-name", s.SkylineAdminRoleName, "skyline authnz admin role name, default is admin")
	s.GenericServerRunOptions.AddFlags(fs, s.GenericServerRunOptions)
	s.DatabaseOptions.AddFlags(fss.FlagSet("db"), s.DatabaseOptions)
	s.NotifyOptions.AddFlags(fss.FlagSet("notify"))
	fs = fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		fs.AddGoFlag(fl)
	})
	return fss
}
