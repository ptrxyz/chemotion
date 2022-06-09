package cli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/chigopher/pathlib"
	"github.com/google/uuid"
	vercompare "github.com/hashicorp/go-version"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

var versionSuffix string = " --version"

// check if file exists, and is a file (keep it simple, runs before logging starts!
func existingFile(filePath string) (exists bool) {
	exists, _ = pathlib.NewPath(filePath).IsFile()
	return
}

// check if the CLI is running interactively; if not, then exit. Wrapper around currentState.quiet.
func confirmInteractive() {
	if currentState.quiet {
		zboth.Fatal().Err(fmt.Errorf("incomplete in quiet mode")).Msgf("%s is in quiet mode. Give all arguments to specify the desired action; use '--help' flag for more. ABORT!", nameCLI)
	}
	if currentState.isInside {
		zboth.Fatal().Err(fmt.Errorf("inside container in interactive mode")).Msgf("%s CLI is not meant to executed interactively from within a container. Use the `-q` flag. ABORT!", nameCLI)
	}
}

func confirmInstalled() {
	if firstRun {
		// Println output so that user is not discouraged by a FATAL error on-screen... especially when beginning with the tool.
		msg := fmt.Sprintf("Please install %s by running `%s` before using it.", nameCLI, "chemotion install")
		fmt.Println(msg)
		zlog.Fatal().Err(fmt.Errorf("chemotion not installed")).Msgf(msg)
	}
}

// check if a string is an array of strings, if yes, return the 1st index, else -1.
func stringInArray(str string, strings *[]string) int {
	for index, element := range *strings {
		if element == str {
			return index
		}
	}
	return -1
}

// check if a int is an array of int, if yes, return the 1st index, else -1.
func uint16InArray(num uint16, array *[]uint16) int {
	for index, element := range *array {
		if element == num {
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

// find version of a given software (using its command)
func findVersion(software string) (version string) {
	ver, err := execShell(software + versionSuffix)
	version = strings.TrimSpace(strings.Split(strings.TrimPrefix(strings.TrimPrefix(string(ver), "v"), "Docker version "), ",")[0]) // TODO: Regexify!
	if err != nil {
		zlog.Debug().Err(err).Msgf("Version determination of %s failed", software)
		if virtualizer == "Docker" && err.Error() == "exit status 1" {
			version = "Docker on WSL not running!"
		} else if err.Error() == "exit status 127" {
			version = "Unknown / not installed or found!" // 127 is software not found
		} else {
			version = err.Error()
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

// download a file, filepath is respective to current working directory
func downloadFile(fileURL string, downloadLocation string) (filepath pathlib.Path) {
	if resp, err := grab.Get(downloadLocation, fileURL); err == nil {
		zboth.Info().Msgf("Downloaded file saved as: %s", resp.Filename)
		filepath = *pathlib.NewPath(downloadLocation).Join(resp.Filename)
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to download file from: %s. Check log. ABORT!", fileURL)
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

// to manage config files as loaded into Viper
func getKeysValues(configuration *viper.Viper, key string) (keys, values []string) {
	for k, v := range configuration.GetStringMapString(key) {
		keys = append(keys, k)
		values = append(values, v)
	}
	return
}

// join keys so as to access them in a viper configuration
func joinKey(s ...string) (result string) {
	result = strings.Join(s, ".")
	return
}

// to lower case, same as strings.ToLower
var toLower = strings.ToLower

// to select an instance, gives a list to select from when less than 5, else a text input
func selectInstance(action string) (instance string) {
	existingInstances := allInstances()
	if len(existingInstances) < 5 {
		fmt.Printf("Please pick the instance to %s:\n", action)
		instance = selectOpt(existingInstances)
	} else {
		zlog.Debug().Msgf("String prompt to select instance")
		prompt := promptui.Prompt{
			Label: "Please name the instance to " + action,
			Validate: func(input string) (err error) {
				if input == "^C" {
					err = fmt.Errorf("^C")
				} else if len(input) == 0 {
					err = fmt.Errorf("can not accept empty value")
				} else {
					if stringInArray(input, &existingInstances) > -1 {
						err = nil
					} else {
						err = fmt.Errorf("no instance called %s", input)
					}
				}
				return
			},
		}
		if inst, err := prompt.Run(); err == nil {
			zlog.Debug().Msgf("Given answer: %s", inst)
			instance = inst
		} else if err.Error() == "^C" {
			zboth.Fatal().Err(fmt.Errorf("string prompt cancelled")).Msg("Input cancelled. Can't proceed without. Abort!")
		} else {
			zboth.Fatal().Err(err).Msgf("Prompt failed for unknown reason. Check log. Abort!")
		}
	}
	return
}

// to get all existing instances as determined by the configuration file
func allInstances() (instances []string) {
	instances, _ = getKeysValues(&conf, "instances")
	return
}

// to get all existing used ports
func allPorts() (ports []uint16) {
	existingInstances := allInstances()
	for _, instance := range existingInstances {
		ports = append(ports, uint16(conf.GetUint32(joinKey("instances", instance, "port"))))
	}
	return
}

//
func internalName(given_name string) (name string) {
	name = conf.GetString(joinKey("instances", given_name, "name"))
	return
}

// TODO

// Start shell for user
// var shellSystemRootCmd = &cobra.Command{
// 	Use:        "shell",
// 	SuggestFor: []string{"she"},
// 	Args:       cobra.NoArgs,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		logWhere()
// 		confirmInstalled()
// 		fmt.Println("We are now going to start shell")
// 		//TODO
// 	},
// }

// Start a rails shell for user
// var railsSystemRootCmd = &cobra.Command{
// 	Use:        "rails",
// 	SuggestFor: []string{"rai"},
// 	Args:       cobra.NoArgs,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		logWhere()
// 		confirmInstalled()
// 		fmt.Println("We are now going to start Rails shell")
// 		//TODO
// 	},
// }

// example starter

// var uninstallAdvancedRootCmd = &cobra.Command{
// 	Use:   "uninstall",
// 	Args:  cobra.NoArgs,
// 	Short: fmt.Sprintf("Uninstall %s completely.", nameCLI),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		logWhere()
// 		confirmInstalled()
// 		confirmInteractive()
// 	},
// }
