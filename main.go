package main

import (
	"context"
	"github.com/x893675/gocron/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

var ClientOptions = []grpc.DialOption{
	grpc.WithInsecure(),
	grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}),
	grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.Config{MaxDelay: 3 * time.Second}}),
}

func main() {
	endpoint := "192.168.234.137:8080"
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, endpoint, ClientOptions...)
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewTaskClient(conn)
	resp, err := c.Run(context.TODO(), &pb.TaskRequest{
		Command: "echo \"hello\"",
		Timeout: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
