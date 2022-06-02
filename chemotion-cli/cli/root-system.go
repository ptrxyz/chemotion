package cli

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// helper function that is also used by infoInstanceRootCmd
func systemInfo() {
	// CPU
	fmt.Println("- CPU Cores:", runtime.NumCPU())
	// Memory
	if mem, err := execShell("free -h"); err == nil {
		mem := strings.Fields(string(mem))
		fmt.Println("- Memory:\n  -", mem[7], "(total),", mem[9], "(free)")
	} else {
		fmt.Println("Couldn't determine memory usage.")
	}
	fmt.Println("Used software versions:")
	printVersionOf := []string{"ruby", "passenger", "node", "npm"}
	for _, software := range printVersionOf {
		fmt.Printf("- %s: %s\n", strings.ToTitle(software), findVersion(software))
	}
}

// Show host machine information to the user
// See also, chemotion instance info
var infoSystemRootCmd = &cobra.Command{
	Use:                   "info",
	Args:                  cobra.MaximumNArgs(0),
	Short:                 "get information about the system",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		fmt.Println("This is what we know about the host machine:")
		systemInfo()
	},
}

// Start shell for user
var shellSystemRootCmd = &cobra.Command{
	Use:        "shell",
	SuggestFor: []string{"she"},
	Args:       cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		fmt.Println("We are now going to start shell")
		//TODO
	},
}

// Start a rails shell for user
var railsSystemRootCmd = &cobra.Command{
	Use:        "rails",
	SuggestFor: []string{"rai"},
	Args:       cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		fmt.Println("We are now going to start Rails shell")
		//TODO
	},
}

// Backbone for system-related commands
var systemRootCmd = &cobra.Command{
	Use:        "system",
	Aliases:    []string{"s"},
	SuggestFor: []string{"s"},
	Short:      "Perform system-oriented actions",
	Long:       "Perform system-oriented actions using one of the available actions",
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		confirmInteractive()
		fmt.Println("Chemotion. Available system resources.")
		acceptedOpts := []string{"info", "shell", "rails", "exit"}
		selected := selectOpt(acceptedOpts)
		switch selected {
		case "info":
			infoSystemRootCmd.Run(&cobra.Command{}, []string{})
		case "shell":
			shellSystemRootCmd.Run(&cobra.Command{}, []string{})
		case "rails":
			railsSystemRootCmd.Run(&cobra.Command{}, []string{})
		case "exit":
			zlog.Debug().Msg("Chose to exit")
		}
	},
}

func init() {
	rootCmd.AddCommand(systemRootCmd)
	systemRootCmd.AddCommand(infoSystemRootCmd)
	systemRootCmd.AddCommand(shellSystemRootCmd)
	systemRootCmd.AddCommand(railsSystemRootCmd)
}
