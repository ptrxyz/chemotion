package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

// call to docker
func callDocker(args string) (success bool) {
	zboth.Info().Msgf("%s will now fork the execution of the command to %s and then pick it up once it is completed.", projectName, shell)
	commandArgs := strings.Split(args, " ")
	commandExec := exec.Command("docker", commandArgs...)
	// out, err := commandExec.CombinedOutput()
	// fmt.Println(string(out))
	var stdoutBuf, stderrBuf bytes.Buffer
	commandExec.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	commandExec.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := commandExec.Run()
	if err == nil {
		success = true
	} else {
		zboth.Fatal().Err(err).Msg("Docker command failed! Check log. ABORT!")
	}
	return
}

// confirm that Docker is installed
func confirmDocker() (installed bool) {
	if _, err := execShell("docker-compose --version"); err == nil {
		installed = true
	} else {
		zboth.Fatal().Err(err).Msg("Docker compose is not installed / accessible in path.")
		installed = false
	}
	// TODO check on the version number as well
	// docker engine 17.12.0 && compose 3.5
	return
}
