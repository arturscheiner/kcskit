package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	ctrl "github.com/arturscheiner/kcskit/internal/controller"
	"github.com/arturscheiner/kcskit/internal/model"
)

var clustersOutput string

var (
	flagClusterPage   int
	flagClusterLimit  int
	flagClusterSort   string
	flagClusterBy     string
	flagClusterScopes []string
)

var clustersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
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

		v := url.Values{}
		v.Set("page", strconv.Itoa(flagClusterPage))
		v.Set("limit", strconv.Itoa(flagClusterLimit))
		if flagClusterSort != "" {
			v.Set("sort", flagClusterSort)
		}
		if flagClusterBy != "" {
			v.Set("by", flagClusterBy)
		}
		for _, s := range flagClusterScopes {
			v.Add("scopes[]", s)
		}
		rawQuery := v.Encode()

		items, body, endpoint, err := ctrl.ListClusters(cfg, InvalidCert, rawQuery)
		if err != nil {
			fmt.Println("failed to list clusters:", err)
			if body != "" {
				fmt.Println("response body:", body)
			}
			os.Exit(1)
		}

		if clustersOutput == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
				fmt.Println(body)
			} else {
				fmt.Println(pretty.String())
			}
			return
		} else if clustersOutput == "ollama" {
			var clusterNames []string
			var risks []string
			for _, item := range items {
				clusterNames = append(clusterNames, item.ClusterName)
				risks = append(risks, item.RiskRating)
			}

			header := model.OllamaHeader{
				Command:     strings.Join(os.Args, " "),
				Cluster:     strings.Join(clusterNames, ", "),
				Risk:        strings.Join(risks, ", "),
				ReportTitle: "Kaspersky Container Security Cluster Assessment Report.",
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tOrchestrator\tNamespaces\tRisk")
		for _, it := range items {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", it.ID, it.ClusterName, it.Orchestrator, it.Namespaces, it.RiskRating)
		}
		_ = w.Flush()
	},
}

func init() {
	clustersCmd.AddCommand(clustersListCmd)

	clustersListCmd.Flags().IntVar(&flagClusterPage, "page", 1, "page number")
	clustersListCmd.Flags().IntVar(&flagClusterLimit, "limit", 50, "items per page")
	clustersListCmd.Flags().StringVar(&flagClusterSort, "sort", "clusterName", "sort by (clusterName|orchestrator|namespaces|riskRating)")
	clustersListCmd.Flags().StringVar(&flagClusterBy, "by", "asc", "sort order (asc|desc)")
	clustersListCmd.Flags().StringSliceVar(&flagClusterScopes, "scopes", nil, "filter by scopes (repeatable)")

	clustersListCmd.Flags().StringVarP(&clustersOutput, "output", "o", "", "output format (\"json\" for raw JSON output, \"ollama\" to send to Ollama). Default: tabbed table")
}
