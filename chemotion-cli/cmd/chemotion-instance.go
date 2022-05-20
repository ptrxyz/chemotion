package cmd

import (
	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:   "instance {create|status|upgrade|switch|start|pause|stop|restart|delete} <name_of_instance>",
	Short: "Manipulate instances of " + nameCLI,
	Long:  "Manipulate instances of " + nameCLI + " using one of the available actions",
	Run: func(cmd *cobra.Command, args []string) {
		confirmInteractive()
		acceptedOpts := []string{"new"} //, "status", "upgrade", "switch", "start", "pause", "stop", "restart", "delete"}
		switch selectOpt(acceptedOpts) {
		case "new":
			newInstanceCmd.Run(&cobra.Command{}, []string{})
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
		}
	},
}

func init() {
	instanceCmd.AddCommand(newInstanceCmd)
}
