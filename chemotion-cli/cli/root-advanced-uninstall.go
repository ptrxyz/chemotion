package cli

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var uninstallAdvancedRootCmd = &cobra.Command{
	Use:   "uninstall (accepts no flags)",
	Args:  cobra.NoArgs,
	Short: fmt.Sprintf("uninstall %s completely", nameCLI),
	Run: func(cmd *cobra.Command, args []string) {
		if currentState.quiet {
			zboth.Info().Msgf("For security reasons, this command will not run in silent mode.")
		}
		zerolog.SetGlobalLevel(zerolog.DebugLevel) // uninstall operates in debug mode
		fmt.Println("Uninstall operates in debug mode!")
		logWhere()
		confirmInstalled()
		confirmInteractive()
		removelog := selectYesNo("Do you want to remove the log file as well", false)
		if selectYesNo("Are you sure you want to uninstall "+nameCLI, false) {
			chosen := conf.GetString(selector_key)
			instances := append(allInstances(), chosen)
			skip := true
			for _, inst := range instances {
				if inst == chosen && skip { // contraption to make sure that the chosen instance is deleted last
					skip = false
					continue
				}
				zboth.Info().Msgf("Removing instance called %s.", inst)
				_chemotion_instance_remove_force_ = true
				if !instanceRemove(inst) {
					zboth.Fatal().Err(fmt.Errorf("uninstalled failed")).Msgf("Uninstall failed while trying to remove %s", inst)
					break
				}
				workDir.Join(instancesFolder).RemoveAll()
				workDir.Join(conf.ConfigFileUsed()).Remove()
			}
			if removelog {
				workDir.Join(logFilename).Remove()
			}
			zboth.Info().Msgf("%s successfully removed.", nameCLI)
		} else {
			zboth.Info().Msgf("Nothing was done.")
		}
	},
}

func init() {
	advancedRootCmd.AddCommand(uninstallAdvancedRootCmd)
}
