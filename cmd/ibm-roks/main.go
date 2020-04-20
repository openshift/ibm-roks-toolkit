package main

import (
	"github.com/spf13/cobra"

	"github.com/openshift/ibm-roks-toolkit/pkg/cmd/render"
)

func main() {
	rootCmd := newHypershiftCommand()
	rootCmd.AddCommand(render.NewRenderManifestsCommand())
	rootCmd.Execute()
}

func newHypershiftCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "ibm-roks",
		Short: "IBM ROKS is a toolkit that enables running OpenShift 4.x in a hyperscale manner with many control planes hosted on a central management cluster.",
	}
}
