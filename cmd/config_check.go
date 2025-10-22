package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	ctrl "github.com/arturscheiner/kcskit/internal/controller"
	"github.com/arturscheiner/kcskit/internal/model"
)

var invalidCert bool
var outputFlag string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if kcskit is configured and test connection to endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := ctrl.LoadConfig()
		if err != nil {
			fmt.Println("not configured:", err)
			os.Exit(1)
		}
		if err := ctrl.ValidateConfig(cfg); err != nil {
			fmt.Println("not configured:", err)
			os.Exit(1)
		}

		body, err := ctrl.TestConfigConnection(cfg, invalidCert)
		if err != nil {
			fmt.Println("connection test failed:", err)
			if body != "" {
				fmt.Println("response body:", body)
			}
			os.Exit(1)
		}

		// if user requested JSON output, print pretty JSON and exit
		if outputFlag == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
				// fallback to raw body
				fmt.Println(body)
			} else {
				fmt.Println(pretty.String())
			}
			return
		}

		var hr model.HealthResponse
		if err := json.Unmarshal([]byte(body), &hr); err != nil {
			fmt.Println("failed to parse health JSON:", err)
			fmt.Println("response body:", body)
			os.Exit(1)
		}

		// print tabbed table: Name | Pod | Status | Version | Error
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Name\tPod\tStatus\tVersion\tError")
		for _, it := range hr.Items {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", it.ComponentName, it.PodName, it.Status, it.Version, it.ErrorMessage)
		}
		_ = w.Flush()
	},
}

func init() {
	configCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&invalidCert, "invalid-cert", "i", false, "ignore TLS certificate validation when performing connection test")
	checkCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "output format (\"json\" for raw JSON output). Default: tabbed table")
}
