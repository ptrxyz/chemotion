package cli

import (
	"github.com/spf13/cobra"
)

func instanceStart(givenName string) {
	status := instanceStatus(givenName)
	if status == "Up" {
		zboth.Warn().Msgf("The instance called %s is already running.", givenName)
	} else {
		if _, success, _ := gotoFolder(givenName), callVirtualizer("compose up -d"), gotoFolder("workdir"); success {
			seconds := 40
			if status == "Exited" {
				seconds = 20
			}
			zboth.Info().Msgf("Starting instance called %s. Please give it %d seconds to initialize.", givenName, seconds)
			waitProgressBar(seconds, []string{"Starting", givenName})
			zboth.Info().Msgf("Successfully started instance called %s.", givenName)
		} else {
			zboth.Fatal().Msgf("Failed to start instance called %s.", givenName)
		}
	}
}

func instanceStop(givenName string) {
	status := instanceStatus(givenName)
	if status == "Up" {
		if _, success, _ := gotoFolder(givenName), callVirtualizer("compose stop"), gotoFolder("workdir"); success {
			zboth.Info().Msgf("Successfully stopped instance called %s.", givenName)
		} else {
			zboth.Fatal().Msgf("Failed to stop instance called %s.", givenName)
		}
	} else {
		zboth.Warn().Msgf("Cannot stop instance %s. It seems to be %s.", givenName, status)
	}
}

func instanceRestart(givenName string) {
	instanceStop(givenName)
	instanceStart(givenName)
}

var restartRootCmd = &cobra.Command{
	Use:   "restart [-i <instance_name>]",
	Args:  cobra.NoArgs,
	Short: "Restart the selected instance of " + nameCLI,
	Run: func(_ *cobra.Command, _ []string) {
		instanceRestart(currentInstance)
	},
	// TODO: add a force restart flag
}

var onRootCmd = &cobra.Command{
	Use:   "on [-i <instance_name>]",
	Args:  cobra.NoArgs,
	Short: "Start the selected instance of " + nameCLI,
	Run: func(_ *cobra.Command, _ []string) {
		instanceStart(currentInstance)
	},
}

var offRootCmd = &cobra.Command{
	Use:   "off [-i <instance_name>]",
	Args:  cobra.NoArgs,
	Short: "Stop the selected instance of " + nameCLI,
	Run: func(_ *cobra.Command, _ []string) {
		instanceStop(currentInstance)
	},
}

func init() {
	rootCmd.AddCommand(onRootCmd)
	rootCmd.AddCommand(offRootCmd)
	rootCmd.AddCommand(restartRootCmd)
}
