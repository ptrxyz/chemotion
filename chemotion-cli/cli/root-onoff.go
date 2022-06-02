package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func startInstance() {
	// TODO: check if it is already running
	confirmInstalled()
	os.Chdir(workDir.Join(instancesFolder, conf.GetString(joinKey("instances", currentState.name, "name"))).String())
	confirmVirtualizer(minimumVirtualizer)
	callVirtualizer("compose up -d")
	os.Chdir("../..")
}

func stopInstance() {
	// TODO: check if it is already running
	confirmInstalled()
	os.Chdir(workDir.Join(instancesFolder, conf.GetString(joinKey("instances", currentState.name, "name"))).String())
	confirmVirtualizer(minimumVirtualizer)
	callVirtualizer("compose down")
	os.Chdir("../..")
}

var onRootCmd = &cobra.Command{
	Use:   "on",
	Short: "start chemotion",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		startInstance()
	},
}

var offRootCmd = &cobra.Command{
	Use:   "off",
	Short: "stop chemotion",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		stopInstance()
	},
}

func init() {
	rootCmd.AddCommand(onRootCmd)
	rootCmd.AddCommand(offRootCmd)
}
