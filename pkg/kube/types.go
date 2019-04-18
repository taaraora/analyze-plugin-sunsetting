package kube

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/labels"
)

func GetNodeAgentLabelsSet() labels.Set {
	return labels.Set{
		"app.kubernetes.io/part-of":   "analyze",
		"app.kubernetes.io/name":      "analyze-nodeagent",
		"app.kubernetes.io/component": "nodeagent",
	}
}

type NodeResourceRequirements struct {
	Name                  string
	IsRecommendedToSunset bool
	// extracted from ProviderID
	Region     string
	InstanceID string

	PodsResourceRequirements []*PodResourceRequirements
	AllocatableCPU           int64
	AllocatableMemory        int64

	//nolint
	cpuReqs int64
	//nolint
	cpuLimits    int64
	memoryReqs   int64
	memoryLimits int64

	fractionCPUReqs      float64
	fractionCPULimits    float64
	fractionMemoryReqs   float64
	fractionMemoryLimits float64

	internalIPAddress string
}

type PodResourceRequirements struct {
	PodName      string `json:"podName"`
	CPUReqs      int64  `json:"cpuRequests"`
	CPULimits    int64  `json:"cpuLimits"`
	MemoryReqs   int64  `json:"memoryRequests"`
	MemoryLimits int64  `json:"memoryLimits"`
}

// RefreshTotals recalculates total node requests and limits and their fractional representation
// need to be invoked every time when PodsResourceRequirements or AllocatableMemory or AllocatableCPU where changed
func (n *NodeResourceRequirements) RefreshTotals() {
	n.cpuReqs = 0
	n.cpuLimits = 0
	n.memoryReqs = 0
	n.memoryLimits = 0
	n.fractionCPUReqs = 0
	n.fractionCPULimits = 0
	n.fractionMemoryReqs = 0
	n.fractionMemoryLimits = 0

	for _, podRR := range n.PodsResourceRequirements {
		n.cpuReqs += podRR.CPUReqs
		n.cpuLimits += podRR.CPULimits
		n.memoryReqs += podRR.MemoryReqs
		n.memoryLimits += podRR.MemoryLimits
	}

	if n.AllocatableCPU != 0 {
		n.fractionCPUReqs = float64(n.cpuReqs) / float64(n.AllocatableCPU) * 100
		n.fractionCPULimits = float64(n.cpuLimits) / float64(n.AllocatableCPU) * 100
	}

	if n.AllocatableMemory != 0 {
		n.fractionMemoryReqs = float64(n.memoryReqs) / float64(n.AllocatableMemory) * 100
		n.fractionMemoryLimits = float64(n.memoryLimits) / float64(n.AllocatableMemory) * 100
	}
}

func (n *NodeResourceRequirements) CPUReqs() int64 {
	return n.cpuReqs
}

func (n *NodeResourceRequirements) CPULimits() int64 {
	return n.cpuLimits
}

func (n *NodeResourceRequirements) MemoryReqs() int64 {
	return n.memoryReqs
}

func (n *NodeResourceRequirements) MemoryLimits() int64 {
	return n.memoryLimits
}

func (n *NodeResourceRequirements) FractionCPUReqs() float64 {
	return n.fractionCPUReqs
}

func (n *NodeResourceRequirements) FractionCPULimits() float64 {
	return n.fractionCPULimits
}

func (n *NodeResourceRequirements) FractionMemoryReqs() float64 {
	return n.fractionMemoryReqs
}

func (n *NodeResourceRequirements) FractionMemoryLimits() float64 {
	return n.fractionMemoryLimits
}

func (n *NodeResourceRequirements) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":                     n.Name,
		"isRecommendedToSunset":    n.IsRecommendedToSunset,
		"region":                   n.Region,
		"instanceId":               n.InstanceID,
		"podsResourceRequirements": n.PodsResourceRequirements,
		"allocatableCpu":           n.AllocatableCPU,
		"allocatableMemory":        n.AllocatableMemory,
		"cpuRequests":              n.cpuReqs,
		"cpuLimits":                n.cpuLimits,
		"memoryRequests":           n.memoryReqs,
		"memoryLimits":             n.memoryLimits,
		"fractionCpuRequests":      n.fractionCPUReqs,
		"fractionCpuLimits":        n.fractionCPULimits,
		"fractionMemoryRequests":   n.fractionMemoryReqs,
		"fractionMemoryLimits":     n.fractionMemoryLimits,
	})
}

func (n *NodeResourceRequirements) IPAddress() string {
	return n.internalIPAddress
}
