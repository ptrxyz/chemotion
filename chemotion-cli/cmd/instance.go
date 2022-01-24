package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Create a new instance of Chemotion
var createInstance = &cobra.Command{
	Use:        "create <name_of_instance>",
	SuggestFor: []string{"cre"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are now going to create an instance called", args[0])
		//TODO
	},
}

// Upgrade an existing instance of Chemotion
var updateInstance = &cobra.Command{
	Use:        "upgrade <name_of_instance>",
	SuggestFor: []string{"upg"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are now upgrade the instance called", args[0])
		//TODO
	},
}

var instanceCmd = &cobra.Command{
	Use:        "instance",
	Aliases:    []string{"i"},
	SuggestFor: []string{"ins"},
	Short:      "Manipulate instances of " + baseCommand,
	Long:       "Manipulate instances of " + baseCommand + "using one of the available actions",
	Args:       cobra.MinimumNArgs(1),
	// status|create|upgrade|start|pause|stop|restart|delete
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(createInstance)
	instanceCmd.AddCommand(updateInstance)
}
