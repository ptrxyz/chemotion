package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Show system related information to the user
var infoSystem = &cobra.Command{
	Use:        "info",
	SuggestFor: []string{"inf"},
	Args:       cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ncpus := runtime.NumCPU()
		mem := strings.Fields(execShell("free -h"))
		rubyVersion := findVersion("ruby")
		passengerVersion := findVersion("passenger")
		nodeVersion := findVersion("node")
		npmVersion := findVersion("npm")

		fmt.Println("This is what we know about the system")
		fmt.Println("- CPU Cores:", ncpus)
		fmt.Println("- Memory:\n  -", mem[7], "(total),", mem[9], "(free)")
		fmt.Println("Used software versions:")
		fmt.Println("- Ruby:", rubyVersion)
		fmt.Println("- Passenger:", passengerVersion)
		fmt.Println("- Node:", nodeVersion)
		fmt.Println("- npm:", npmVersion)

	},
}

// Start shell for user
var shellSystem = &cobra.Command{
	Use:        "shell",
	SuggestFor: []string{"she"},
	Args:       cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are now going to start shell")
		//TODO
	},
}

// Start a rails shell for user
var railsSystem = &cobra.Command{
	Use:        "rails",
	SuggestFor: []string{"rai"},
	Args:       cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("We are now going to start Rails shell")
		//TODO
	},
}

// Backbone for system-related commands
var systemCmd = &cobra.Command{
	Use:        "system {info|shell|rails}",
	Aliases:    []string{"s"},
	SuggestFor: []string{"s"},
	Short:      "Perform system-oriented actions",
	Long:       "Perform system-oriented actions using one of the available actions",
	Run: func(cmd *cobra.Command, args []string) {
		confirmInteractive()
		fmt.Println("Chemotion. Available system resources.")
		// acceptedOpts := []string{"info", "shell", "rails"}
		switch "info" {
		case "info":
			infoSystem.Run(&cobra.Command{}, []string{})
		case "shell":
			shellSystem.Run(&cobra.Command{}, []string{})
		case "rails":
			railsSystem.Run(&cobra.Command{}, []string{})
		}
	},
}

func init() {
	chemotionCmd.AddCommand(systemCmd)
	systemCmd.AddCommand(infoSystem)
	systemCmd.AddCommand(shellSystem)
	systemCmd.AddCommand(railsSystem)
}
