package kube

import (
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/pkg/errors"
	corev1api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientSet *kubernetes.Clientset
	logger    logrus.FieldLogger
}

func NewKubeClient(logger logrus.FieldLogger) (*Client, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientSet: clientSet,
		logger:    logger,
	}, nil
}

func (c *Client) GetDaemonset(labelsSet labels.Set) (v1beta1.DaemonSet, error) {
	var labelsSelector = labels.SelectorFromSet(labelsSet)
	var options = metav1.ListOptions{
		LabelSelector: labelsSelector.String(),
	}

	dss, err := c.clientSet.ExtensionsV1beta1().DaemonSets("").List(options)
	if err != nil {
		return v1beta1.DaemonSet{}, errors.Wrap(err, "can't list daemonSets")
	}

	if len(dss.Items) != 1 {
		return v1beta1.DaemonSet{}, errors.Wrap(err, "can't multiple instances of nodeagent are deployed")
	}

	return dss.Items[0], nil
}

// GetDaemonsetPods returns map kubernetes node ip to daemnoset pod on that node
func (c *Client) GetDaemonsetPods(daemonSet v1beta1.DaemonSet) (map[string]corev1api.Pod, error) {
	var result = map[string]corev1api.Pod{}
	fieldSelector, err := fields.ParseSelector(
		"status.phase!=" +
			string(corev1api.PodSucceeded) +
			",status.phase!=" +
			string(corev1api.PodFailed),
	)
	if err != nil {
		return nil, err
	}

	nonTerminatedPodsList, err := c.clientSet.CoreV1().Pods("").List(
		metav1.ListOptions{FieldSelector: fieldSelector.String()},
	)
	if err != nil {
		return nil, err
	}

	if len(nonTerminatedPodsList.Items) == 0 {
		return nil, errors.New("there are no running pods at cluster")
	}

	for _, pod := range nonTerminatedPodsList.Items {
		if len(pod.OwnerReferences) == 1 && pod.OwnerReferences[0].Name == daemonSet.Name {
			result[pod.Status.HostIP] = pod
		}
	}

	if len(result) == 0 {
		return nil, errors.Errorf("there are no running pods for daemonSet %v", daemonSet.Name)
	}

	return result, nil
}

func (c *Client) GetNodeResourceRequirements() (map[string]*NodeResourceRequirements, error) {
	var instanceEntries = map[string]*NodeResourceRequirements{}

	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		fieldSelector, err := fields.ParseSelector(
			"spec.nodeName=" +
				node.Name +
				",status.phase!=" +
				string(corev1api.PodSucceeded) +
				",status.phase!=" +
				string(corev1api.PodFailed),
		)
		if err != nil {
			return nil, err
		}

		nonTerminatedPodsForNode, err := c.clientSet.CoreV1().Pods("").List(
			metav1.ListOptions{FieldSelector: fieldSelector.String()},
		)
		if err != nil {
			return nil, err
		}

		nodeResourceRequirements, err := getNodeResourceRequirements(node, nonTerminatedPodsForNode.Items)
		if err != nil {
			return nil, err
		}

		instanceEntries[nodeResourceRequirements.InstanceID] = nodeResourceRequirements
	}

	return instanceEntries, nil
}

func (c *Client) GetRandomMaster() (*corev1api.Node, error) {

	nodes, err := c.clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var randomMasterNode *corev1api.Node

	for nodeIndex, node := range nodes.Items {
		if node.Status.Phase != corev1api.NodeRunning {
			continue
		}

		for label, value := range node.GetLabels() {
			if label == "node-role.kubernetes.io/master" || (label == "kubernetes.io/role" && value == "master") {
				randomMasterNode = &nodes.Items[nodeIndex]
				break
			}
		}

		if randomMasterNode != nil {
			break
		}
	}

	if randomMasterNode == nil {
		c.logger.Errorf("nodes: %+v", nodes.Items)
		return nil, errors.New("can't identify master node, check that label node-role.kubernetes.io/master is set")
	}

	return randomMasterNode, nil
}

func getNodeResourceRequirements(node corev1api.Node, pods []corev1api.Pod) (*NodeResourceRequirements, error) {
	var nodeResourceRequirements = &NodeResourceRequirements{
		Name:                     node.Name,
		PodsResourceRequirements: []*PodResourceRequirements{},
	}
	var err error

	nodeResourceRequirements.Region, nodeResourceRequirements.InstanceID, err = parseProviderID(node.Spec.ProviderID)
	if err != nil {
		return nil, err
	}

	// calculate worker node requests/limits
	nodeResourceRequirements.PodsResourceRequirements = getPodsRequestsAndLimits(pods)

	var allocatable = node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}

	nodeResourceRequirements.AllocatableCPU = allocatable.Cpu().MilliValue()
	nodeResourceRequirements.AllocatableMemory = allocatable.Memory().Value()

	nodeResourceRequirements.RefreshTotals()

	var internalIP string
	for _, address := range node.Status.Addresses {
		if address.Type == corev1api.NodeInternalIP {
			internalIP = address.Address
		}
	}

	nodeResourceRequirements.internalIPAddress = internalIP

	return nodeResourceRequirements, nil
}

// TODO: add checks and errors
// for aws ProviderID has format - aws:///us-west-1b/i-0c912bfd4048b97e5
// TODO: implement other possible formats of ProviderID
// kubernetesInstanceID represents the id for an instance in the kubernetes API;
// the following form
//  * aws:///<zone>/<awsInstanceId>
//  * aws:////<awsInstanceId>
//  * <awsInstanceId>
func parseProviderID(providerID string) (string, string, error) {
	var s = strings.TrimPrefix(providerID, "aws:///")
	ss := strings.Split(s, "/")
	if len(ss) != 2 {
		return "", "", errors.Errorf("Cant parse ProviderID: %s", providerID)
	}
	return ss[0], ss[1], nil
}

func getPodsRequestsAndLimits(podList []corev1api.Pod) []*PodResourceRequirements {
	var result = []*PodResourceRequirements{}
	for _, pod := range podList {
		p := pod
		var podRR = &PodResourceRequirements{
			PodName: p.Name,
		}

		podReqs, podLimits := PodRequestsAndLimits(&p)
		cpuReqs, cpuLimits := podReqs[corev1api.ResourceCPU], podLimits[corev1api.ResourceCPU]
		memoryReqs, memoryLimits := podReqs[corev1api.ResourceMemory], podLimits[corev1api.ResourceMemory]
		podRR.CPUReqs, podRR.CPULimits = cpuReqs.MilliValue(), cpuLimits.MilliValue()
		podRR.MemoryReqs, podRR.MemoryLimits = memoryReqs.Value(), memoryLimits.Value()

		result = append(result, podRR)
	}

	return result
}

// PodRequestsAndLimits returns a dictionary of all defined resources summed up for all
// containers of the pod.
func PodRequestsAndLimits(pod *corev1api.Pod) (reqs corev1api.ResourceList, limits corev1api.ResourceList) {
	reqs, limits = corev1api.ResourceList{}, corev1api.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}
	return
}

// addResourceList adds the resources in newList to list
func addResourceList(list, new corev1api.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = *quantity.Copy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

// maxResourceList sets list to the greater of list/newList for every resource
// either list
func maxResourceList(list, new corev1api.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = *quantity.Copy()
			continue
		} else if quantity.Cmp(value) > 0 {
			list[name] = *quantity.Copy()
		}
	}
}
