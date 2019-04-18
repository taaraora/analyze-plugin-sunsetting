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

func CheckEachPodOneByOne(unsortedEntries []*InstanceEntry) []InstanceEntry {
	var entriesByWastedRAM = NewSortedEntriesByWastedRAM(unsortedEntries)
	var entriesByRequestedRAM = NewSortedEntriesByRequestedRAM(unsortedEntries)
	var res = make([]InstanceEntry, 0)

	for _, maxWastedRAMEntry := range entriesByWastedRAM {
		// sort pods in descending order by requested memory to more effective pods packing on nodes
		podsRR := maxWastedRAMEntry.WorkerNode.PodsResourceRequirements
		sort.Slice(
			podsRR,
			func(i, j int) bool {
				return podsRR[i].MemoryReqs > podsRR[j].MemoryReqs
			},
		)

		// check
		for i := 0; i < len(maxWastedRAMEntry.WorkerNode.PodsResourceRequirements); i++ {
			var podRR = maxWastedRAMEntry.WorkerNode.PodsResourceRequirements[i]
			for _, maxRequestedRAMEntry := range entriesByRequestedRAM {
				if maxRequestedRAMEntry.RAMWasted() >= podRR.MemoryReqs &&
					maxRequestedRAMEntry.CPUWasted() >= podRR.CPUReqs {
					// we can move the pod
					// delete it from  maxWastedRAMEntry
					maxWastedRAMEntry.WorkerNode.PodsResourceRequirements = append(
						maxWastedRAMEntry.WorkerNode.PodsResourceRequirements[:i],
						maxWastedRAMEntry.WorkerNode.PodsResourceRequirements[i+1:]...,
					)

					maxWastedRAMEntry.WorkerNode.RefreshTotals()
					//and add to maxRequestedRAMEntry
					maxRequestedRAMEntry.WorkerNode.PodsResourceRequirements = append(
						maxRequestedRAMEntry.WorkerNode.PodsResourceRequirements,
						podRR,
					)
					maxRequestedRAMEntry.WorkerNode.RefreshTotals()
					// if last pod was moved out exit from loop
					// and basically because of i == o we also exit from outer loop
					if len(maxWastedRAMEntry.WorkerNode.PodsResourceRequirements) == 0 {
						break
					}
					// reset index to zero because after maxWastedRAMEntry.WorkerNode.PodsResourceRequirements slice
					// modification we need to start from scratch
					i = 0
					continue
				}
			}
		}
		if len(maxWastedRAMEntry.WorkerNode.PodsResourceRequirements) == 0 {
			res = append(res, *maxWastedRAMEntry)
		}
	}

	return res
}
