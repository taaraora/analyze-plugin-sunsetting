package server

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"go.etcd.io/etcd/clientv3"

	"github.com/supergiant/analyze-plugin-sunsetting/pkg/cloudprovider"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/cloudprovider/aws"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/info"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/kube"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/nodeagent"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/reshuffle"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/storage"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/storage/etcd"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/sunsetting"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/supergiant/analyze/pkg/plugin/proto"
)

type server struct {
	config                 *proto.PluginConfig
	nodeAgentClient        *nodeagent.Client
	awsClient              *aws.Client
	kubeClient             *kube.Client
	computeInstancesPrices map[string][]cloudprovider.ProductPrice
	logger                 logrus.FieldLogger
	etcdStorage            storage.Interface
}

type msg []byte

func (m msg) Payload() []byte {
	return m
}

func checkResult() *proto.CheckResult {
	return &proto.CheckResult{
		ExecutionStatus: "OK",
		Status:          proto.CheckStatus_UNKNOWN_CHECK_STATUS,
		Name:            "Underutilized nodes sunsetting Check",
		Description: &any.Any{
			TypeUrl: "io.supergiant.analyze.plugin.requestslimitscheck",
			Value: []byte(
				"Resources (CPU/RAM) total capacity and allocatable where checked on nodes of k8s cluster. Results:",
			),
		},
	}
}

func NewServer(logger logrus.FieldLogger) proto.PluginServer {
	return &server{
		logger: logger,
	}
}

func (u *server) Check(ctx context.Context, in *proto.CheckRequest) (*proto.CheckResponse, error) {
	u.logger.Info("got Check request")
	defer u.logger.Info("Check request done")
	var nodeResourceRequirements, err = u.kubeClient.GetNodeResourceRequirements()
	if err != nil {
		u.logger.Errorf("unable to get nodeResourceRequirements, %v", err)
		return nil, errors.Wrap(err, "unable to get nodeResourceRequirements")
	}

	nodeAgentsDaemonSet, err := u.kubeClient.GetDaemonset(kube.GetNodeAgentLabelsSet())
	if err != nil {
		u.logger.Errorf("unable to get nodeAgentsDaemonSet, %v", err)
		return nil, errors.Wrap(err, "unable to get nodeAgentsDaemonSet")
	}

	nodeagentPods, err := u.kubeClient.GetDaemonsetPods(nodeAgentsDaemonSet)
	if err != nil {
		u.logger.Errorf("unable to get GetDaemonsetPods, %v", err)
		return nil, errors.Wrap(err, "unable to get GetDaemonsetPods")
	}

	var computeInstances = make(map[string]cloudprovider.ComputeInstance)
	for instanceID, resourceRequirements := range nodeResourceRequirements {
		nodeagentPod, exists := nodeagentPods[resourceRequirements.IPAddress()]
		if !exists {
			u.logger.Errorf(
				"There is no analyze nodeAgent is running for nodeIP, %s",
				resourceRequirements.IPAddress(),
			)
			continue
		}
		var nodeAgentInstance = nodeagent.Instance{
			HostIP: nodeagentPod.Status.HostIP,
			PodIP:  nodeagentPod.Status.PodIP,
		}
		fetchedInstanceID, err := u.nodeAgentClient.Get(nodeAgentInstance.PodURI() + "/aws/meta-data/instance-id")
		if err != nil {
			u.logger.Errorf("cant fetch ID for node %s, error %v", instanceID, err)
			continue
		}
		if fetchedInstanceID != instanceID {
			u.logger.Errorf(
				"fetched ec2 instanceID: %s not equal to instanceID from providerID %s",
				fetchedInstanceID,
				instanceID,
			)
			continue
		}

		instanceType, err := u.nodeAgentClient.Get(nodeAgentInstance.PodURI() + "/aws/meta-data/instance-type")
		if err != nil {
			u.logger.Errorf("can't ec2 instance Type: %+v", err)
			continue
		}
		computeInstances[instanceID] = cloudprovider.ComputeInstance{
			InstanceID:   instanceID,
			InstanceType: instanceType,
		}
	}

	var unsortedEntries = make([]*sunsetting.InstanceEntry, 0)
	var result = make([]sunsetting.InstanceEntry, 0)

	// create InstanceEntries by combining nodeResourceRequirements with ec2 instance type and price
	for InstanceID, computeInstance := range computeInstances {
		var kubeNode, exists = nodeResourceRequirements[InstanceID]
		if !exists {
			continue
		}

		// TODO: fix me when prices collecting will be clear
		// TODO: We need to match it with instance tenancy?
		var instanceTypePrice cloudprovider.ProductPrice
		var instanceTypePrices, exist = u.computeInstancesPrices[computeInstance.InstanceType]
		if exist {
			for _, priceItem := range instanceTypePrices {
				if strings.Contains(priceItem.UsageType, "BoxUsage") && priceItem.ValuePerUnit != "0.0000000000" {
					instanceTypePrice = priceItem
				}
			}
			if instanceTypePrice.InstanceType == "" && len(instanceTypePrices) > 0 {
				instanceTypePrice = instanceTypePrices[0]
			}
		}

		result = append(result, sunsetting.InstanceEntry{
			CloudProvider: computeInstance,
			Price:         instanceTypePrice,
			WorkerNode:    *kubeNode,
		})
		unsortedEntries = append(unsortedEntries, &sunsetting.InstanceEntry{
			CloudProvider: computeInstance,
			Price:         instanceTypePrice,
			WorkerNode:    *kubeNode,
		})
	}

	//TODO: double check logic, is it really needed?
	var instancesToSunset = sunsetting.CheckEachPodOneByOne(unsortedEntries)
	if len(instancesToSunset) == 0 {
		instancesToSunset = sunsetting.CheckAllPodsAtATime(unsortedEntries)
	}

	workerNodesToSunset := make([]string, 0)

	// mark nodes selected node with IsRecommendedToSunset == true
	for i := range result {
		for _, entryToSunset := range instancesToSunset {
			if entryToSunset.CloudProvider.InstanceID == result[i].CloudProvider.InstanceID {
				result[i].WorkerNode.IsRecommendedToSunset = true
				workerNodesToSunset = append(workerNodesToSunset, entryToSunset.CloudProvider.InstanceID)
			}
		}
	}

	b, _ := json.Marshal(result)

	checkResult := checkResult()
	checkResult.Description = &any.Any{
		TypeUrl: "io.supergiant.analyze.plugin.sunsetting",
		Value:   b,
	}

	if len(instancesToSunset) == 0 {
		checkResult.Status = proto.CheckStatus_GREEN
	} else {
		checkResult.Status = proto.CheckStatus_YELLOW
	}

	master, err := u.kubeClient.GetRandomMaster()
	if err != nil {
		u.logger.Errorf("can't find master: %v", err)
		return &proto.CheckResponse{Result: checkResult}, nil
	}

	if master != nil {
		rc := reshuffle.ReshufflePodsCommand{
			ClusterID:      master.Status.NodeInfo.MachineID,
			WorkerNodesIDs: workerNodesToSunset,
			AZ:             u.awsClient.Region(),
		}

		b, err := json.Marshal(rc)
		if err != nil {
			u.logger.Errorf("can't marshal command: %+v", err)
		}

		// TODO: refactor this, to make it testable
		re := reshuffle.CommandEnvelope{
			ID:       uuid.New(),
			Type:     reshuffle.ReshufflePodsCommandType,
			SourceID: master.Name,
			Payload:  b,
		}

		r, err := json.Marshal(re)
		if err != nil {
			u.logger.Errorf("can't marshal envelope: %+v", err)
		}

		err = u.etcdStorage.Put(context.Background(), storage.ReshuffleCommandPrefix, master.Name, msg(r))
		if err != nil {
			u.logger.Errorf("can't write envelope: %+v", err)
		}
	}

	return &proto.CheckResponse{Result: checkResult}, nil
}

func (u *server) Configure(ctx context.Context, pluginConfig *proto.PluginConfig) (*empty.Empty, error) {
	u.logger.Info("got Configure request")
	defer u.logger.Info("Configure request done")
	u.config = pluginConfig

	nodeAgentClient, err := nodeagent.NewClient(logrus.New())
	if err != nil {
		return nil, err
	}
	u.nodeAgentClient = nodeAgentClient

	var awsClientConfig = pluginConfig.GetAwsConfig()
	awsClient, err := aws.NewClient(awsClientConfig, u.logger.WithField("component", "awsClient"))
	if err != nil {
		return nil, err
	}
	u.awsClient = awsClient

	//TODO: may be we need just log warning?
	u.computeInstancesPrices, err = u.awsClient.GetPrices()
	if err != nil {
		return nil, err
	}

	u.kubeClient, err = kube.NewKubeClient(u.logger.WithField("component", "kubeClient"))
	if err != nil {
		return nil, err
	}

	u.logger.Infof("got configuration %+v", *pluginConfig)

	// TODO: make it correctly stopped when configure is invoked multiple times
	u.etcdStorage, err = etcd.NewETCDStorage(
		clientv3.Config{Endpoints: u.config.EtcdEndpoints},
		u.logger.WithField("component", "etcdClient"),
	)
	if err != nil {
		u.logger.Errorf("unable to create ETCD client, err: %v", err)
	}

	return &empty.Empty{}, nil
}

func (u *server) Info(ctx context.Context, in *empty.Empty) (*proto.PluginInfo, error) {
	u.logger.Info("got info request")
	defer u.logger.Info("info request done")

	return &proto.PluginInfo{
		Id:          info.Info().ID,
		Version:     info.Info().Version,
		Name:        info.Info().Name,
		Description: info.Info().Description,
		DefaultConfig: &proto.PluginConfig{
			ExecutionInterval: ptypes.DurationProto(2 * time.Minute),
		},
	}, nil
}

func (u *server) Stop(ctx context.Context, in *proto.Stop_Request) (*proto.Stop_Response, error) {
	err := u.etcdStorage.Close()
	return &proto.Stop_Response{}, err
}
