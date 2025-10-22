/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags. Default is "dev".
var Version = "dev"

// optional: include commit/date fields if desired
var Commit = ""
var Date = ""

// Global flag: ignore TLS certificate validation
var InvalidCert bool

var rootCmd = &cobra.Command{
	Use:   "kcskit",
	Short: "kcskit — lightweight CLI for Kaspersky Container Security (KCS)",
	Long: `kcskit is a small command-line utility to interact with Kaspersky Container Security (KCS).

For help on a specific command:
  kcskit <command> --help
`,
	Version: Version,
	// nice --version output
	// VersionTemplate: "kcskit version: {{.Version}}{{if .Commit}}\ncommit: {{.Commit}}{{end}}{{if .Date}}\nbuilt:  {{.Date}}{{end}}\n",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// register global persistent flag for ignoring invalid TLS certificates
	rootCmd.PersistentFlags().BoolVarP(&InvalidCert, "invalid-cert", "i", false, "ignore TLS certificate validation for all commands (use with caution)")

	// set a version template (VersionTemplate field is unexported; use setter)
	rootCmd.SetVersionTemplate("kcskit version: {{.Version}}\n")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kcskit.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
