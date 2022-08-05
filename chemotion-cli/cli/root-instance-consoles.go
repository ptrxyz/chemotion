package cli

import (
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var consoleInstanceCmdTable = make(cmdTable)

var consoleInstanceRootCmd = &cobra.Command{
	Use:       "console",
	Short:     "Allow users to interact with an instance's command line interface",
	ValidArgs: maps.Keys(consoleInstanceCmdTable),
	Run: func(cmd *cobra.Command, args []string) {
		isInteractive(true)
		acceptedOpts := []string{"shell", "ruby", "psql"}
		consoleInstanceCmdTable["shell"] = shellConsoleInstanceRootCmd.Run
		consoleInstanceCmdTable["ruby"] = railsConsoleInstanceRootCmd.Run
		consoleInstanceCmdTable["psql"] = psqlConsoleInstanceRootCmd.Run
		if cmd.Use == cmd.CalledAs() { // || elementInSlice(cmd.CalledAs(), &cmd.Aliases) > -1 { { // there are no aliases at the moment
			acceptedOpts = append(acceptedOpts, "exit")
		} else {
			acceptedOpts = append(acceptedOpts, []string{"back", "exit"}...)
			consoleInstanceCmdTable["back"] = cmd.Run
		}
		consoleInstanceCmdTable[selectOpt(acceptedOpts, "")](cmd, args)
	},
}

func init() {
	instanceRootCmd.AddCommand(consoleInstanceRootCmd)
}
