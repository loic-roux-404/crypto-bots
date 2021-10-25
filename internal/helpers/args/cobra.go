package args

import (
	"os"

	"github.com/spf13/cobra"
)

// ShowFullHelp if no args
func ShowFullHelp(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}
}