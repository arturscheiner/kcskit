package cmd

import (
	"github.com/spf13/cobra"
)

var cicdCmd = &cobra.Command{
	Use:   "cicd",
	Short: "Manage CI/CD scans",
	Long:  "Commands to list and inspect CI/CD scans in Kaspersky Container Security.",
}

func init() {
	rootCmd.AddCommand(cicdCmd)
}
