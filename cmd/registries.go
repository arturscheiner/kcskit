package cmd

import (
	"github.com/spf13/cobra"
)

var registriesCmd = &cobra.Command{
	Use:   "registries",
	Short: "Manage image registries integrations",
	Long:  "Manage image registries integrations configured in Kaspersky Container Security (list, show, ...).",
}

func init() {
	rootCmd.AddCommand(registriesCmd)
}
