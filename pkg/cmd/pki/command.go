package pki

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/hypershift-toolkit/pkg/cmd/util"
	"github.com/openshift/hypershift-toolkit/pkg/config"
	"github.com/openshift/hypershift-toolkit/pkg/pki"
)

func NewPKICommand() *cobra.Command {
	var outputDir, configFile string
	cmd := &cobra.Command{
		Use:   "pki",
		Short: "Generates PKI artifacts given an output directory",
		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureDir(outputDir)

			params, err := config.ReadFrom(configFile)
			if err != nil {
				log.WithError(err).Fatal("Cannot read config file")
			}

			if err := pki.GeneratePKI(params, outputDir); err != nil {
				log.WithError(err).Fatal("Failed to generate PKI")
			}
		},
	}
	cmd.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "Specify the directory where PKI artifacts should be output")
	cmd.Flags().StringVar(&configFile, "config", defaultConfigFile(), "Specify the config file for this cluster")
	return cmd
}

func defaultOutputDir() string {
	return filepath.Join(util.WorkingDir(), "pki")
}

func defaultConfigFile() string {
	return filepath.Join(util.WorkingDir(), "cluster.yaml")
}
