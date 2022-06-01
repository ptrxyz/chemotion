package cli

import (
	"github.com/spf13/cobra"
)

var instanceRootCmd = &cobra.Command{
	Use:   "instance {create|status|upgrade|switch|start|pause|stop|restart|delete} <name_of_instance>",
	Short: "Manipulate instances of " + nameCLI,
	Long:  "Manipulate instances of " + nameCLI + " using one of the available actions",
	Run: func(cmd *cobra.Command, args []string) {
		logCall(cmd.Use, cmd.CalledAs())
		confirmInstalled()
		confirmInteractive()
		acceptedOpts := []string{"new", "exit"} //, "status", "upgrade", "switch", "start", "pause", "stop", "restart", "delete"}
		switch selectOpt(acceptedOpts) {
		case "new":
			newInstanceRootCmd.Run(&cobra.Command{}, []string{})
			// case "status":
			// 	statusInstance.Run(&cobra.Command{}, []string{})
			// case "upgrade":
			// 	upgradeInstance.Run(&cobra.Command{}, []string{})
			// case "switch":
			// 	switchInstance.Run(&cobra.Command{}, []string{})
			// case "start":
			// 	startInstance.Run(&cobra.Command{}, []string{})
			// case "pause":
			// 	pauseInstance.Run(&cobra.Command{}, []string{})
			// case "stop":
			// 	stopInstance.Run(&cobra.Command{}, []string{})
			// case "restart":
			// 	restartInstance.Run(&cobra.Command{}, []string{})
			// case "delete":
			// 	deleteInstance.Run(&cobra.Command{}, []string{})
		case "exit":
			zlog.Debug().Msg("Chose to exit.")
		}
	},
}

func init() {
	rootCmd.AddCommand(instanceRootCmd)
}
