package app

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/x893675/gocron/cmd/gocron-server/app/options"
	"github.com/x893675/gocron/internal/apiserver"
	"github.com/x893675/gocron/internal/apiserver/service/task"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/client/notify"
	serverConfig "github.com/x893675/gocron/pkg/config"
	"github.com/x893675/gocron/pkg/utils/signals"
	"github.com/x893675/gocron/pkg/utils/term"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	cliflag "k8s.io/component-base/cli/flag"
	"net/http"
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
		return err
	}
	return apiServer.Run(stopCh)
}

func NewApiServer(s *options.ServerRunOptions, stopCh <-chan struct{}) (*apiserver.APIServer, error) {
	apiServer := &apiserver.APIServer{
		Config:    s.Config,
		JwtSecret: s.JwtSecret,
	}

	dbClient, err := database.NewDatabaseClient(s.DatabaseOptions, stopCh)
	if err != nil {
		return nil, err
	}
	apiServer.Db = dbClient
	informer := notify.NewInform(s.NotifyOptions)
	go informer.Run(stopCh)
	apiServer.TaskService = task.NewTaskService(dbClient, informer)
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.GenericServerRunOptions.BindAddress, s.GenericServerRunOptions.InsecurePort),
	}
	if s.GenericServerRunOptions.SecurePort != 0 {
		certificate, err := tls.LoadX509KeyPair(s.GenericServerRunOptions.TlsCertFile, s.GenericServerRunOptions.TlsPrivateKey)
		if err != nil {
			return nil, err
		}
		server.TLSConfig.Certificates = []tls.Certificate{certificate}
	}
	apiServer.Server = server

	return apiServer, nil
}
