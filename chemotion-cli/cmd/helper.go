package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/afero"
)

const (
	filePerm   = os.FileMode(0644)
	folderPerm = os.FileMode(0755)
)

var (
	shell         string = "bash"
	notInstallErr string = "Not installed!"
)

// check if file exists
func existingFile(filePath string) (exists bool) {
	exists, _ = afero.Exists(fs, filePath) // no error if file exists, fs pointer is same as used by logger
	return
}

// check if the CLI is running interactively; if not, then exit.
func confirmInteractive() (silent bool) {
	if currentState.Quiet {
		zboth.Fatal().Msg(projectName + " is in quiet mode. Give all arguments to specify the desired action; use '--help' for more. ABORT!")
		silent = false
	} else {
		silent = true
	}
	return
}

// check if a string is an array of strings
func stringInArray(str string, strings []string) int {
	for index, element := range strings {
		if element == str {
			return index
		}
	}
	return -1
}

// execute command in bash shell
func execShell(command string) (result string, err error) {
	res, err := exec.Command(shell, "-c", command).Output()
	result = strings.TrimSpace(string(res))
	if err == nil {
		zlog.Debug().Str("instance", currentState.name).Msg("Executed shell command: " + command)
	}
	return
}

// find version of a software using --version in bash shell
func findVersion(software string) (version string) {
	command := software + " --version"
	version, err := execShell(command)
	if err == nil {
	} else if err.Error() == "exit status 127" { // 127 is generally command not found
		version = notInstallErr
	} else {
		zboth.Warn().Err(err).Msg("Check log. Failed to execute command: " + command)
	}
	return
}
