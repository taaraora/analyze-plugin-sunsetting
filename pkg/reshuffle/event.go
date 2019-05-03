package reshuffle

// TODO: share this package with analyze sunsetting plugin and qbox portal

import (
	"encoding/json"

	"github.com/supergiant/analyze-plugin-sunsetting/pkg/kube"
)

type CommandType string

const ReshufflePodsCommandType CommandType = "ReshufflePodsCommand"
const EsClusterMaintenanceModeCommandType CommandType = "EsClusterMaintenanceModeCommand"

type CommandEnvelope struct {
	// unique ID of command
	ID string `json:"id,omitempty"`
	// type of command
	Type CommandType `json:"commandType,omitempty"`
	// address of service which has produced the command
	// it can be DNS name ir IP address of service
	SourceID string `json:"sourceId,omitempty"`
	// CommandType dependant payload
	Payload json.RawMessage `json:"payload,omitempty"`
}

//nolint
type ReshufflePodsCommand struct {
	ClusterID      string            `json:"clusterId"`
	WorkerNodesIDs []string          `json:"workerNodesIds"`
	AZ             string            `json:"az"`
	ClientConfig   kube.ClientConfig `json:"clientConfig"`
}
