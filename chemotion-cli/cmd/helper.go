package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/chigopher/pathlib"
	"github.com/google/uuid"
	vercompare "github.com/hashicorp/go-version"
)

var versionSuffix string = " --version"

// check if file exists, and is a file (keep it simple, runs before logging starts!
func existingFile(filePath string) (exists bool) {
	p := pathlib.NewPath(filePath)
	exists, _ = p.IsFile()
	return
}

// check if the CLI is running interactively; if not, then exit. Wrapper around currentstate.Quiet.
func confirmInteractive() {
	if currentState.Quiet {
		fmt.Println(nameCLI + " is in quiet mode. Give all arguments to specify the desired action; use '--help' for more. ABORT!")
		zboth.Fatal().Err(fmt.Errorf("incomplete in quiet mode")).Msgf("%s is in quiet mode. Give all arguments to specify the desired action; use '--help' for more. ABORT!", nameCLI)
	}
}

// check if a string is an array of strings, if yes, return the 1st index, else -1.
func stringInArray(str string, strings []string) int {
	for index, element := range strings {
		if element == str {
			return index
		}
	}
	return -1
}

// execute a command in shell
func execShell(command string) (result []byte, err error) {
	if result, err = exec.Command(shell, "-c", command).Output(); err == nil {
		zlog.Debug().Msgf("Sucessfully executed shell command: %s", command)
	} else if !strings.HasSuffix(command, versionSuffix) {
		zboth.Warn().Err(err).Msgf("Failed execution of command: %s", command)
	}
	return
}

func findVersion(software string) (version string) {
	ver, err := execShell(software + versionSuffix)
	version = strings.Split(strings.TrimPrefix(strings.TrimPrefix(string(ver), "v"), "Docker version "), ",")[0] // TODO: Regexify!
	if err != nil {
		zlog.Debug().Err(err).Msgf("Version determination of %s failed.", software)
		if virtualizer == "docker" && err.Error() == "exit status 1" {
			version = "docker on WSL not running!"
		} else {
			version = "Unknown / not installed or found!"
		}
	}
	return
}

// confirm a minimum version for a given software
func compareSoftwareVersion(min string, software string) error {
	var minimum, current *vercompare.Version
	if ver, err := vercompare.NewVersion(findVersion(software)); err == nil {
		current = ver
	} else {
		return err
	}
	if ver, err := vercompare.NewVersion(min); err == nil {
		minimum = ver
	} else {
		return err
	}
	if current.LessThan(minimum) {
		return fmt.Errorf("current version of %s: %s is less than the minimum required: %s", software, current.String(), minimum.String())
	}
	return nil
}

// generate a new UID (of the form xxxxxxxx) as a string
func getNewUniqueID() string {
	id, _ := uuid.NewRandom()
	return strings.Split(id.String(), "-")[0]
}

// download a file
func downloadFile(fileURL string, downloadLocation string) (filename string) {
	if resp, err := grab.Get(downloadLocation, fileURL); err == nil {
		zboth.Info().Msgf("Downloaded file saved as: %s", resp.Filename)
		filename = resp.Filename
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to download file from: %s. Check log. ABORT!", fileURL)
	}
	return
}

// copy a text file
func copyTextFile(source *pathlib.Path, target *pathlib.Path) (err error) {
	if reader, errRead := source.ReadFile(); err == nil {
		if errWrite := target.WriteFile(reader); err != nil {
			err = errWrite
		}
	} else {
		err = errRead
	}
	return
}
