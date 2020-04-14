package ignition

import (
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/hypershift-toolkit/pkg/cmd/util"
	"github.com/openshift/hypershift-toolkit/pkg/config"
	"github.com/openshift/hypershift-toolkit/pkg/ignition"
)

func NewIgnitionCommand() *cobra.Command {
	var pkiDir, outputDir, configFile, pullSecretFile, sshPublicKeyFile string
	cmd := &cobra.Command{
		Use:   "ignition",
		Short: "Generates an ignition file to be used by RHCOS workers on boot",
		Run: func(cmd *cobra.Command, args []string) {
			util.EnsureDir(outputDir)

			params, err := config.ReadFrom(configFile)
			if err != nil {
				log.WithError(err).Fatal("Cannot read config file")
			}

			sshPublicKey, err := ioutil.ReadFile(sshPublicKeyFile)
			if err != nil {
				log.WithError(err).Fatal("Cannot read SSH public key file")
			}

			if err := ignition.GenerateIgnition(params, sshPublicKey, pullSecretFile, pkiDir, outputDir); err != nil {
				log.WithError(err).Fatal("Failed to generate ignition")
			}
		},
	}
	cmd.Flags().StringVar(&outputDir, "output-dir", defaultOutputDir(), "Specify the directory where the ignition file should be output")
	cmd.Flags().StringVar(&configFile, "config", defaultConfigFile(), "Specify the config file for this cluster")
	cmd.Flags().StringVar(&sshPublicKeyFile, "ssh-public-key", defaultSSHPublicKeyFile(), "Specify the config file for this cluster")
	cmd.Flags().StringVar(&pkiDir, "pki-dir", defaultPKIDir(), "Specify the directory containing PKI files")
	cmd.Flags().StringVar(&pullSecretFile, "pull-secret", defaultPullSecretFile(), "Specify the config file for this cluster")
	return cmd
}

func defaultOutputDir() string {
	return filepath.Join(util.WorkingDir(), "ignition")
}

func defaultConfigFile() string {
	return filepath.Join(util.WorkingDir(), "cluster.yaml")
}

func defaultSSHPublicKeyFile() string {
	return filepath.Join(util.WorkingDir(), "ssh-public-key")
}

func defaultPKIDir() string {
	return filepath.Join(util.WorkingDir(), "pki")
}

func defaultPullSecretFile() string {
	return filepath.Join(util.WorkingDir(), "pull-secret.txt")
}
