package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Create a new instance of Chemotion
var createInstance = &cobra.Command{
	Use:        "create <name_of_instance>",
	SuggestFor: []string{"cre"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getArg(args, "Please enter name of the instance you want to create")
		fmt.Println("We are now going to create an instance called", actOn)
		//TODO
	},
}

// Determine status of the active/named instance of Chemotion
var statusInstance = &cobra.Command{
	Use:        "status <name_of_instance>",
	SuggestFor: []string{"stat"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getInstance(args)
		fmt.Println("This is what we know about the instance called", actOn)
		//TODO
	},
}

// Upgrade the active/named instance of Chemotion
var upgradeInstance = &cobra.Command{
	Use:        "upgrade <name_of_instance>",
	SuggestFor: []string{"upg"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getInstance(args)
		fmt.Println("We now upgrade the instance called", actOn)
		//TODO
	},
}

// Start an existing named of Chemotion
var startInstance = &cobra.Command{
	Use:        "start <name_of_instance>",
	SuggestFor: []string{"star"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getArg(args, "Name the instance you want to start.")
		fmt.Printf("Starting %s...\n", actOn)
		//TODO
	},
}

// Pause an existing instance of Chemotion
var pauseInstance = &cobra.Command{
	Use:        "pause <name_of_instance>",
	SuggestFor: []string{"pau"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getInstance(args)
		fmt.Printf("Pausing %s...\n", actOn)
		//TODO
	},
}

// Stop an existing instance of Chemotion
var stopInstance = &cobra.Command{
	Use:        "stop <name_of_instance>",
	SuggestFor: []string{"sto"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getInstance(args)
		fmt.Printf("Stopping %s...\n", actOn)
		//TODO
	},
}

// Restart an existing instance of Chemotion
var restartInstance = &cobra.Command{
	Use:        "restart <name_of_instance>",
	SuggestFor: []string{"res"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getInstance(args)
		fmt.Printf("Restarting %s...\n", actOn)
		//TODO
	},
}

// Delete an existing instance of Chemotion
var deleteInstance = &cobra.Command{
	Use:        "delete <name_of_instance>",
	SuggestFor: []string{"del"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getInstance(args)
		fmt.Printf("Deleting %s...\n", actOn)
		// Remember to change back to the fallback instance
		//TODO
	},
}

// Change currently selected instance
var switchInstance = &cobra.Command{
	Use:        "switch <name_of_instance>",
	SuggestFor: []string{"swi"},
	Args:       cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		actOn := getArg(args, "Name the instance you want to switch to")
		fmt.Printf("Switching to %s...\n", actOn)
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
	Run: func(cmd *cobra.Command, args []string) {
		confirmInteractive()
		fmt.Println("Chemotion. Actions on instance:")
		acceptedOpts := []string{"create", "status", "upgrade", "switch", "start", "pause", "stop", "restart", "delete"}
		switch selectOpt(acceptedOpts, args) {
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
