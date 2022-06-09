package cli

import (
	"github.com/spf13/cobra"
)

func instanceRestart(given_name string) {
	instanceStop(given_name)
	instanceStart(given_name)
}

var restartInstanceRootCmd = &cobra.Command{
	Use:   "restart",
	Args:  cobra.NoArgs,
	Short: "Restart (the selected) instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		instanceRestart(currentState.name)
	},
	// TODO: add a force restart flag
}

func init() {
	instanceRootCmd.AddCommand(restartInstanceRootCmd)
}
