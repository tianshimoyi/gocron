package agent

import (
	"context"
	"fmt"
	"github.com/x893675/gocron/cmd/gocron-agent/app/options"
	"github.com/x893675/gocron/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
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

	return server.Serve(l)

}
