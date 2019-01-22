package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main()  {
	command := &cobra.Command{
		Use:          "analyze-sunsetting",
		Short:        "analyze-sunsetting plugin",
		RunE:         runCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().StringP(
		"api-port",
		"p",
		"9999", //TODO: make it configurable using configmap
		"tcp port where sunsetting plugin grpc server is serving")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {

	apiPort, err := cmd.Flags().GetString("api-port")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag api-port")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) //TODO: make it configurable
	var mainLogger = logger.WithField("app", "sunsetting-plugin")

	listener, err := net.Listen("tcp", ":"+apiPort)
	if err != nil {
		return err
	}

	var grpcLogger = mainLogger.WithField("component", "grpc_server")
	grpc_logrus.ReplaceGrpcLogger(grpcLogger)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(grpcLogger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	reflection.Register(grpcServer)

	var pluginService = NewServer(mainLogger.WithField("component", "plugin_server"))

	proto.RegisterPluginServer(grpcServer, pluginService)


	return grpcServer.Serve(listener)
}
