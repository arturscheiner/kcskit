// filepath: /home/arturscheiner/Development/kaspersky/kcskit/cmd/images_scan.go
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	ctrl "github.com/arturscheiner/kcskit/internal/controller"
	"github.com/arturscheiner/kcskit/internal/model"
)

var imagesScanOutput string
var flagArtifact string
var flagRegistryID string

var imagesScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Create a new scanning job for an artifact in a registry",
	Long:  "Create a manual scan job. Both --artifact and --registry are required.",
	Run: func(cmd *cobra.Command, args []string) {
		if flagArtifact == "" || flagRegistryID == "" {
			fmt.Fprintln(os.Stderr, "error: --artifact and --registry are required")
			_ = cmd.Help()
			os.Exit(1)
		}

		cfg, err := ctrl.LoadConfig()
		if err != nil {
			fmt.Println("not configured:", err)
			os.Exit(1)
		}
		if err := ctrl.ValidateConfig(cfg); err != nil {
			fmt.Println("not configured:", err)
			os.Exit(1)
		}

		job, body, endpoint, err := ctrl.CreateScan(cfg, InvalidCert, flagArtifact, flagRegistryID)
		if err != nil {
			fmt.Println("failed to create scan:", err)
			if body != "" {
				fmt.Println("response body:", body)
			}
			os.Exit(1)
		}

		if imagesScanOutput == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
				fmt.Println(body)
			} else {
				fmt.Println(pretty.String())
			}
			return
		} else if imagesScanOutput == "ollama" {
			header := model.OllamaHeader{
				Command:     strings.Join(os.Args, " "),
				Cluster:     "",
				Risk:        job.Status,
				ReportTitle: "Kaspersky Container Security Image Scan Assessment Report.",
				ApiEndpoint: endpoint,
			}
			response, err := ctrl.SendToOllama(body, header)
			if err != nil {
				fmt.Println("failed to send to ollama:", err)
				os.Exit(1)
			}
			fmt.Println(response)
			return
		}

		// default tabbed table: ID | Artifact | Scanner | Status
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tArtifact\tScanner\tStatus")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", job.ID, job.ArtifactName, job.ScannerName, job.Status)
		_ = w.Flush()
	},
}

func init() {
	imagesCmd.AddCommand(imagesScanCmd)

	imagesScanCmd.Flags().StringVar(&flagArtifact, "artifact", "", "artifact reference, e.g. nginx:latest (required)")
	imagesScanCmd.Flags().StringVar(&flagRegistryID, "registry", "", "registry ID where the artifact resides (required)")
	imagesScanCmd.Flags().StringVarP(&imagesScanOutput, "output", "o", "", "output format (\"json\" for raw JSON output, \"ollama\" to send to Ollama). Default: tabbed table")

	_ = imagesScanCmd.MarkFlagRequired("artifact")
	_ = imagesScanCmd.MarkFlagRequired("registry")
}
