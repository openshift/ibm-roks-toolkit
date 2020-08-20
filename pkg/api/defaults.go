package api

import (
	"github.com/google/uuid"
)

// NewClusterParams returns a new default cluster params struct
func NewClusterParams() *ClusterParams {
	p := &ClusterParams{}
	p.DefaultFeatureGates = []string{
		"SupportPodPidsLimit=true",
		"LocalStorageCapacityIsolation=false",
		"RotateKubeletServerCertificate=true",
	}
	p.ImageRegistryHTTPSecret = uuid.New().String()
	return p
}
