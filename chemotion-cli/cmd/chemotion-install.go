package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	_chemotion_install_name_ string
	_chemotion_install_use_  string
)

// command to install a new container of Chemotion
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the first instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		if !firstRun {
			zboth.Fatal().Err(fmt.Errorf("config file found")).Msgf("This option `install` is only available for initial installation. Use `chemotion instance new` to create more instances of %s.", nameCLI)
		}
		zlog.Debug().Msg("In installCmd.")
		create := true
		if currentState.Quiet {
			zboth.Info().Msgf("You chose do first run of chemotion in quiet mode. Will go ahead and install it!")
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
				zboth.Info().Msgf("All done! Now you can do `%s on` and `%s off` to start/stop %s.", cmd.Use, cmd.Use, nameCLI)
			}
		}
	},
}

func init() {
	installCmd.Flags().StringVar(&_chemotion_install_name_, "name", instanceDefault, "Name of the first instance to create")
	installCmd.Flags().StringVar(&_chemotion_install_use_, "use", composeURL, "URL or filepath to use for creating the instance")
}
