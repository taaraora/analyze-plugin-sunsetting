package sunsetting

import "sort"

// CheckAllPodsAtATime makes simple check that it is possible to move all pods of a node to another node.
func CheckAllPodsAtATime(unsortedEntries []*InstanceEntry) []InstanceEntry {
	var entriesByWastedRAM = NewSortedEntriesByWastedRAM(unsortedEntries)
	var res = make([]InstanceEntry, 0)

	for _, maxWastedRAMEntry := range entriesByWastedRAM {
		maxWastedRAMEntry.WorkerNode.RefreshTotals()

		for i := len(entriesByWastedRAM) - 1; i > 0; i-- {
			var instanceWithMinimumWastedRAM = entriesByWastedRAM[i]
			instanceWithMinimumWastedRAM.WorkerNode.RefreshTotals()
			// check that all requested memory of instance can be moved to another instance
			var wastedRAM = instanceWithMinimumWastedRAM.RAMWasted()
			var wastedCPU = instanceWithMinimumWastedRAM.CPUWasted()
			if maxWastedRAMEntry.WorkerNode.MemoryReqs() <= wastedRAM &&
				maxWastedRAMEntry.WorkerNode.CPUReqs() <= wastedCPU {
				//sunset this instance
				res = append(res, *maxWastedRAMEntry)
				//change memory requests of node which receive all workload
				instanceWithMinimumWastedRAM.WorkerNode.PodsResourceRequirements = append(
					instanceWithMinimumWastedRAM.WorkerNode.PodsResourceRequirements,
					maxWastedRAMEntry.WorkerNode.PodsResourceRequirements...,
				)
				instanceWithMinimumWastedRAM.WorkerNode.RefreshTotals()

				maxWastedRAMEntry.WorkerNode.PodsResourceRequirements = nil
				maxWastedRAMEntry.WorkerNode.RefreshTotals()
				break
			}
		}
	}

	return res
}

// CheckEachPodOneByOne returns instances which are possible to sunset.
func CheckEachPodOneByOne(unsortedEntries []*InstanceEntry) []InstanceEntry {
	var entriesByWastedRAM = NewSortedEntriesByWastedRAM(unsortedEntries)
	var entriesByRequestedRAM = NewSortedEntriesByRequestedRAM(unsortedEntries)
	var res = make([]InstanceEntry, 0)

	for maxWastedRAMEntryIndex := range entriesByWastedRAM {
		// sort pods in descending order by requested memory to more effective pods packing on nodes
		podsRR := entriesByWastedRAM[maxWastedRAMEntryIndex].WorkerNode.PodsResourceRequirements

		sort.Slice(
			podsRR,
			func(i, j int) bool {
				return podsRR[i].MemoryReqs > podsRR[j].MemoryReqs
			},
		)

		// check
		for 0 < len(podsRR) {
			var podRR = podsRR[0]
			for _, maxRequestedRAMEntry := range entriesByRequestedRAM {
				// we need not to move pod between the same node in different slices
				if entriesByWastedRAM[maxWastedRAMEntryIndex].CloudProvider.InstanceID ==
					maxRequestedRAMEntry.CloudProvider.InstanceID {
					continue
				}
				if maxRequestedRAMEntry.RAMWasted() >= podRR.MemoryReqs &&
					maxRequestedRAMEntry.CPUWasted() >= podRR.CPUReqs {
					// we can move the pod
					// delete it from  maxWastedRAMEntry
					podsRR = append(podsRR[:0], podsRR[1:]...)

					entriesByWastedRAM[maxWastedRAMEntryIndex].WorkerNode.RefreshTotals()
					//and add to maxRequestedRAMEntry
					maxRequestedRAMEntry.WorkerNode.PodsResourceRequirements = append(
						maxRequestedRAMEntry.WorkerNode.PodsResourceRequirements,
						podRR,
					)
					maxRequestedRAMEntry.WorkerNode.RefreshTotals()

					goto NextPod
				}
			}
			goto NextNode
		NextPod:
		}
		if len(podsRR) == 0 {
			res = append(res, *entriesByWastedRAM[maxWastedRAMEntryIndex])

			for i, maxRequestedRAMEntry := range entriesByRequestedRAM {
				if entriesByWastedRAM[maxWastedRAMEntryIndex].CloudProvider.InstanceID ==
					maxRequestedRAMEntry.CloudProvider.InstanceID {
					entriesByRequestedRAM = append(entriesByRequestedRAM[:i], entriesByRequestedRAM[i+1:]...)
				}
			}

		}
	NextNode:
	}

	return res
}
