package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// confirm that virtualizer is the required minimum version
func confirmVirtualizer(minimum string) {
	if version := findVersion(toLower(virtualizer)); version == "docker on WSL not running!" {
		zboth.Fatal().Err(fmt.Errorf(version)).Msgf("Docker is not running in your WSL environment. Hint: Turn on WSL integration setting in Docker Desktop.")
	} else if version == "Unknown / not installed or found!" {
		zboth.Fatal().Err(fmt.Errorf(version)).Msgf("%s is necessary to run %s", virtualizer, nameCLI)
	} else {
		zlog.Debug().Msgf("%s is installed", virtualizer)
	}
	if err := compareSoftwareVersion(minimum, toLower(virtualizer)); err != nil {
		zboth.Fatal().Err(err).Msgf(err.Error())
	}
	zlog.Debug().Msgf("%s version requirement met", virtualizer)
}

// call to virtualizer
func callVirtualizer(args string) (success bool) {
	zboth.Info().Msgf("%s will now fork the execution with command `%s %s` sent to shell. Will return once execution is completed.", nameCLI, toLower(virtualizer), args)
	commandArgs := strings.Split(args, " ")
	commandExec := exec.Command(toLower(virtualizer), commandArgs...)
	// see https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html#:~:text=Capture%20output%20but%20also%20show%20progress%20%233
	var stdoutBuf, stderrBuf bytes.Buffer
	commandExec.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	commandExec.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	if err := commandExec.Run(); err == nil {
		success = true
	} else {
		success = false
		zboth.Warn().Err(err).Msgf("%s command failed! Check log. ABORT!", toLower(virtualizer))
	}
	// TODO-v2: make this more elegent by using a virtual terminal, see https://github.com/creack/pty
	return
}
