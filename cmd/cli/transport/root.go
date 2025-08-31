package transport

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/edcnogueira/valoguard-api/cmd/cli/transport/player"
)

var rootCmd = &cobra.Command{
	Use:   "valoguard",
	Short: "ValoGuard CLI - analyze Valorant players locally",
}

func init() {
	rootCmd.AddCommand(player.CmdPlayer)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
