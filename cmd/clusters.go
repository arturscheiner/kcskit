package cmd

import (
	"github.com/spf13/cobra"
)

var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Manage clusters",
	Long:  "Commands to list and inspect clusters configured in Kaspersky Container Security.",
}

func init() {
	rootCmd.AddCommand(clustersCmd)
}
