package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/afero"
)

var (
	shell string = "bash"
)

// check if file exists
func existingFile(filePath string) (exists bool) {
	exists, _ = afero.Exists(fs, filePath) // no error if file exists, fs pointer is same as used by logger
	return
}

// check if the CLI is running interactively; if not, then exit.
func confirmInteractive() (silent bool) {
	if config.Quiet {
		zboth.Fatal().Msg(baseCommand + " is in quiet mode. Give all arguments to specify the desired action; use '--help' for more. ABORT!")
		silent = false
	} else {
		silent = true
	}
	return
}

// execute command in bash shell
func execShell(command string) (result string) {
	out, err := exec.Command(shell, "-c", command).Output()
	if err == nil {
		zlog.Debug().Str("container", containerName).Str("instance", config.Instance).Msg("Executed shell command: ")
		return string(out)
	} else {
		zboth.Fatal().Err(err).Msg("Check log. Failed to execute command: " + command)
		return ""
	}
}

// find version of a software using --version in bash shell
func findVersion(software string) (result string) {
	command := software + " --version"
	out, err := exec.Command(shell, "-c", command).Output()
	if err == nil {
		zlog.Debug().Str("container", containerName).Str("instance", config.Instance).Msg("Executed shell command: ")
		return fmt.Sprint(out)
	} else if err.Error() == "127" {
		return "Not installed!"
	} else {
		zboth.Fatal().Err(err).Msg("Check log. Failed to execute command: " + command)
		return ""
	}
}
