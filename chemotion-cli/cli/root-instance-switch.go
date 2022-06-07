package cli

import (
	"github.com/spf13/cobra"
)

func instanceSwitch() {
	conf.Set(selector_key, currentState.name)
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
				instanceSwitch()
			}
		} else {
			confirmInteractive()
			if cmd.Flags().Lookup("select-instance").Changed {
				if selectYesNo("Confirm switching selected instance to "+currentState.name, false) {
					instanceSwitch()
				}
			} else {
				currentState.name = selectInstance("switch to")
				instanceSwitch()
			}
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(switchInstanceRootCmd)
}
