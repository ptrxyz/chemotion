package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var _root_instance_remove_name_ string
var _root_instance_remove_force_ bool

func instanceRemove(givenName string) (success bool) {
	name := internalName(givenName)
	if !_root_instance_remove_force_ {
		if instanceStatus(givenName) == "Up" {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete an instance that is currently running. Please use `chemotion -i %s stop` to stop the instance.")
		}
		if givenName == currentState.name {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete the currently selected instance. Use `chemotion switch` to switch selection to another instance before proceeding.")
		}
		if len(allInstances()) == 1 {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete the only instance. Use `chemotion advanced uninstall` remove %s entirely", nameCLI)
		}
	}
	changeDir(workDir.Join(instancesFolder, name).String())
	if _root_instance_remove_force_ {
		success = callVirtualizer("compose kill")
	} else {
		success = true
	}
	if success {
		success = callVirtualizer("compose down --remove-orphans --volumes")
		if success {
			zboth.Info().Msgf("Successfully removed instance called %s.", givenName)
		}
	}
	// delete folder
	if success {
		pwd, _ := os.Getwd()
		zboth.Info().Msgf("Removing folder associated with %s. (arcane procedure!)", givenName)
		success = callVirtualizer("run --rm -v " + pwd + ":/x --name chemotion-helper-safe-to-remove busybox rm -rf x/shared")
	}
	changeDir("../..")
	if success {
		if err := workDir.Join(instancesFolder, name).RemoveAll(); err != nil {
			zboth.Warn().Err(err).Msgf("Failed to delete associated folder `%s` in `%s`.", name, instancesFolder)
		}
	}
	// delete entry in config
	if success {
		configMap := conf.GetStringMap("instances")
		delete(configMap, givenName)
		conf.Set("instances", configMap)
		if err := conf.WriteConfig(); err == nil {
			zboth.Info().Msgf("Modified configuration file `%s` to remove entry for `%s`.", conf.ConfigFileUsed(), givenName)
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to update the configuration file.")
		}
	}
	if !success {
		zboth.Info().Msgf("Clean deletion of %s failed. Check log to see what went wrong.", givenName)
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
		if _root_instance_remove_name_ == "" {
			confirmInteractive()
			_root_instance_remove_name_ = selectInstance("remove")
		}
		instanceRemove(_root_instance_remove_name_)
	},
}

func init() {
	instanceRootCmd.AddCommand(removeInstanceRootCmd)
	removeInstanceRootCmd.Flags().StringVar(&_root_instance_remove_name_, "name", "", "name of the instance to remove")
	removeInstanceRootCmd.Flags().BoolVar(&_root_instance_remove_force_, "force", false, "force remove an instance (very risky)")
}
