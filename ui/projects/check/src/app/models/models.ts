export interface WebComponentInfo {
  selector: string
  webComponentName: string
  pluginName: string
  pluginVersion: string
}

export interface CheckResult {
  checkStatus: CheckStatus
  completedAt: Date
  //TODO: explore better solution
  description: Array<NodeCheckResult> | string
  executionStatus: string
  id: string
  name: string
  possibleActions: PluginAction
}

export enum CheckStatus {
  RED = "RED",
  YELLOW = "YELLOW",
  GREEN = "GREEN"
}

export interface PluginAction {
  description: string
  id: string
  name: string
}

export interface NodeCheckResult {
  cloudProvider: CloudProvider
  price: Price
  kube: KubeNode

}

export interface CloudProvider {
  instanceId: string
  instanceType: string
}

export interface Price {
  "instanceType": string
  "memory": string
  "vcpu": string
  "unit": string
  "currency": string
  "valuePerUnit": string
  "usageType": string
  "tenancy": string
}

export interface KubeNode {
  allocatableCpu: number
  allocatableMemory: number
  cpuLimits: number
  cpuRequests: number
  fractionCpuLimits: number
  fractionCpuRequests: number
  fractionMemoryLimits: number
  fractionMemoryRequests: number
  instanceId: string
  isRecommendedToSunset: boolean
  memoryLimits: number
  memoryRequests: number
  name: string
  region: string
  podsResourceRequirements: PodResourceRequirement
}

export interface PodResourceRequirement {
  podName: string
  cpuRequests: number
  cpuLimits: number
  memoryRequests: number
  memoryLimits: number
}
