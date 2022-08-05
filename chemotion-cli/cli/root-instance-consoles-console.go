package cli

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func dropIntoConsole(givenName string, consoleName string) {
	commandExec := exec.Command(toLower(virtualizer), []string{"compose", "exec", "eln", "chemotion", consoleName}...)
	commandExec.Stdin, commandExec.Stdout, commandExec.Stderr = os.Stdin, os.Stdout, os.Stderr
	if consoleName == "psql" {
		consoleName = "postgreSQL" // use proper name for psql when printing to user
	}
	if instanceStatus(givenName) == "Up" {
		zboth.Info().Msgf("Starting %s console for instance `%s`.", consoleName, givenName)
		if _, err, _ := gotoFolder(givenName), commandExec.Run(), gotoFolder("workdir"); err == nil {
			zboth.Debug().Msgf("Successfuly closed console for %s in `%s`.", consoleName, givenName)
		} else {
			zboth.Fatal().Err(err).Msgf("Console ended with exit message: %s.", err.Error())
		}
	} else {
		zboth.Warn().Err(toError("instance not running")).Msgf("Cannot start a %s console for `%s`. Instance is not running.", consoleName, givenName)
	}
}

var shellConsoleInstanceRootCmd = &cobra.Command{
	Use:     "shell",
	Aliases: []string{"bash"},
	Short:   "Drop into a shell (bash) console",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		dropIntoConsole(currentInstance, "shell")
	},
}

var railsConsoleInstanceRootCmd = &cobra.Command{
	Use:     "rails",
	Aliases: []string{"ruby", "railsc"},
	Short:   "Drop into a Ruby on Rails console",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		dropIntoConsole(currentInstance, "rails")
	},
}

var psqlConsoleInstanceRootCmd = &cobra.Command{
	Use:     "psql",
	Aliases: []string{"postgresql", "sql", "postgres", "postgreSQL"},
	Short:   "Drop into a PostgreSQL console",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		dropIntoConsole(currentInstance, "psql")
	},
}

func init() {
	consoleInstanceRootCmd.AddCommand(shellConsoleInstanceRootCmd)
	consoleInstanceRootCmd.AddCommand(railsConsoleInstanceRootCmd)
	consoleInstanceRootCmd.AddCommand(psqlConsoleInstanceRootCmd)
}
