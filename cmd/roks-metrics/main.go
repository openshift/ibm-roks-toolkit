package main

import (
	"flag"

	"github.com/spf13/cobra"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/openshift/ibm-roks-toolkit/pkg/metrics"
)

func main() {
	log.SetLogger(zap.New(zap.UseDevMode(true)))
	metricsLog := ctrl.Log.WithName("metrics")
	if err := newMetricsServerCommand().Execute(); err != nil {
		metricsLog.Error(err, "Operator failed")
	}
}

func newMetricsServerCommand() *cobra.Command {
	metricsServer := metrics.MetricsServer{}
	cmd := &cobra.Command{
		Use:          "roks-metrics",
		Short:        "Metrics server for ROKS clusters",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return metricsServer.Run()
		},
	}
	flags := cmd.Flags()
	flags.AddGoFlagSet(flag.CommandLine)
	flags.StringVar(&metricsServer.ListenAddress, "listen", ":8443", "Address/port to listen on")
	flags.StringVar(&metricsServer.CertFile, "cert", "/var/run/secrets/serving-cert/tls.crt", "Serving certificate file")
	flags.StringVar(&metricsServer.KeyFile, "key", "/var/run/secrets/serving-cert/tls.key", "Serving certificate key file")
	return cmd
}
