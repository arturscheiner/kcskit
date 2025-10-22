/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	ctrl "github.com/arturscheiner/kcskit/internal/controller"
)

var tokenFlag string
var endpointFlag string
var caCertFlag string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage kcskit configuration (save token/endpoint)",
	Long: `Save token and endpoint into $HOME/.kcskit/config

CA certificate examples (ca_cert flag accepts literal PEM, a file path, or '-' to read from stdin):

# Save PEM from a file
kcskit config --ca_cert /path/to/ca.pem --endpoint https://kcs.example.com/api/ --token kcs_...

# Pipe PEM via stdin (useful in scripts)
cat /path/to/ca.pem | kcskit config --ca_cert - --endpoint https://kcs.example.com/api/ --token kcs_...

# Provide PEM inline (shell-escaped)
kcskit config --ca_cert "$(cat /path/to/ca.pem)" --endpoint https://kcs.example.com/api/ --token kcs_...

The ca_cert value is stored as the 'ca_cert' field in the YAML config at $HOME/.kcskit/config.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no flags provided, show help
		if tokenFlag == "" && endpointFlag == "" && caCertFlag == "" {
			_ = cmd.Help()
			return
		}

		// support: literal PEM, file path, or "-" for stdin
		caCertContent := caCertFlag
		if caCertFlag != "" {
			if caCertFlag == "-" {
				b, err := io.ReadAll(os.Stdin)
				if err != nil {
					fmt.Println("error reading ca_cert from stdin:", err)
					return
				}
				caCertContent = string(b)
			} else if fi, err := os.Stat(caCertFlag); err == nil && !fi.IsDir() {
				b, err := os.ReadFile(caCertFlag)
				if err != nil {
					fmt.Println("error reading ca_cert file:", err)
					return
				}
				caCertContent = string(b)
			}
			// otherwise: treat caCertFlag as the literal PEM text (backwards compatible)
		}

		if err := ctrl.SaveConfig(tokenFlag, endpointFlag, caCertContent); err != nil {
			fmt.Println("error writing config file:", err)
		} else {
			fmt.Println("configuration saved")
		}
	},
}

func init() {
	// add flags to config command
	configCmd.Flags().StringVar(&tokenFlag, "token", "", "the API token value defined in the KCS web console's user my profile")
	configCmd.Flags().StringVar(&endpointFlag, "endpoint", "", "the API endpoint URL, e.g. https://kcs.example.com/api/")
	configCmd.Flags().StringVar(&caCertFlag, "ca_cert", "", "CA certificate PEM text or path to a PEM file. Use '-' to read from stdin.")
	rootCmd.AddCommand(configCmd)
}
