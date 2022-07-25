package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func advancedPullImage(use string) {
	tempCompose := getCompose(use)
	services := getSubHeadings(&tempCompose, "services")
	if len(services) == 0 {
		zboth.Warn().Err(fmt.Errorf("no services found")).Msgf("Please check that %s is a valid compose file with named services.", tempCompose.ConfigFileUsed())
	}
	for _, service := range services {
		zboth.Info().Msgf("Pulling image for the service called %s", service)
		if success := callVirtualizer(fmt.Sprintf("pull %s", tempCompose.GetString(joinKey("services", service, "image")))); !success {
			zboth.Warn().Err(fmt.Errorf("pull failed")).Msgf("Failed to pull image for the service called %s", service)
		}
	}
}

var pullimageAdvancedRootCmd = &cobra.Command{
	Use:   "pull-image",
	Args:  cobra.NoArgs,
	Short: fmt.Sprintf("pull latest image of %s from the internet", nameCLI),
	Run: func(cmd *cobra.Command, _ []string) {
		var use string
		if ownCall(cmd) {
			use = cmd.Flag("use").Value.String()
		} else {
			use = composeURL
		}
		zboth.Info().Msgf("Using compose file from: %s", use)
		advancedPullImage(use)
	},
}

func init() {
	pullimageAdvancedRootCmd.Flags().String("use", composeURL, "URL or filepath of the compose file to use when pulling the instance")
	advancedRootCmd.AddCommand(pullimageAdvancedRootCmd)
}
