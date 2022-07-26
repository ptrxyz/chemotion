package cli

import (
	"os"
	"os/exec"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/chigopher/pathlib"
)

// debug level logging of where we are running at the moment
func logwhere() {
	if isInContainer {
		if currentInstance == "" {
			zlog.Debug().Msgf("Running inside an unknown container") // TODO: read .version file or get from environment
		} else {
			zlog.Debug().Msgf("Running inside `%s`", currentInstance)
		}
	} else {
		if currentInstance == "" {
			zlog.Debug().Msgf("Running on host machine; no instance selected yet")
		} else {
			zlog.Debug().Msgf("Running on host machine; selected instance: %s", currentInstance)
		}
	}
	zlog.Debug().Msgf("Called as: %s", strings.Join(os.Args, " "))
}

// to rewrite the configuration file
func rewriteConfig() (err error) {
	if err = conf.WriteConfig(); err == nil {
		zboth.Debug().Msgf("Modified configuration file `%s`.", conf.ConfigFileUsed())
	} else {
		zboth.Warn().Err(err).Msgf("Failed to update the configuration file.")
	}
	return
}

// check if file exists, and is a file (keep it simple, runs before logging starts!
func existingFile(filePath string) (exists bool) {
	exists, _ = pathlib.NewPath(filePath).IsFile()
	return
}

// download a file, filepath is respective to current working directory
func downloadFile(fileURL string, downloadLocation string) (filepath pathlib.Path) {
	if resp, err := grab.Get(downloadLocation, fileURL); err == nil {
		zboth.Info().Msgf("Downloaded file saved as: %s", resp.Filename)
		filepath = *pathlib.NewPath(resp.Filename)
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to download file from: %s. Check log. ABORT!", fileURL)
	}
	return
}

// change directory with logging
func gotoFolder(givenName string) (pwd string) {
	var folder string
	if givenName == "workdir" {
		folder = "../.."
	} else {
		folder = workDir.Join(instancesWord, getInternalName(givenName)).String()
	}
	if err := os.Chdir(folder); err == nil {
		pwd, _ = os.Getwd()
		zboth.Debug().Msgf("Changed working directory to: %s", pwd)
	} else {
		zboth.Fatal().Msgf("Failed to changed working directory as required.")
	}
	return
}

// execute a command in shell
func execShell(command string) (result []byte, err error) {
	if result, err = exec.Command(shell, "-c", command).CombinedOutput(); err == nil {
		zboth.Debug().Msgf("Sucessfully executed shell command: %s", command)
	} else {
		zboth.Warn().Err(err).Msgf("Failed execution of command: %s", command)
	}
	return
}

// copy a text file
// func copyTextFile(source *pathlib.Path, target *pathlib.Path) (err error) {
// 	fmt.Println(source.String(), target.String())
// 	if reader, errRead := source.ReadFile(); err == nil {
// 		if errWrite := target.WriteFile(reader); err != nil {
// 			err = errWrite
// 		}
// 	} else {
// 		err = errRead
// 	}
// 	return
// }
