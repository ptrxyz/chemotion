package cli

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func replceString(str string) (out string) {
	// make chemotion command names more readable
	switch str {
	case "shell":
		out = "shell console"
	case "railsc":
		out = "ruby console"
	case "psql":
		out = "postgreSQL console"
	}
	return out
}

func startInteractiveConsole(givenName string, prg string) {
	// Generic function for chemotion commands [i.e, shell, psql and railsc]

	arg1 := "docker"
	arg2 := "compose"
	arg3 := "exec"
	arg4 := "eln"
	arg5 := "chemotion"

	cmd := exec.Command(arg1, arg2, arg3, arg4, arg5, prg)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if instanceStatus(givenName) == "Up" {
		gotoFolder(givenName)
		zboth.Info().Msgf("[%s] is getting start for service '%s'.", replceString(prg), givenName)
		err := cmd.Run()

		if err != nil {
			zboth.Fatal().Msgf("Failed to start interactive shell for service %s.", givenName)
		}

		gotoFolder("workdir")

	} else {
		zboth.Warn().Msgf("Instance '%s' is either stopped or not yet started", givenName)
	}
}

var shellConsoleAdvancedRootCmd = &cobra.Command{
	Use:   "shell console",
	Short: "Allow user to interact with bash in ELN service " + nameCLI,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {

		startInteractiveConsole(currentInstance, "shell")
	},
}

var rubylConsoleAdvancedRootCmd = &cobra.Command{
	Use:   "ruby console",
	Short: "Allow user to interact with ruby console inside ELN service " + nameCLI,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {

		startInteractiveConsole(currentInstance, "railsc")
	},
}

var psqlConsoleAdvancedRootCmd = &cobra.Command{
	Use:   "PostgreSQL console",
	Short: "Allow user to interact with PostgreSQL console inside ELN service " + nameCLI,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {

		startInteractiveConsole(currentInstance, "psql")
	},
}

func init() {
	consoleAdvancedRootCmd.AddCommand(shellConsoleAdvancedRootCmd)
}
