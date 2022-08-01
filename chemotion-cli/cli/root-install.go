package cli

import (
	"github.com/spf13/cobra"
)

// command to install a new container of Chemotion
var installRootCmd = &cobra.Command{
	Use:   "install",
	Args:  cobra.NoArgs,
	Short: "Initialize the configuration file and install the first instance of " + nameCLI,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if cmd.Flag("selected-instance").Changed {
			zboth.Warn().Msgf("The `-i` flag is not supported for the `install` command.")
		}
	},
	Run: func(cmd *cobra.Command, _ []string) {
		var create bool
		if firstRun {
			if isInteractive(false) {
				create = newInstanceInteraction(cmd)
			} else {
				zboth.Info().Msgf("You chose do first run of %s in quiet mode. Will go ahead and install it!", nameCLI)
				create = true
			}
			if create {
				if success := instanceCreateProduction(cmd); success {
					zboth.Info().Msgf("All done! Now you can do `%s on` and `%s off` to start/stop %s.", commandForCLI, commandForCLI, nameCLI)
				}
			}
		} else {
			zboth.Fatal().Err(toError("config file found")).Msgf("This option `%s` is only available for initial installation. Use `%s %s %s` if you wish to create more instances of %s.", cmd.Name(), rootCmd.Name(), instanceRootCmd.Name(), newInstanceRootCmd.Name(), nameCLI)
		}
	},
}

func init() {
	rootCmd.AddCommand(installRootCmd)
	installRootCmd.Flags().String("name", instanceDefault, "Name of the first instance to create")
	installRootCmd.Flags().String("use", composeURL, "URL or filepath of the compose file to use for creating the instance")
	installRootCmd.Flags().String("address", addressDefault, "Web-address (or hostname) for accessing the instance")
	installRootCmd.Flags().String("env", "", ".env file for the first instance")
	installRootCmd.Flags().Uint("expose", 0, "Port that is exposed by the instance to access it")
}
