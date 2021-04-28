package render

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/ibm-roks-toolkit/pkg/cmd/util"
	"github.com/openshift/ibm-roks-toolkit/pkg/config"
	"github.com/openshift/ibm-roks-toolkit/pkg/render"
)

type RenderManifestsOptions struct {
	OutputDir       string
	ConfigFile      string
	PullSecretFile  string
	IncludeRegistry bool
}

func NewRenderManifestsCommand() *cobra.Command {
	opt := &RenderManifestsOptions{}
	cmd := &cobra.Command{
		Use: "render",
		Run: func(cmd *cobra.Command, args []string) {
			if err := opt.Run(); err != nil {
				log.WithError(err).Fatal("Error occurred rendering manifests")
			}
		},
	}
	cmd.Flags().StringVar(&opt.OutputDir, "output-dir", defaultManifestsDir(), "Specify the directory where manifest files should be output")
	cmd.Flags().StringVar(&opt.ConfigFile, "config", defaultConfigFile(), "Specify the config file for this cluster")
	cmd.Flags().StringVar(&opt.PullSecretFile, "pull-secret", defaultPullSecretFile(), "Specify the config file for this cluster")
	cmd.Flags().BoolVar(&opt.IncludeRegistry, "include-registry", false, "If true, includes a default registry config to deploy into the user cluster")
	return cmd
}

func (o *RenderManifestsOptions) Run() error {
	util.EnsureDir(o.OutputDir)
	params, err := config.ReadFrom(o.ConfigFile)
	if err != nil {
		log.WithError(err).Fatalf("Error occurred reading configuration")
	}
	externalOauth := params.ExternalOauthPort != 0
	if len(params.ExternalOauthDNSName) == 0 {
		params.ExternalOauthDNSName = params.ExternalAPIDNSName
	}
	err = render.RenderClusterManifests(params, o.PullSecretFile, o.OutputDir, externalOauth, o.IncludeRegistry, params.KonnectivityEnabled)
	if err != nil {
		return err
	}
	return nil
}

func defaultManifestsDir() string {
	return filepath.Join(util.WorkingDir(), "manifests")
}

func defaultConfigFile() string {
	return filepath.Join(util.WorkingDir(), "cluster.yaml")
}

func defaultPullSecretFile() string {
	return filepath.Join(util.WorkingDir(), "pull-secret.txt")
}

func defaultPKIDir() string {
	return filepath.Join(util.WorkingDir(), "pki")
}
