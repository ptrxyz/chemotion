package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var _chemotion_instance_remove_name_ string
var _chemotion_instance_remove_force_ bool

func instanceRemove(given_name string) {
	if !_chemotion_instance_remove_force_ {
		if instanceStatus(given_name) == "Up" {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete an instance that is currently running. Please use `chemotion -i %s stop` to stop the instance.")
		}
		if given_name == currentState.name {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete the currently selected instance. Use `chemotion switch` to switch selection to another instance before proceeding.")
		}
		if len(allInstances()) == 1 {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete the only instance. Use `chemotion advanced uninstall` remove %s entirely", nameCLI)
		}
	}
	name := internalName(given_name)
	os.Chdir(workDir.Join(instancesFolder, name).String())
	confirmVirtualizer(minimumVirtualizer) // TODO if required: set virtualizer depending on compose file requirements
	if _chemotion_instance_remove_force_ {
		callVirtualizer("compose kill")
	}
	success := callVirtualizer("compose down --volumes")
	zboth.Info().Msgf("Successfully removed instance called %s.", given_name)
	os.Chdir("../..")
	if err := workDir.Join(instancesFolder, name).RemoveAll(); err != nil { // doesn't work because of permission issues!
		zboth.Warn().Err(err).Msgf("Failed to delete associated folder %s in %s.", name, instancesFolder)
	}
	if success {
		configMap := conf.GetStringMap("instances")
		delete(configMap, given_name)
		conf.Set("instances", configMap)
		if err := conf.WriteConfig(); err == nil {
			zboth.Info().Msgf("Modified configuration file %s to remove entry for %s.", conf.ConfigFileUsed(), given_name)
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to update the configuration file.")
		}
	}
}

var removeInstanceRootCmd = &cobra.Command{
	Use:   "remove",
	Args:  cobra.NoArgs,
	Short: "Remove an existing instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		if _chemotion_instance_remove_name_ == "" {
			confirmInteractive()
			_chemotion_instance_remove_name_ = selectInstance("remove")
		}
		instanceRemove(_chemotion_instance_remove_name_)
	},
}

func init() {
	instanceRootCmd.AddCommand(removeInstanceRootCmd)
	removeInstanceRootCmd.Flags().StringVar(&_chemotion_instance_remove_name_, "name", "", "name of the instance to remove")
	removeInstanceRootCmd.Flags().BoolVar(&_chemotion_instance_remove_force_, "force", false, "force remove an instance (very risky)")
}
