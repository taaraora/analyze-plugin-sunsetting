package main

import (
	"encoding/json"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/supergiant/analyze-plugin-sunsetting/cmd/analyze-sunsetting/server"
	"github.com/supergiant/analyze/pkg/plugin/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

//TODO: this was copy pasted form analyze repository, need to think how it is better to extract it as dependency
type Plugin struct {
	// detailed plugin description
	Description string `json:"description,omitempty"`

	// unique ID of installed plugin
	// basically it is slugged URI of plugin repository name e. g. supergiant-request-limits-check
	//
	ID string `json:"id,omitempty"`

	// date/Time the plugin was installed
	// Format: date-time
	InstalledAt time.Time `json:"installedAt,omitempty"`

	// name is the name of the plugin.
	Name string `json:"name,omitempty"`

	// service labels
	ServiceLabels map[string]string `json:"serviceLabels,omitempty"`

	// name of k8s service which is front of plugin deployment
	ServiceName string `json:"serviceName,omitempty"`

	// plugin status
	Status string `json:"status,omitempty"`

	// plugin version, major version shall be equal to robots version
	Version string `json:"version,omitempty"`
}

func main()  {
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

	grpcApiPort, err := cmd.Flags().GetString("grpc-api-port")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag grpc-api-port")
	}

	restApiPort, err := cmd.Flags().GetString("rest-api-port")
	if err != nil {
		return errors.Wrap(err, "unable to get config flag rest-api-port")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel) //TODO: make it configurable
	var mainLogger = logger.WithField("app", "sunsetting-plugin")

	mainLogger.Infof("grpc-api-port: %v, rest-api-port: %v", grpcApiPort, restApiPort)

	//TODO: for now it returns only info about registered plugin, but also need to serve bundles
	handler := func (w http.ResponseWriter, r *http.Request) {
		p := &Plugin{
			Description:   "detailed plugin description",
			ID:            "unique ID of installed plugin",
			InstalledAt:   time.Now(),
			Name:          "name is the name of the plugin",
			ServiceLabels: nil,
			ServiceName:   "",
			Status:        "OK",
			Version:       "v2.0.1",
		}
		b, _ := json.Marshal(p)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}

	http.HandleFunc("/api/v1/info", handler)
	go func() {
		mainLogger.Fatal(http.ListenAndServe(":" + restApiPort, nil))
	}()


	listener, err := net.Listen("tcp", ":"+grpcApiPort)
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
	var pluginService = server.NewServer(mainLogger.WithField("component", "plugin_server"))
	proto.RegisterPluginServer(grpcServer, pluginService)


	return grpcServer.Serve(listener)
}
