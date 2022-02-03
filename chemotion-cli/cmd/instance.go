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

// Determine status of an existing instance of Chemotion
var statusInstance = &cobra.Command{
	Use:        "status <name_of_instance>",
	SuggestFor: []string{"stat"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is what we know about the instance called", args[0])
		//TODO
	},
}

// Upgrade an existing instance of Chemotion
var upgradeInstance = &cobra.Command{
	Use:        "upgrade <name_of_instance>",
	SuggestFor: []string{"upg"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are now upgrade the instance called", args[0])
		//TODO
	},
}

// Start an existing instance of Chemotion
var startInstance = &cobra.Command{
	Use:        "start <name_of_instance>",
	SuggestFor: []string{"star"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting", args[0]+"...")
		//TODO
	},
}

// Pause an existing instance of Chemotion
var pauseInstance = &cobra.Command{
	Use:        "pause <name_of_instance>",
	SuggestFor: []string{"pau"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pausing", args[0]+"...")
		//TODO
	},
}

// Stop an existing instance of Chemotion
var stopInstance = &cobra.Command{
	Use:        "stop <name_of_instance>",
	SuggestFor: []string{"sto"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping", args[0]+"...")
		//TODO
	},
}

// Restart an existing instance of Chemotion
var restartInstance = &cobra.Command{
	Use:        "restart <name_of_instance>",
	SuggestFor: []string{"res"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Restarting", args[0]+"...")
		//TODO
	},
}

// Delete an existing instance of Chemotion
var deleteInstance = &cobra.Command{
	Use:        "delete <name_of_instance>",
	SuggestFor: []string{"del"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting", args[0]+"...")
		//TODO
	},
}

// Change currently selected instance
var switchInstance = &cobra.Command{
	Use:        "switch <name_of_instance>",
	SuggestFor: []string{"swi"},
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Switiching to", args[0]+"...")
		//TODO
	},
}

// Backbone for instance-related commands
var instanceCmd = &cobra.Command{
	Use:        "instance {create|status|upgrade|switch|start|pause|stop|restart|delete} <name_of_instance>",
	Aliases:    []string{"i"},
	SuggestFor: []string{"i"},
	Short:      "Manipulate instances of " + baseCommand,
	Long:       "Manipulate instances of " + baseCommand + " using one of the available actions",
	Args:       cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chemotion. Actions on instance:")
		switch prompter([]string{"create", "status", "upgrade", "switch", "start", "pause", "stop", "restart", "delete"}) {
		case "create":
			createInstance.Run(&cobra.Command{}, []string{})
		case "status":
			statusInstance.Run(&cobra.Command{}, []string{})
		case "upgrade":
			upgradeInstance.Run(&cobra.Command{}, []string{})
		case "switch":
			switchInstance.Run(&cobra.Command{}, []string{})
		case "start":
			startInstance.Run(&cobra.Command{}, []string{})
		case "pause":
			pauseInstance.Run(&cobra.Command{}, []string{})
		case "stop":
			stopInstance.Run(&cobra.Command{}, []string{})
		case "restart":
			restartInstance.Run(&cobra.Command{}, []string{})
		case "delete":
			deleteInstance.Run(&cobra.Command{}, []string{})
		}
	},
}

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(createInstance)
	instanceCmd.AddCommand(statusInstance)
	instanceCmd.AddCommand(upgradeInstance)
	instanceCmd.AddCommand(switchInstance)
	instanceCmd.AddCommand(startInstance)
	instanceCmd.AddCommand(pauseInstance)
	instanceCmd.AddCommand(stopInstance)
	instanceCmd.AddCommand(restartInstance)
	instanceCmd.AddCommand(deleteInstance)
}
