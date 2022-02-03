package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

// prompting function
func prompter(items []string) (result string) {
	active_instance := "default"
	prompt := promptui.Select{
		Label: fmt.Sprintf("[%s] ", active_instance) + "Select one of the following:",
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		panic("Prompt failed!")
	}
	return
}

var failedMessage = "Failed to execute command: "

func execShell(command string) (result string) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return fmt.Sprintf(failedMessage+"%s with message %s", command, err.Error())
	}
	return string(out)
}

func findVersion(software string) (result string) {
	result = execShell(software + " --version")
	if strings.HasPrefix(result, failedMessage) && strings.HasSuffix(result, "127") {
		return "Not installed!"
	} else {
		return result
	}
}
