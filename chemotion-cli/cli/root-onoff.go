package cli

import (
	"github.com/spf13/cobra"
)

func instanceStart(given_name string) {
	name := internalName(given_name)
	if instanceStatus(given_name) == "Up" {
		zboth.Warn().Msgf("The instance called %s is already running.", given_name)
	} else {
		changeDir(workDir.Join(instancesFolder, name).String())
		confirmVirtualizer(minimumVirtualizer) // TODO if required: set virtualizer depending on compose file requirements
		callVirtualizer("compose up -d")
		zboth.Info().Msgf("Successfully started instance called %s. Please give it a minute to initialize.", given_name)
		changeDir("../..")
	}
}

func instanceStop(given_name string) {
	name := internalName(given_name)
	if instanceStatus(given_name) == "Up" {
		changeDir(workDir.Join(instancesFolder, name).String())
		confirmVirtualizer(minimumVirtualizer) // TODO if required: set virtualizer depending on compose file requirements
		callVirtualizer("compose stop")
		zboth.Info().Msgf("Successfully stopped instance called %s.", given_name)
		changeDir("../..")
	} else {
		zboth.Warn().Msgf("It seems that the instance %s is not running. Please check its status.", given_name)
	}
}

var onRootCmd = &cobra.Command{
	Use:   "on",
	Args:  cobra.NoArgs,
	Short: "Start (the selected instance of) chemotion",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		instanceStart(currentState.name)
	},
}

var offRootCmd = &cobra.Command{
	Use:   "off",
	Args:  cobra.NoArgs,
	Short: "Stop (the selected instance of) chemotion",
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
