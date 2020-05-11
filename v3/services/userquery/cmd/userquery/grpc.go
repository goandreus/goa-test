package main

import (
	"context"
	"net"
	"net/url"
	"sync"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	userquerypb "gitlab.com/wiserskills/v3/services/userquery/gen/grpc/userquery/pb"
	userquerysvr "gitlab.com/wiserskills/v3/services/userquery/gen/grpc/userquery/server"
	log "gitlab.com/wiserskills/v3/services/userquery/gen/log"
	userquery "gitlab.com/wiserskills/v3/services/userquery/gen/userquery"
	grpcmdlwr "goa.design/goa/v3/grpc/middleware"
	"goa.design/goa/v3/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// handleGRPCServer starts configures and starts a gRPC server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleGRPCServer(ctx context.Context, u *url.URL, userqueryEndpoints *userquery.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = logger
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to gRPC requests and
	// responses.
	var (
		userqueryServer *userquerysvr.Server
	)
	{
		userqueryServer = userquerysvr.New(userqueryEndpoints, nil)
	}

	// Initialize gRPC server with the middleware.
	srv := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcmdlwr.UnaryRequestID(),
			grpcmdlwr.UnaryServerLog(adapter),
		),
	)

	// Register the servers.
	userquerypb.RegisterUserqueryServer(srv, userqueryServer)

	for svc, info := range srv.GetServiceInfo() {
		for _, m := range info.Methods {
			logger.Infof("serving gRPC method %s", svc+"/"+m.Name)
		}
	}

	// Register the server reflection service on the server.
	// See https://grpc.github.io/grpc/core/md_doc_server-reflection.html.
	reflection.Register(srv)

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start gRPC server in a separate goroutine.
		go func() {
			lis, err := net.Listen("tcp", u.Host)
			if err != nil {
				errc <- err
			}
			logger.Infof("gRPC server listening on %q", u.Host)
			errc <- srv.Serve(lis)
		}()

		<-ctx.Done()
		logger.Infof("shutting down gRPC server at %q", u.Host)
		srv.Stop()
	}()
}
