package sunsetting

import (
	"testing"

	"github.com/supergiant/analyze-plugin-sunsetting/cloudprovider"
	"github.com/supergiant/analyze-plugin-sunsetting/kube"
)

func unsortedEntries() []*InstanceEntry {
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

func TestCheckAllPodsAtATime_CandidateForSunsettingFound(t *testing.T) {
	in := unsortedEntries()
	var result = CheckAllPodsAtATime(in)
	if !(len(result) != 0 && result[0].CloudProvider.InstanceID == "i-03fb8e89232700cc3") {
		t.Fatal(result)
	}
}

func TestCheckEachPodOneByOne_CandidateForSunsettingFound(t *testing.T) {
	in := unsortedEntries()
	var result = CheckEachPodOneByOne(in)
	if !(len(result) != 0 && result[0].CloudProvider.InstanceID == "i-03fb8e89232700cc3") {
		t.Fatal(result)
	}
}
