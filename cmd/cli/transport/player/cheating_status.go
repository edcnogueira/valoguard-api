package player

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/edcnogueira/valoguard-api/internal/models"
	"github.com/edcnogueira/valoguard-api/internal/providers/henrik"
	"github.com/edcnogueira/valoguard-api/internal/service/analysisservice"
)

var name string
var tag string
var region string

var (
	analyzeCmd = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze a player cheating probability",
		RunE: func(cmd *cobra.Command, args []string) error {
			if apiKey == "" {
				apiKey = os.Getenv("HENRIK_API_KEY")
			}
			if apiKey == "" {
				return fmt.Errorf("define HENRIK_API_KEY or pass --api-key")
			}
			if name == "" || tag == "" {
				return fmt.Errorf("flags --name and --tag are required")
			}

			hClient := henrik.New(http.DefaultClient, apiKey)
			svc := analysisservice.New(hClient)

			ctx := context.Background()
			resp, err := svc.AnalyzePlayer(ctx, &models.AnalyzeRequest{
				Name:   name,
				Tag:    tag,
				Region: region,
			})
			if err != nil {
				return err
			}

			b, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}

			fmt.Println(string(b))
			return nil
		},
	}
)

func init() {
	analyzeCmd.Flags().StringVar(&name, "name", "", "Riot ID name (e.g., player)")
	analyzeCmd.Flags().StringVar(&tag, "tag", "", "Riot ID tag (e.g., 1234)")
	analyzeCmd.Flags().StringVar(&region, "region", "", "Region (default in service: br)")
	_ = analyzeCmd.MarkFlagRequired("name")
	_ = analyzeCmd.MarkFlagRequired("tag")

	cmdCheatingStatus.AddCommand(analyzeCmd)
}
