/*
Copyright © 2023 Dhiraj Bokde dhirajsb@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	gwHost   = "localhost"
	grpcHost = "localhost"
	gwPort   = 8081
	grpcPort = 8080

	// serveCmd represents the serve command
	gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Starts the gRPC gateway server",
		Long: `This command launches gRPC gateway server.

The server connects to a runnning instance of gRPC server. '`,
		RunE: runGrpcGatewayServer,
	}
)

func runGrpcGatewayServer(cmd *cobra.Command, args []string) error {
	glog.Info("starting server...")

	// Create a channel to receive signals
	signalChannel := make(chan os.Signal, 1)

	// Notify the channel on SIGINT (Ctrl+C) and SIGTERM signals
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, gwPort))
	if err != nil {
		log.Fatalf("server listen failed: %v", err)
	}

	// gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcGatewayServer := createGrpcGatewayMux(ctx)

	// start cmux listeners
	g := new(errgroup.Group)
	g.Go(func() error {
		glog.Info("starting gRPC gateway server...")
		return grpcGatewayServer.Serve(listener)
	})

	go func() {
		err := g.Wait()
		// error starting server
		if err != nil || err != http.ErrServerClosed || err != grpc.ErrServerStopped || err != cmux.ErrServerClosed {
			glog.Errorf("server listener error: %v", err)
		}
		signalChannel <- syscall.SIGINT
	}()

	// Wait for a signal
	receivedSignal := <-signalChannel
	glog.Infof("received signal: %s\n", receivedSignal)

	// Perform cleanup or other graceful shutdown actions here
	glog.Info("shutting down services...")
	_ = grpcGatewayServer.Shutdown(context.Background())

	glog.Info("shutdown!")
	return nil
}

func withLogging(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// gather information about request and log it
		uri := r.URL.String()
		method := r.Method

		fmt.Printf("%v: [%s] %s\n", time.Now().Local(), method, uri)

		// call the original http.Handler we're wrapping
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func createGrpcGatewayMux(ctx context.Context) *http.Server {
	mux := runtime.NewServeMux(
		runtime.WithUnescapingMode(runtime.UnescapingModeAllExceptReserved),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	lopts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadReceived, logging.PayloadSent, logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	grpcServerAddr := fmt.Sprintf("%s:%d", host, grpcPort)
	conn, err := grpc.DialContext(
		ctx,
		grpcServerAddr,
		grpc.WithReturnConnectionError(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(InterceptorLogger(logger), lopts...),
		),
	)
	if err != nil {
		log.Fatalf("Error dialing connection to grpc server %s: %v", grpcServerAddr, err)
	}

	err = proto.RegisterMetadataStoreServiceHandler(ctx, mux, conn)
	if err != nil {
		log.Fatal(err)
	}

	return &http.Server{
		Handler: withLogging(mux),
	}
}

func init() {
	rootCmd.AddCommand(gatewayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	gatewayCmd.Flags().StringVarP(&gwHost, "hostname", "n", gwHost, "Server listen hostname")
	gatewayCmd.Flags().IntVarP(&gwPort, "port", "p", gwPort, "Server listen port")
	gatewayCmd.Flags().StringVarP(&grpcHost, "grpc-hostname", "m", grpcHost, "Server listen hostname")
	gatewayCmd.Flags().IntVarP(&grpcPort, "grpc-port", "q", grpcPort, "Server listen port")
}
