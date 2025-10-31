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

var registriesOutput string
var registriesInvalidCert bool

var registriesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured image registries",
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

		items, body, endpoint, err := ctrl.ListRegistries(cfg, registriesInvalidCert)
		if err != nil {
			fmt.Println("failed to list registries:", err)
			if body != "" {
				fmt.Println("response body:", body)
			}
			os.Exit(1)
		}

		if registriesOutput == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
				fmt.Println(body)
			} else {
				fmt.Println(pretty.String())
			}
			return
		} else if registriesOutput == "ollama" {
			header := model.OllamaHeader{
				Command:     strings.Join(os.Args, " "),
				Cluster:     "",
				Risk:        "",
				ReportTitle: "Kaspersky Container Security Registries Assessment Report.",
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

		// default: tabbed table with columns: ID, Name, Type, Url
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tType\tUrl")
		for _, it := range items {
			urlToShow := it.ApiUrl
			if urlToShow == "" {
				urlToShow = it.RegistryUrl
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", it.ID, it.RegistryName, it.RegistryType, urlToShow)
		}
		_ = w.Flush()
	},
}

func init() {
	registriesCmd.AddCommand(registriesListCmd)
	registriesListCmd.Flags().StringVarP(&registriesOutput, "output", "o", "", "output format (\"json\" for raw JSON output, \"ollama\" to send to Ollama). Default: tabbed table")
	registriesListCmd.Flags().BoolVarP(&registriesInvalidCert, "invalid-cert", "i", false, "ignore TLS certificate validation when performing API requests")
}
