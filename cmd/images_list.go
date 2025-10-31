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

var imagesOutput string

// flags for /v1/images/registry
var (
	flagPage             int
	flagLimit            int
	flagSort             string
	flagBy               string
	flagScopes           []string
	flagName             string
	flagRegistry         string
	flagRepositoriesWith string
	flagScannedAt        string
	flagRisks            []string
)

var imagesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List images for a registry (scan results)",
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

		// build query params
		v := url.Values{}
		v.Set("page", strconv.Itoa(flagPage))
		v.Set("limit", strconv.Itoa(flagLimit))
		if flagSort != "" {
			v.Set("sort", flagSort)
		}
		if flagBy != "" {
			v.Set("by", flagBy)
		}
		for _, s := range flagScopes {
			v.Add("scopes[]", s)
		}
		if flagName != "" {
			v.Set("name", flagName)
		}
		if flagRegistry != "" {
			v.Set("registry", flagRegistry)
		}
		if flagRepositoriesWith != "" {
			v.Set("repositoriesWith", flagRepositoriesWith)
		}
		if flagScannedAt != "" {
			v.Set("scannedAt", flagScannedAt)
		}
		for _, r := range flagRisks {
			v.Add("risks[]", r)
		}
		rawQuery := v.Encode()

		items, body, endpoint, err := ctrl.ListImages(cfg, InvalidCert, rawQuery)
		if err != nil {
			fmt.Println("failed to list images:", err)
			if body != "" {
				fmt.Println("response body:", body)
			}
			os.Exit(1)
		}

		if imagesOutput == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
				fmt.Println(body)
			} else {
				fmt.Println(pretty.String())
			}
			return
		} else if imagesOutput == "ollama" {
			var risks []string
			for _, item := range items {
				risks = append(risks, item.RiskRating)
			}

			header := model.OllamaHeader{
				Command:     strings.Join(os.Args, " "),
				Cluster:     "",
				Risk:        strings.Join(risks, ", "),
				ReportTitle: "Kaspersky Container Security Image Assessment Report.",
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

		// default: tabbed table with columns: ID, Name, Registry, Risk
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tName\tRegistry\tRisk")
		for _, it := range items {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", it.ID, it.Name, it.ImageRegistryName, it.RiskRating)
		}
		_ = w.Flush()
	},
}

func init() {
	imagesCmd.AddCommand(imagesListCmd)

	imagesListCmd.Flags().IntVar(&flagPage, "page", 1, "page number")
	imagesListCmd.Flags().IntVar(&flagLimit, "limit", 50, "items per page")
	imagesListCmd.Flags().StringVar(&flagSort, "sort", "", "sort by (name|riskRating)")
	imagesListCmd.Flags().StringVar(&flagBy, "by", "asc", "sort order (asc|desc)")
	imagesListCmd.Flags().StringSliceVar(&flagScopes, "scopes", nil, "filter by scopes (repeatable)")
	imagesListCmd.Flags().StringVar(&flagName, "name", "", "filter by registry name")
	imagesListCmd.Flags().StringVar(&flagRegistry, "registry", "", "filter by registry ID")
	imagesListCmd.Flags().StringVar(&flagRepositoriesWith, "repositoriesWith", "", "filter by repository risk rating (compliant|non-compliant|error|process)")
	imagesListCmd.Flags().StringVar(&flagScannedAt, "scannedAt", "", "filter by scan timeframe (hour|day|week)")
	imagesListCmd.Flags().StringSliceVar(&flagRisks, "risks", nil, "filter by risk types (malware|vulnerabilities|sensitive-data|misconfiguration) (repeatable)")

	imagesListCmd.Flags().StringVarP(&imagesOutput, "output", "o", "", "output format (\"json\" for raw JSON output, \"ollama\" to send to Ollama). Default: tabbed table")
}
