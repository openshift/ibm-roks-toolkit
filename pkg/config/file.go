package config

import (
	"os"

	"github.com/ghodss/yaml"

	"github.com/openshift/ibm-roks-toolkit/pkg/api"
)

func ReadFrom(fileName string) (*api.ClusterParams, error) {
	result := api.NewClusterParams()
	b, err := os.ReadFile(fileName) // #nosec G304 We control the contents of any files read by this function
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
