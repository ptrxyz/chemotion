package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

// err check function
func check_err(err error, message string) {
	if err != nil {
		if len(message) != 0 {
			message = "An error occured!"
		}
		fmt.Println(message)
		fmt.Println("Error was:", err)
		os.Exit(1)
	}
}

// check if the CLI is running quietly
func is_quiet(called_as string) bool {
	if quiet {
		fmt.Println(baseCommand, "is running in non-interactive mode.")
		fmt.Printf("Called command: `%s` is incomplete. Use the `--help` flag after it to learn more.", called_as)
		os.Exit(1)
		return false
	} else {
		return true
	}
}

// selection function
func selectOpt(items []string) (result string) {
	selection := promptui.Select{
		Label: fmt.Sprintf("[%s] Select one of the following:", active_instance),
		Items: items,
	}
	_, result, err := selection.Run()
	check_err(err, "Selection failed!")
	return
}

// string prompt function
func promptStr(label string) (result string) {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	check_err(err, "Prompt failed!")
	return
}

// process args
func getArgs(args []string, prompt string) (label string) {
	if len(args) == 0 {
		is_quiet(rootCmd.CalledAs())
		label = promptStr(prompt)
	} else {
		label = args[0]
	}
	return
}

var failedMessage = "Failed to execute command: "

// execute command in bash shell
func execShell(command string) (result string) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return fmt.Sprintf(failedMessage+"%s with message %s", command, err.Error())
	}
	return string(out)
}

// find version of a software using --version in bash shell
func findVersion(software string) (result string) {
	result = execShell(software + " --version")
	if strings.HasPrefix(result, failedMessage) && strings.HasSuffix(result, "127") {
		return "Not installed!"
	} else {
		return result
	}
}
