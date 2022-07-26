package cli

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// show (and then remove) a progress bar that waits for an instance to start
func waitStartSpinner(seconds int, givenName, message string) (waitTime int) {
	url := getURL(givenName)
	var (
		err  error
		code int
	)
	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription(fmt.Sprintf("%s %s...", message, givenName)),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionClearOnFinish(),
	)
	bar.RenderBlank()
	for i := 0; i < seconds; i++ {
		if code, err = instancePing(url); code == 200 {
			bar.Finish()
			waitTime = i
			return
		}
		bar.Add(1)
		time.Sleep(1 * time.Second)
	}
	bar.Finish()
	zboth.Warn().Err(err).Msgf("Response from the instance is %d", code)
	waitTime = -1
	return
}

func instanceStart(givenName string) {
	status := instanceStatus(givenName)
	if status == "Up" {
		zboth.Warn().Msgf("The instance called %s is already running.", givenName)
	} else {
		if _, success, _ := gotoFolder(givenName), callVirtualizer("compose up -d"), gotoFolder("workdir"); success {
			waitFor := 120 // in seconds
			if status == "Exited" {
				waitFor = 20
			}
			zboth.Info().Msgf("Starting instance called %s.", givenName)
			waitTime := waitStartSpinner(waitFor, givenName, "Starting")
			if waitTime >= 0 {
				zboth.Info().Msgf("Successfully started instance called %s in %d seconds.", givenName, waitTime)
			} else {
				zboth.Fatal().Msgf("Failed to start instance called %s. Please check logs using `%s instance %s`.", givenName, commandForCLI, logInstanceRootCmd.Use)
			}
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
