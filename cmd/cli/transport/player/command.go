package player

import "github.com/spf13/cobra"

// defaultAPIKey can be set at build time via:
//   go build -ldflags "-X 'github.com/edcnogueira/valoguard-api/cmd/cli/transport/player.defaultAPIKey=YOUR_KEY'" ./cmd/cli
var defaultAPIKey string

var apiKey string

var (
	CmdPlayer = &cobra.Command{
		Use:   "player",
		Short: "Player related commands",
	}

	cmdCheatingStatus = &cobra.Command{
		Use:   "cheating",
		Short: "Get cheating status of a player",
	}
)

func init() {
	CmdPlayer.AddCommand(cmdCheatingStatus)

	CmdPlayer.PersistentFlags().StringVar(&apiKey, "api-key", defaultAPIKey, "Henrik API key (or set HENRIK_API_KEY env var)")
}
