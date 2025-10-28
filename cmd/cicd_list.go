package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	ctrl "github.com/arturscheiner/kcskit/internal/controller"
	"github.com/spf13/cobra"
)

var (
	cicdOutput      string
	flagCicdPage         int
	flagCicdLimit        int
	flagCicdSort         string
	flagCicdBy           string
	flagCicdBuildNumber  string
	flagCicdBuildPipeline string
)

var cicdListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CI/CD scans",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := ctrl.LoadConfig()
		if err != nil {
			return fmt.Errorf("not configured: %w", err)
		}
		if err := ctrl.ValidateConfig(cfg); err != nil {
			return fmt.Errorf("not configured: %w", err)
		}

		page := strconv.Itoa(flagCicdPage)
		limit := strconv.Itoa(flagCicdLimit)

		items, body, err := ctrl.ListCicd(cfg, InvalidCert, page, limit, flagCicdSort, flagCicdBy, flagCicdBuildNumber, flagCicdBuildPipeline)
		if err != nil {
			fmt.Println("failed to list clusters:", err)
			if body != "" {
				fmt.Println("response body:", body)
			}
			os.Exit(1)
		}

		if cicdOutput == "json" {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
				fmt.Println(body)
			} else {
				fmt.Println(pretty.String())
			}
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tArtifact\tRisk\tStatus")
		for _, it := range items.Items {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", it.ID, it.ArtifactName, it.RiskRating, it.Status)
		}
		_ = w.Flush()
		return nil
	},
}

func init() {
	cicdCmd.AddCommand(cicdListCmd)

	cicdListCmd.Flags().IntVar(&flagCicdPage, "page", 1, "The page number to retrieve for paginated results.")
	cicdListCmd.Flags().IntVar(&flagCicdLimit, "limit", 50, "The number of items to include per page.")
	cicdListCmd.Flags().StringVar(&flagCicdSort, "sort", "createdAt", "Sort by value (createdAt|updatedAt|artifactName|name|artifactType|status|riskRating)")
	cicdListCmd.Flags().StringVar(&flagCicdBy, "by", "desc", "Sort by order (asc|desc)")
	cicdListCmd.Flags().StringVar(&flagCicdBuildNumber, "build-number", "", "Filter by build number.")
	cicdListCmd.Flags().StringVar(&flagCicdBuildPipeline, "build-pipeline", "", "Filter by build pipeline.")

	cicdListCmd.Flags().StringVarP(&cicdOutput, "output", "o", "", "output format (\"json\" for raw JSON output). Default: tabbed table")
}
