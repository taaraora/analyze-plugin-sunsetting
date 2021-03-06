package sunsetting

import (
	"math"
	"testing"

	"github.com/supergiant/analyze-plugin-sunsetting/pkg/cloudprovider"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/kube"
)

func fixture() *InstanceEntry {
	return &InstanceEntry{
		CloudProvider: cloudprovider.ComputeInstance{
			InstanceID:   "i-028bc20adaf2311d6",
			InstanceType: "m4.large",
		},
		Price: cloudprovider.ProductPrice{
			InstanceType: "m4.large",
			Memory:       "8 GiB",
			Vcpu:         "2",
			Unit:         "Hrs",
			Currency:     "USD",
			ValuePerUnit: "0.0000000000",
			UsageType:    "USW1-HostBoxUsage:m4.large",
			Tenancy:      "Host",
		},
		WorkerNode: kube.NodeResourceRequirements{
			Name:              "ip-172-20-1-21.us-west-1.compute.internal",
			Region:            "us-west-1b",
			InstanceID:        "i-028bc20adaf2311d6",
			AllocatableCPU:    2000,
			AllocatableMemory: 8260218880,
			PodsResourceRequirements: []*kube.PodResourceRequirements{
				&kube.PodResourceRequirements{
					PodName:      "dmts-es-1-elasticsearch-client-848f4d5db6-bvksx",
					CPUReqs:      25,
					CPULimits:    1000,
					MemoryReqs:   536870912,
					MemoryLimits: 0,
				},
				&kube.PodResourceRequirements{
					PodName:      "dmts-es-1-elasticsearch-master-0",
					CPUReqs:      25,
					CPULimits:    1000,
					MemoryReqs:   536870912,
					MemoryLimits: 0,
				},
				&kube.PodResourceRequirements{
					PodName:      "dmts-es-2-elasticsearch-client-56db9bfb6-h79wr",
					CPUReqs:      25,
					CPULimits:    1000,
					MemoryReqs:   536870912,
					MemoryLimits: 0,
				},
				&kube.PodResourceRequirements{
					PodName:      "dmts-es-2-elasticsearch-data-0",
					CPUReqs:      25,
					CPULimits:    1000,
					MemoryReqs:   1610612736,
					MemoryLimits: 0,
				},
				&kube.PodResourceRequirements{
					PodName:      "dmts-es-2-elasticsearch-master-0",
					CPUReqs:      25,
					CPULimits:    1000,
					MemoryReqs:   536870912,
					MemoryLimits: 0,
				},
				&kube.PodResourceRequirements{
					PodName:      "errant-sheep-analyze-764c68787c-4c2gs",
					CPUReqs:      300,
					CPULimits:    0,
					MemoryReqs:   268435456,
					MemoryLimits: 0,
				},
				&kube.PodResourceRequirements{
					PodName:      "kube-proxy-ip-172-20-1-21.us-west-1.compute.internal",
					CPUReqs:      0,
					CPULimits:    0,
					MemoryReqs:   0,
					MemoryLimits: 0,
				},
			},
		},
	}
}
func fixtures() []*InstanceEntry {
	return []*InstanceEntry{
		&InstanceEntry{
			CloudProvider: cloudprovider.ComputeInstance{
				InstanceID:   "i-03fb8e89232700cc3",
				InstanceType: "m4.large",
			},
			Price: cloudprovider.ProductPrice{
				InstanceType: "m4.large",
				Memory:       "8 GiB",
				Vcpu:         "2",
				Unit:         "Hrs",
				Currency:     "USD",
				ValuePerUnit: "0.0000000000",
				UsageType:    "USW1-HostBoxUsage:m4.large",
				Tenancy:      "Host",
			},
			WorkerNode: kube.NodeResourceRequirements{
				Name:              "ip-172-20-1-44.us-west-1.compute.internal",
				Region:            "us-west-1b",
				InstanceID:        "i-03fb8e89232700cc3",
				AllocatableCPU:    2000,
				AllocatableMemory: 8260214784,
				PodsResourceRequirements: []*kube.PodResourceRequirements{
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-data-1",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   1610612736,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-master-2",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-client-56db9bfb6-wc462",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-master-2",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "coredns-77cd44d8df-bsl5j",
						CPUReqs:      100,
						CPULimits:    0,
						MemoryReqs:   73400320,
						MemoryLimits: 178257920,
					},
					&kube.PodResourceRequirements{
						PodName:      "heapster-v11-2j9kb",
						CPUReqs:      100,
						CPULimits:    100,
						MemoryReqs:   222298112,
						MemoryLimits: 222298112,
					},
					&kube.PodResourceRequirements{
						PodName:      "monitoring-influxdb-grafana-v2-vr9jh",
						CPUReqs:      200,
						CPULimits:    200,
						MemoryReqs:   314572800,
						MemoryLimits: 314572800,
					},
					&kube.PodResourceRequirements{
						PodName:      "tiller-deploy-677f9cb999-jv6ct",
						CPUReqs:      0,
						CPULimits:    0,
						MemoryReqs:   0,
						MemoryLimits: 0,
					},
				},
			},
		},
		&InstanceEntry{
			CloudProvider: cloudprovider.ComputeInstance{
				InstanceID:   "i-028bc20adaf2311d6",
				InstanceType: "m4.large",
			},
			Price: cloudprovider.ProductPrice{
				InstanceType: "m4.large",
				Memory:       "8 GiB",
				Vcpu:         "2",
				Unit:         "Hrs",
				Currency:     "USD",
				ValuePerUnit: "0.0000000000",
				UsageType:    "USW1-HostBoxUsage:m4.large",
				Tenancy:      "Host",
			},
			WorkerNode: kube.NodeResourceRequirements{
				Name:              "ip-172-20-1-21.us-west-1.compute.internal",
				Region:            "us-west-1b",
				InstanceID:        "i-028bc20adaf2311d6",
				AllocatableCPU:    2000,
				AllocatableMemory: 8260218880,
				PodsResourceRequirements: []*kube.PodResourceRequirements{
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-client-848f4d5db6-bvksx",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-master-0",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-client-56db9bfb6-h79wr",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-data-0",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   1610612736,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-master-0",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "errant-sheep-analyze-764c68787c-4c2gs",
						CPUReqs:      300,
						CPULimits:    0,
						MemoryReqs:   268435456,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "kube-proxy-ip-172-20-1-21.us-west-1.compute.internal",
						CPUReqs:      0,
						CPULimits:    0,
						MemoryReqs:   0,
						MemoryLimits: 0,
					},
				},
			},
		},
		&InstanceEntry{
			CloudProvider: cloudprovider.ComputeInstance{
				InstanceID:   "i-0898d927727329231",
				InstanceType: "m4.large",
			},
			Price: cloudprovider.ProductPrice{
				InstanceType: "m4.large",
				Memory:       "8 GiB",
				Vcpu:         "2",
				Unit:         "Hrs",
				Currency:     "USD",
				ValuePerUnit: "0.0000000000",
				UsageType:    "USW1-HostBoxUsage:m4.large",
				Tenancy:      "Host",
			},
			WorkerNode: kube.NodeResourceRequirements{
				Name:              "ip-172-20-1-242.us-west-1.compute.internal",
				Region:            "us-west-1b",
				InstanceID:        "i-0898d927727329231",
				AllocatableCPU:    2000,
				AllocatableMemory: 8260218880,
				PodsResourceRequirements: []*kube.PodResourceRequirements{
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-client-848f4d5db6-xll2f",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-data-0",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   1610612736,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-1-elasticsearch-master-1",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-data-0",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   1610612736,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-data-1",
						CPUReqs:      25,
						CPULimits:    1000,
						MemoryReqs:   1610612736,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "dmts-es-2-elasticsearch-master-1",
						CPUReqs:      300,
						CPULimits:    0,
						MemoryReqs:   536870912,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "etcd-d74c5fbff-85sws",
						CPUReqs:      250,
						CPULimits:    0,
						MemoryReqs:   256000000,
						MemoryLimits: 0,
					},
					&kube.PodResourceRequirements{
						PodName:      "kube-proxy-ip-172-20-1-242.us-west-1.compute.internal",
						CPUReqs:      0,
						CPULimits:    0,
						MemoryReqs:   0,
						MemoryLimits: 0,
					},
				},
			},
		},
	}
}

func TestWorkerNode_RefreshTotals_CorrectCounts(t *testing.T) {
	var in = fixture()

	if in.RAMRequested() != 536870912+536870912+536870912+1610612736+536870912+268435456+0 {
		t.Fatal("RAMRequested returned incorrect value", in)
	}

	if in.RAMWasted() != 8260218880-(536870912+536870912+536870912+1610612736+536870912+268435456+0) {
		t.Fatal("RAMWasted returned incorrect value", in)
	}
}

func TestNewSortedEntriesByWastedRAM_singleElement(t *testing.T) {
	var in = []*InstanceEntry{fixture()}

	var sorted = NewSortedEntriesByWastedRAM(in)

	if len(in) != len(sorted) {
		t.Fatal("incorrect number of items")
	}

	var f = in[0]
	var s = sorted[0]

	compareInstanceEntries(t, f, s)
}

func TestNewSortedEntriesByRequestedRAM_elementsAreEqual(t *testing.T) {
	var in = []*InstanceEntry{fixture()}

	var sorted = NewSortedEntriesByRequestedRAM(in)

	if len(in) != len(sorted) {
		t.Fatal("incorrect number of items")
	}

	var f = in[0]
	var s = sorted[0]

	compareInstanceEntries(t, f, s)
}

func TestNewSortedEntriesByWastedRAM_multipleElements(t *testing.T) {
	var in = fixtures()

	var sorted = NewSortedEntriesByWastedRAM(in)

	if len(in) != len(sorted) && len(sorted) == 1 {
		t.Fatal("incorrect number of items")
	}

	for _, sortedItem := range sorted {
		for _, unsortedItem := range in {
			if sortedItem.CloudProvider.InstanceID == unsortedItem.CloudProvider.InstanceID {
				compareInstanceEntries(t, unsortedItem, sortedItem)
			}
		}
	}
}

func TestNewSortedEntriesByRequestedRAM_multipleElements(t *testing.T) {
	var in = fixtures()

	var sorted = NewSortedEntriesByRequestedRAM(in)

	if len(in) != len(sorted) && len(sorted) == 1 {
		t.Fatal("incorrect number of items")
	}

	for _, sortedItem := range sorted {
		for _, unsortedItem := range in {
			if sortedItem.CloudProvider.InstanceID == unsortedItem.CloudProvider.InstanceID {
				compareInstanceEntries(t, unsortedItem, sortedItem)
			}
		}
	}
}

func TestNewSortedEntriesByWastedRAM_multipleInstanceEntriesAreSorted(t *testing.T) {
	var in = fixtures()
	var wastedRAM int64 = math.MaxInt64

	var sorted = NewSortedEntriesByWastedRAM(in)

	if len(in) != len(sorted) && len(sorted) == 1 {
		t.Fatal("incorrect number of items")
	}

	// sorting order is descending
	for _, item := range sorted {
		if item.RAMWasted() > wastedRAM {
			t.Log(sorted)
			t.Fatalf(
				"looks like order of elemnts is not based of wasted ram. current item wasted memory: %v, previos: %v",
				item.RAMWasted(),
				wastedRAM,
			)
		}
		wastedRAM = item.RAMWasted()
	}
}

func TestNewSortedEntriesByRequestedRAM_multipleInstanceEntriesAreSorted(t *testing.T) {
	var in = fixtures()
	var requestedRAM int64 = math.MaxInt64

	var sorted = NewSortedEntriesByRequestedRAM(in)

	if len(in) != len(sorted) && len(sorted) == 1 {
		t.Fatal("incorrect number of items")
	}

	// sorting order is descending
	for _, item := range sorted {
		if item.RAMRequested() > requestedRAM {
			t.Log(sorted)
			t.Fatalf(
				"looks like order of elemnts is not based of wasted ram. current item wasted memory: %v, previos: %v",
				item.RAMWasted(),
				requestedRAM,
			)
		}
		requestedRAM = item.RAMRequested()
	}
}

func compareInstanceEntries(t *testing.T, f *InstanceEntry, s *InstanceEntry) {
	t.Helper()

	if s.Price != f.Price {
		t.Fatal("price structures are not unequal ")
	}

	if s.CloudProvider != f.CloudProvider {
		t.Fatal("cloud provider structures are not unequal ")

	}

	// we do not need to compare all non exportable fields
	if f.WorkerNode.InstanceID != s.WorkerNode.InstanceID ||
		f.WorkerNode.Name != s.WorkerNode.Name ||
		f.WorkerNode.AllocatableCPU != s.WorkerNode.AllocatableCPU ||
		f.WorkerNode.AllocatableMemory != s.WorkerNode.AllocatableMemory ||
		f.WorkerNode.Region != s.WorkerNode.Region {
		t.Logf("node 1: %v, node 2: %v ", f.WorkerNode, s.WorkerNode)
		t.Fatal("WorkerNodes are not equal")
	}

	for i, podRR := range f.WorkerNode.PodsResourceRequirements {
		if *podRR != *s.WorkerNode.PodsResourceRequirements[i] {
			t.Fatal("pods RRs are not equal")
		}
	}
}
