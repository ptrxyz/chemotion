package cli

import (
	"github.com/spf13/cobra"
)

func instanceSwitch(given_name string) {
	conf.Set(selector_key, given_name)
	if err := conf.WriteConfig(); err == nil {
		zboth.Info().Msgf("Modified configuration file %s.", conf.ConfigFileUsed())
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to update the selected instance.")
	}
}

var switchInstanceRootCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to an instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		if currentState.quiet {
			if cmd.Flags().Lookup("select-instance").Changed { // this implies a non-interactive run
				instanceSwitch(currentState.name)
			}
		} else {
			confirmInteractive()
			if cmd.Flags().Lookup("select-instance").Changed {
				if selectYesNo("Confirm switching selected instance to "+currentState.name, false) {
					instanceSwitch(currentState.name)
				}
			} else {
				currentState.name = selectInstance("switch to")
				instanceSwitch(currentState.name)
			}
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(switchInstanceRootCmd)
}
