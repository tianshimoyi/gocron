package options

import (
	"flag"
	genericoptions "github.com/x893675/gocron/pkg/server/options"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
	"strings"
)

type AgentOptions struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
}

func NewAgentOptions() *AgentOptions {
	s := &AgentOptions{GenericServerRunOptions: genericoptions.NewServerRunOptions()}
	return s
}

func (s *AgentOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("generic")
	s.GenericServerRunOptions.AddFlags(fs, s.GenericServerRunOptions)
	fs = fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		fs.AddGoFlag(fl)
	})
	return fss
}
