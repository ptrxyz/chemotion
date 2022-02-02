package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

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
