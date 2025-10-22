package cmd

import (
	"github.com/spf13/cobra"
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Manage images (scan results) operations",
	Long:  "Commands to list and inspect images scanning results from KCS.",
}

func init() {
	rootCmd.AddCommand(imagesCmd)
}
