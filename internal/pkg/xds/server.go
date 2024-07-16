package xds

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

// defaultServerPort is the default port for the gRPC server.
 const defaultServerPort string = "15010"

// Run starts the gRPC server and the controllers.
func Run(ctx context.Context) error {
	var routinesGroup errgroup.Group
	grpcServer := grpc.NewServer()
	adsServerImpl := &adsServer{}

	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, adsServerImpl)

	listener, err := net.Listen("tcp", fmt.Sprint(":", defaultServerPort))
	if err != nil {
		return fmt.Errorf("creating TCP listener: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	routinesGroup.Go(func() error {
		defer cancel()
		log.Println("Running MCP RPC server")
		return grpcServer.Serve(listener)
	})

	routinesGroup.Go(func() error {
		return nil
	})
	
	routinesGroup.Go(func() error {
		defer log.Print("MCP gRPC server was shut down")
		<-ctx.Done()
		grpcServer.GracefulStop()
		return nil
	})

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	loop:
		for {
			select {
				case <-ctx.Done() :
					adsServerImpl.closeSubscribers()
					break loop

				case <-ticker.C :
					log.Println("Pushing to subscribers")
					if err := adsServerImpl.pushToSubscribers(); err != nil {
						log.Print("Error pushing to subscribers: ", err)
					}
			}
		}

		return routinesGroup.Wait()
}