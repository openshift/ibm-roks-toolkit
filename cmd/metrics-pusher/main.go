package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/openshift/ibm-roks-toolkit/pkg/metricspusher"
)

func main() {
	if err := newMetricsPusherCommand().Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "Error: %v", err)
	}
}

func newMetricsPusherCommand() *cobra.Command {

	metricsPusher := metricspusher.MetricsPusher{}
	cmd := &cobra.Command{
		Use:          "metrics-pusher",
		Short:        "Pushes metrics to a given endpoint inside a cluster",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.SetLogger(zap.New())
			metricsPusher.Log = ctrl.Log.WithName("metrics-pusher")
			return metricsPusher.Run()
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&metricsPusher.SourceURL, "source-url", "", "Source URL to poll for metrics")
	flags.StringVar(&metricsPusher.SourcePath, "source-path", "", "Source path to poll for metrics (using kubeconfig)")
	flags.StringVar(&metricsPusher.DestinationPath, "destination-path", "", "URL path to use as a destination on the target cluster")
	flags.StringVar(&metricsPusher.Kubeconfig, "kubeconfig", "", "Kubeconfig file for destination cluster")
	flags.DurationVar(&metricsPusher.Frequency, "frequency", 30*time.Second, "Frequency with which to push metrics")
	flags.StringVar(&metricsPusher.Clientca, "client-ca-file", "", "CA certificate for client certificate authentication to the server")
	return cmd
}
