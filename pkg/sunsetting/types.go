package sunsetting

import (
	"fmt"
	"sort"
	"strings"

	"github.com/supergiant/analyze-plugin-sunsetting/pkg/cloudprovider"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/kube"
)

// InstanceEntry struct represents Kelly's "instances to sunset" table entry,
// It consists of some k8s cluster worker node params ans some ec2 instance params
type InstanceEntry struct {
	CloudProvider cloudprovider.ComputeInstance `json:"cloudProvider"`
	Price         cloudprovider.ProductPrice    `json:"price"`
	WorkerNode    kube.NodeResourceRequirements `json:"kube"`
}

func (m *InstanceEntry) RAMWasted() int64 {
	m.WorkerNode.RefreshTotals()
	return m.WorkerNode.AllocatableMemory - m.WorkerNode.MemoryReqs()
}

func (m *InstanceEntry) RAMRequested() int64 {
	m.WorkerNode.RefreshTotals()
	return m.WorkerNode.MemoryReqs()
}

func (m *InstanceEntry) CPUWasted() int64 {
	m.WorkerNode.RefreshTotals()
	return m.WorkerNode.AllocatableCPU - m.WorkerNode.CPUReqs()
}

// EntriesByWastedRAM implements sort.Interface based on the value returned by NodeResourceRequirements.RAMWasted().
type EntriesByWastedRAM []*InstanceEntry

// NewSortedEntriesByWastedRAM makes deep copy of passed slice and invoke sort.Sort on it
// NewSortedEntriesByWastedRAM returns copied slice sorted in descending order
func NewSortedEntriesByWastedRAM(in []*InstanceEntry) EntriesByWastedRAM {
	var res = make([]*InstanceEntry, len(in))
	copy(res, in)
	var entries = EntriesByWastedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByWastedRAM) Len() int           { return len(e) }
func (e EntriesByWastedRAM) Less(i, j int) bool { return e[i].RAMWasted() < e[j].RAMWasted() }
func (e EntriesByWastedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (e EntriesByWastedRAM) String() string {
	res := []string{"max mem wasted: "}
	for i := range e {
		if e[i] != nil {
			e[i].WorkerNode.RefreshTotals()
			bytes := float64(e[i].RAMWasted())
			s := fmt.Sprintf(
				"id:%s, %p, %vGB",
				e[i].CloudProvider.InstanceID,
				e[i],
				bytes/1000000000,
			)
			res = append(res, s)
		}
	}
	return strings.Join(res, ", ")
}

// EntriesByRequestedRAM implements sort.Interface based on the value
// returned by NodeResourceRequirements.RAMRequested().
type EntriesByRequestedRAM []*InstanceEntry

// NewSortedEntriesByRequestedRAM makes deep copy of passed slice and invoke sort.Sort on it
// NewSortedEntriesByRequestedRAM returns copied slice sorted in descending order
func NewSortedEntriesByRequestedRAM(in []*InstanceEntry) EntriesByRequestedRAM {
	var res = make([]*InstanceEntry, len(in))
	copy(res, in)
	var entries = EntriesByRequestedRAM(res)
	sort.Sort(sort.Reverse(entries))

	return entries
}

func (e EntriesByRequestedRAM) Len() int           { return len(e) }
func (e EntriesByRequestedRAM) Less(i, j int) bool { return e[i].RAMRequested() < e[j].RAMRequested() }
func (e EntriesByRequestedRAM) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (e EntriesByRequestedRAM) String() string {
	res := []string{"max mem requested: "}
	for i := range e {
		if e[i] != nil {
			e[i].WorkerNode.RefreshTotals()
			bytes := float64(e[i].WorkerNode.MemoryReqs())
			s := fmt.Sprintf(
				"id:%s, %p, %vGB",
				e[i].CloudProvider.InstanceID,
				e[i],
				bytes/1000000000,
			)
			res = append(res, s)
		}
	}
	return strings.Join(res, ", ")
}
