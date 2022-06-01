package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	_chemotion_install_name_ string
	_chemotion_install_use_  string
)

// command to install a new container of Chemotion
var installRootCmd = &cobra.Command{
	Use:   "install",
	Short: "Initialize the configuration file and install the first instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logCall(cmd.Use, cmd.CalledAs())
		if !firstRun {
			zboth.Fatal().Err(fmt.Errorf("config file found")).Msgf("This option `%s` is only available for initial installation. Use `%s %s %s` if you wish to create more instances of %s.", cmd.Name(), rootCmd.Name(), instanceRootCmd.Name(), newInstanceRootCmd.Name(), nameCLI)
		}
		create := true
		if currentState.quiet {
			zboth.Info().Msgf("You chose do first run of %s in quiet mode. Will go ahead and install it!", nameCLI)
		} else {
			if selectYesNo("Installation process may download containers (of multiple GBs) and can take some time. Continue", false) {
				if _chemotion_install_name_ == instanceDefault {
					_chemotion_install_name_ = getString("Please enter the name of the first instance you want to create")
				}
			} else {
				create = false
				zboth.Info().Msgf("Installation cancelled.")
			}
		}
		if create {
			zboth.Info().Msgf("We are now going to create an instance called %s.", _chemotion_install_name_)
			if success := instanceCreate(_chemotion_install_name_, "Production", _chemotion_install_use_); success {
				zboth.Info().Msgf("All done! Now you can do `%s on` and `%s off` to start/stop %s.", rootCmd.Name(), rootCmd.Name(), nameCLI)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installRootCmd)
	installRootCmd.Flags().StringVar(&_chemotion_install_name_, "name", instanceDefault, "Name of the first instance to create")
	installRootCmd.Flags().StringVar(&_chemotion_install_use_, "use", composeURL, "URL or filepath to use for creating the instance")
}
