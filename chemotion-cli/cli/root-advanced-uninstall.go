package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallAdvancedRootCmd = &cobra.Command{
	Use:   "uninstall",
	Args:  cobra.NoArgs,
	Short: fmt.Sprintf("uninstall %s completely", nameCLI),
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		if currentState.quiet {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("For security reasons, this command cannot be run in silent mode.")
		}
		confirmInteractive()
		if selectYesNo("Are you sure you want to uninstall "+nameCLI, false) {
			// TODO: do actual uninstall
		} else {
			zboth.Info().Msgf("Nothing was done.")
		}
	},
}

func init() {
	advancedRootCmd.AddCommand(uninstallAdvancedRootCmd)
}
