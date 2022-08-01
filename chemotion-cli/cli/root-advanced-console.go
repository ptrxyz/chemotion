package cli

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var advancedCmdShellTable = make(cmdTable)

var consoleAdvancedRootCmd = &cobra.Command{
	Use:       "advanced",
	Short:     "Allow users to interact with a service via shell, rails console, PostgreSQL" + nameCLI,
	ValidArgs: maps.Keys(advancedCmdShellTable),
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Flag("selected-instance").Changed {
			zboth.Warn().Msgf("The `-i` flag is not supported for the `advanced` command and its subcommands.")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive(true)
		acceptedOpts := []string{"shell", "ruby", "psql"}
		if cmd.Use == cmd.CalledAs() { // || elementInSlice(cmd.CalledAs(), &cmd.Aliases) > -1 { { // there are no aliases at the moment
			acceptedOpts = append(acceptedOpts, "exit")
		} else {
			acceptedOpts = append(acceptedOpts, []string{"back", "exit"}...)
			advancedCmdShellTable["back"] = cmd.Run
		}
		advancedCmdShellTable["shell"] = shellConsoleAdvancedRootCmd.Run
		advancedCmdShellTable["ruby"] = rubylConsoleAdvancedRootCmd.Run
		advancedCmdShellTable["psql"] = psqlConsoleAdvancedRootCmd.Run
		advancedCmdShellTable[selectOpt(acceptedOpts, "")](cmd, args)
	},
}

func init() {
	advancedRootCmd.AddCommand(consoleAdvancedRootCmd)
}
