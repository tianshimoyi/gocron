package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/x893675/gocron/cmd/gocron-agent/app/options"
	"github.com/x893675/gocron/internal/agent"
	"github.com/x893675/gocron/pkg/utils/signals"
	"github.com/x893675/gocron/pkg/utils/term"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
)

func NewGoCronAgentCommand() *cobra.Command {
	s := options.NewAgentOptions()

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

func Run(s *options.AgentOptions, stopCh <-chan struct{}) error {

	srv := &agent.Server{
		Options: s,
	}
	return srv.Serve(stopCh)
}
