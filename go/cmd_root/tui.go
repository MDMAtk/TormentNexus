package cmd

import (
	"log"

	"gitlab.com/robertpelloni/HyperNexus/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the highly-advanced Native Go BubbleTea Orchestrator Interface",
	Long:  "Starts the interactive visual TUI leveraging the Native Node.js-Parity HyperNexus Orchestrator engine complete with AutoDrive autonomy.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("[Boot] Spinning up Native HyperNexus TUI...")
		tui.StartREPL()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
