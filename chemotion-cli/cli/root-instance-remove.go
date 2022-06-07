package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var _chemotion_instance_remove_name_ string

func removeInstance(name string) {
	if name == currentState.name {
		allInstances, _ := getKeysValues(&conf, "instances")
		if len(allInstances) > 1 {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete the currently selected instance. Switch selected instance before proceeding.")
		} else {
			zboth.Fatal().Err(fmt.Errorf("illegal operation")).Msgf("Cannot delete the only instance. Use `chemotion advanced uninstall` remove %s entirely", nameCLI)
		}
	}
	// TODO: actually remove the instance
}

var removeInstanceRootCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an existing instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		if _chemotion_instance_remove_name_ == "" {
			confirmInteractive()
			_chemotion_instance_remove_name_ = selectInstance("remove")
		}
		removeInstance(_chemotion_instance_remove_name_)
	},
}

func init() {
	instanceRootCmd.AddCommand(removeInstanceRootCmd)
	removeInstanceRootCmd.Flags().StringVar(&_chemotion_instance_remove_name_, "name", "", "name of the instance to remove")
}
