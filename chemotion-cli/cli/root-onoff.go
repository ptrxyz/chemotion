package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func instanceStart(given_name string) {
	name := internalName(given_name)
	if instanceStatus(given_name) == "Up" {
		zboth.Warn().Msgf("The instance called %s is already running.", given_name)
	} else {
		os.Chdir(workDir.Join(instancesFolder, name).String())
		confirmVirtualizer(minimumVirtualizer) // TODO if required: set virtualizer depending on compose file requirements
		callVirtualizer("compose up -d")
		zboth.Info().Msgf("Successfully started instance called %s.", given_name)
		os.Chdir("../..")
	}
}

func instanceStop(given_name string) {
	name := internalName(given_name)
	if instanceStatus(given_name) == "Up" {
		os.Chdir(workDir.Join(instancesFolder, name).String())
		confirmVirtualizer(minimumVirtualizer) // TODO if required: set virtualizer depending on compose file requirements
		callVirtualizer("compose stop")
		zboth.Info().Msgf("Successfully stopped instance called %s.", given_name)
		os.Chdir("../..")
	} else {
		zboth.Warn().Msgf("It seems that the instance %s is not running. Please check its status.", given_name)
	}
}

var onRootCmd = &cobra.Command{
	Use:   "on",
	Short: "start chemotion",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		instanceStart(currentState.name)
	},
}

var offRootCmd = &cobra.Command{
	Use:   "off",
	Short: "stop chemotion",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		instanceStop(currentState.name)
	},
}

func init() {
	rootCmd.AddCommand(onRootCmd)
	rootCmd.AddCommand(offRootCmd)
}
