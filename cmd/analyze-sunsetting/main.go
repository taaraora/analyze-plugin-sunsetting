package main

import (
	"context"
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
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/supergiant/analyze-plugin-sunsetting/asset"
	"github.com/supergiant/analyze-plugin-sunsetting/cmd/analyze-sunsetting/server"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/info"
)

func main() {
	flagSet := flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	var grpcAPIPort = flagSet.StringP(
		"grpc-api-port",
		"g",
		"9999", //TODO: make it configurable using configmap
		"tcp port where sunsetting plugin grpc server is serving")

	var restAPIPort = flagSet.StringP(
		"rest-api-port",
		"r",
		"80", //TODO: make it configurable using configmap
		"tcp port where sunsetting plugin http server is serving")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("unable to parse flags %v\n", err)
	}

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

	mainLogger.Infof("grpc-api-port: %v, rest-api-port: %v", *grpcAPIPort, *restAPIPort)

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

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	srv := &http.Server{Addr: ":" + *restAPIPort, Handler: mux}

	go func() {
		mainLogger.Info("starting api server")
		mainLogger.Fatal(srv.ListenAndServe())
	}()

	listener, err := net.Listen("tcp", ":"+*grpcAPIPort)
	if err != nil {
		mainLogger.Fatal("can't start tcp listener for gRPC server")
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
	// TODO: think how to make this defer correctly
	defer func() {
		_, err := pluginService.Stop(context.Background(), nil)
		if err != nil {
			mainLogger.Errorf("got error at plugin service stop: %+v", err)
		}
	}()
	defer grpcServer.GracefulStop()

	mainLogger.Info("starting gRPC server")
	mainLogger.Fatal(grpcServer.Serve(listener))
}
