package cli

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var advancedCmdTable = make(cmdTable)

// Backbone for system-related commands
var advancedRootCmd = &cobra.Command{
	Use:       "advanced",
	Short:     "Perform advanced actions related to system and " + nameCLI,
	Args:      cobra.NoArgs,
	ValidArgs: maps.Keys(advancedCmdTable),
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Flag("selected-instance").Changed {
			zboth.Warn().Msgf("The `-i` flag is not supported for the `advanced` command and its subcommands.")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive(true)
		acceptedOpts := []string{"info"}
		if updateRequired() {
			acceptedOpts = append(acceptedOpts, "update cli")
			advancedCmdTable["update cli"] = updateSelfAdvancedRootCmd.Run
		}
		acceptedOpts = append(acceptedOpts, []string{"pull image", "uninstall"}...)
		if cmd.Use == cmd.CalledAs() { // || elementInSlice(cmd.CalledAs(), &cmd.Aliases) > -1 { { // there are no aliases at the moment
			acceptedOpts = append(acceptedOpts, "exit")
		} else {
			acceptedOpts = append(acceptedOpts, []string{"back", "exit"}...)
			advancedCmdTable["back"] = cmd.Run
		}
		advancedCmdTable["info"] = infoAdvancedRootCmd.Run
		advancedCmdTable["uninstall"] = uninstallAdvancedRootCmd.Run
		advancedCmdTable[selectOpt(acceptedOpts, "")](cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(advancedRootCmd)
}
