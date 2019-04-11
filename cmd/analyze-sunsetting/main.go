package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/supergiant/analyze-plugin-sunsetting/asset"
	"github.com/supergiant/analyze-plugin-sunsetting/cmd/analyze-sunsetting/server"
	"github.com/supergiant/analyze-plugin-sunsetting/info"
)

func main() {
	command := &cobra.Command{
		Use:          "analyze-sunsetting",
		Short:        "analyze-sunsetting plugin",
		RunE:         runCommand,
		SilenceUsage: true,
	}

	command.PersistentFlags().StringP(
		"grpc-api-port",
		"g",
		"9999", //TODO: make it configurable using configmap
		"tcp port where sunsetting plugin grpc server is serving")

	command.PersistentFlags().StringP(
		"rest-api-port",
		"r",
		"80", //TODO: make it configurable using configmap
		"tcp port where sunsetting plugin http server is serving")

	if err := command.Execute(); err != nil {
		log.Fatalf("\n%v\n", err)
	}
}

func runCommand(cmd *cobra.Command, _ []string) error {
	logger := &logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp:    true,
		},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel, //TODO: make it configurable
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	var mainLogger = logger.WithField("app", "sunsetting-plugin")

	mainLogger.Infof("%+v", info.Info())

	grpcAPIPort, err := cmd.Flags().GetString("grpc-api-port")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag grpc-api-port")
	}

	restAPIPort, err := cmd.Flags().GetString("rest-api-port")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag rest-api-port")
	}

	mainLogger.Infof("grpc-api-port: %v, rest-api-port: %v", grpcAPIPort, restAPIPort)

	//TODO: extract to separate component
	handler := func(w http.ResponseWriter, r *http.Request) {
		mainLogger.Warnf("%s", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/api/v1/info") {
			pi := info.Info()
			b, _ := json.Marshal(&pi)
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(b)
			if err != nil {
				mainLogger.Errorf("can't write info body: %+v", err)
			}
			return
		}

		fs := http.FileServer(asset.Assets)
		fs.ServeHTTP(w, r)
	}

	http.HandleFunc("/", handler)
	go func() {
		mainLogger.Fatal(http.ListenAndServe(":"+restAPIPort, nil))
	}()

	listener, err := net.Listen("tcp", ":"+grpcAPIPort)
	if err != nil {
		return err
	}

	var grpcLogger = mainLogger.WithField("component", "grpc_server")
	grpc_logrus.ReplaceGrpcLogger(grpcLogger)
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     12 * time.Hour,
			MaxConnectionAge:      0,
			MaxConnectionAgeGrace: 0,
			Time:                  5 * time.Minute,
			Timeout:               60 * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(grpcLogger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	reflection.Register(grpcServer)
	var pluginService = server.NewServer(mainLogger.WithField("component", "plugin_server"))
	proto.RegisterPluginServer(grpcServer, pluginService)

	return grpcServer.Serve(listener)
}
