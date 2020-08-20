package agent

import (
	"context"
	"fmt"
	"github.com/x893675/gocron/cmd/gocron-agent/app/options"
	"github.com/x893675/gocron/pkg/pb"
	"github.com/x893675/gocron/pkg/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"k8s.io/klog/v2"
	"net"
	"time"
)

type Server struct {
	Options *options.AgentOptions
}

func (s *Server) Run(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	return nil, nil
}

func (s *Server) Serve(stopCh <-chan struct{}) error {

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Options.GenericServerRunOptions.BindAddress, s.Options.GenericServerRunOptions.InsecurePort))
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 30 * time.Second,
			Time:              30 * time.Second,
			Timeout:           3 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	server := grpc.NewServer(opts...)
	pb.RegisterTaskServer(server, s)

	go func() {
		<-stopCh
		server.GracefulStop()
	}()

	klog.V(0).Info("gocron agent version is ", version.Version.String())
	return server.Serve(l)

}
