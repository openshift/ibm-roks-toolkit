package api

import (
	"github.com/google/uuid"
)

// NewClusterParams returns a new default cluster params struct
func NewClusterParams() *ClusterParams {
	p := &ClusterParams{}
	p.DefaultFeatureGates = []string{
		"APIPriorityAndFairness=true",
		"SupportPodPidsLimit=true",
		"RotateKubeletServerCertificate=true",
		"LegacyNodeRoleBehavior=false",
		"NodeDisruptionExclusion=true",
		"SCTPSupport=true",
		"ServiceNodeExclusion=true",
	}
	p.ImageRegistryHTTPSecret = uuid.New().String()
	return p
}
