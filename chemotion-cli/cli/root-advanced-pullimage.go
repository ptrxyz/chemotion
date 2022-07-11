package cli

import (
	"fmt"
	"net/url"

	"github.com/chigopher/pathlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _advanced_pull_use_ string

var pullimageAdvancedRootCmd = &cobra.Command{
	Use:   "pull-image",
	Args:  cobra.NoArgs,
	Short: fmt.Sprintf("pull latest image of %s from the internet", nameCLI),
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		var composeFilepath pathlib.Path
		var isURL bool
		if existingFile(_advanced_pull_use_) {
			composeFilepath = *pathlib.NewPath(_advanced_pull_use_)
		} else if _, err := url.ParseRequestURI(_advanced_pull_use_); err == nil {
			isURL = true
			composeFilepath = downloadFile(_advanced_pull_use_, workDir.String()) // downloads to the working directory
			if err := composeFilepath.RenameStr("temporary." + composeFilepath.Name()); err == nil {
				zboth.Info().Msgf("Renamed downloaded file to: %s", composeFilepath.Name())
			}
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to parse the URL/file: %s.", _advanced_pull_use_)
		}
		tempCompose := viper.New()
		tempCompose.SetConfigFile(composeFilepath.String())
		if err := tempCompose.ReadInConfig(); err != nil {
			zboth.Fatal().Err(err).Msgf("Failed to read the compose file.")
		}
		services, _ := getKeysValues(tempCompose, "services")
		for _, service := range services {
			zboth.Info().Msgf("Pulling image for the service called %s", service)
			if success := callVirtualizer(fmt.Sprintf("compose -f %s pull %s", composeFilepath.String(), service)); !success {
				zboth.Warn().Err(fmt.Errorf("pull failed")).Msgf("Failed to pull image for the service called %s", service)
			}
		}
		if isURL {
			if err := composeFilepath.Remove(); err != nil {
				zboth.Warn().Err(err).Msgf("Failed to delete the temporary file: %s.", composeFilepath.Name())
			}
		}
	},
}

func init() {
	pullimageAdvancedRootCmd.Flags().StringVar(&_advanced_pull_use_, "use", composeURL, "URL or filepath of the compose file to use when pulling the instance")
	advancedRootCmd.AddCommand(pullimageAdvancedRootCmd)
}
