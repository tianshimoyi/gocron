package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/x893675/gocron/cmd/gocron-server/app/options"
	"github.com/x893675/gocron/internal/apiserver"
	serverConfig "github.com/x893675/gocron/pkg/config"
	"github.com/x893675/gocron/pkg/utils/signals"
	"github.com/x893675/gocron/pkg/utils/term"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
)

func NewGoCronServerCommand() *cobra.Command {
	s := options.NewServerRunOptions()

	conf, err := serverConfig.TryLoadFromDisk()
	if err == nil {
		s = &options.ServerRunOptions{
			GenericServerRunOptions: s.GenericServerRunOptions,
			Config:                  conf,
		}
	}

	cmd := &cobra.Command{
		Use:  "gocron",
		Long: "The cron task schedule platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			if errs := s.Validate(); len(errs) != 0 {
				return utilerrors.NewAggregate(errs)
			}
			return Run(s, signals.SetupSignalHandler())
		},
		SilenceUsage: true,
	}

	fs := cmd.Flags()
	namedFlagSets := s.Flags()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
	})
	return cmd

}

func Run(s *options.ServerRunOptions, stopCh <-chan struct{}) error {
	apiServer, err := NewApiServer(s, stopCh)
	if err != nil {
		return err
	}
	err = apiServer.PrepareRun(stopCh)
	if err != nil {
		return nil
	}
	return apiServer.Run(stopCh)
}

func NewApiServer(s *options.ServerRunOptions, stopCh <-chan struct{}) (*apiserver.APIServer, error) {
	return nil, nil
}
