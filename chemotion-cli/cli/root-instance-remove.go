package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var _chemotion_instance_remove_name_ string
var _chemotion_instance_remove_force_ bool

func instanceRemove(given_name string) (success bool) {
	name := internalName(given_name)
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
	os.Chdir(workDir.Join(instancesFolder, name).String())
	confirmVirtualizer(minimumVirtualizer) // TODO if required: set virtualizer depending on compose file requirements
	if _chemotion_instance_remove_force_ {
		success = callVirtualizer("compose kill")
	} else {
		success = true
	}
	if success {
		success = callVirtualizer("compose down --volumes")
		if success {
			zboth.Info().Msgf("Successfully removed instance called %s.", given_name)
		}
	}
	// delete folder
	if success {
		pwd, _ := os.Getwd()
		zboth.Info().Msgf("Removing folder associated with %s. (arcane procedure!)", given_name)
		success = callVirtualizer("run --rm -v " + pwd + ":/x --name chemotion-helper-safe-to-remove busybox rm -rf x/shared")
	}
	os.Chdir("../..")
	if success {
		if err := workDir.Join(instancesFolder, name).RemoveAll(); err != nil {
			zboth.Warn().Err(err).Msgf("Failed to delete associated folder `%s` in `%s`.", name, instancesFolder)
		}
	}
	// delete entry in config
	if success {
		configMap := conf.GetStringMap("instances")
		delete(configMap, given_name)
		conf.Set("instances", configMap)
		if err := conf.WriteConfig(); err == nil {
			zboth.Info().Msgf("Modified configuration file `%s` to remove entry for `%s`.", conf.ConfigFileUsed(), given_name)
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to update the configuration file.")
		}
	}
	if !success {
		zboth.Info().Msgf("Clean deletion of %s failed. Check log to see what went wrong.", given_name)
	}
	return
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
