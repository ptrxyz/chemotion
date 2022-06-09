package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/chigopher/pathlib"
	"github.com/spf13/cobra"
)

var _chemotion_instance_new_name_ string
var _chemotion_instance_new_use_ string
var _chemotion_instance_new_development_ bool

func instanceCreate(name string, kind string, use string) (success bool) {
	var port uint16
	if firstRun {
		port = 4000
	} else {
		existingInstances := allInstances()
		if stringInArray(name, &existingInstances) > -1 {
			zboth.Fatal().Err(fmt.Errorf("instance %s already exists", name)).Msgf("An instance with name %s already exists.", name)
			return false
		}
		existingPorts := allPorts()
		if kind == "Production" {
			for i := uint16(4100); i <= maxInstancesOfKind+4100; i++ {
				if uint16InArray(i, &existingPorts) == -1 {
					port = i
					break
				}
			}
		} else if kind == "Development" {
			for i := uint16(4200); i <= maxInstancesOfKind+4200; i++ {
				if uint16InArray(i, &existingPorts) == -1 {
					port = i
					break
				}
			}
		}
		if port == 4100+maxInstancesOfKind || port == 4200+maxInstancesOfKind {
			zboth.Fatal().Err(fmt.Errorf("max instances")).Msgf("A maximum of %d instances of %s are allowed. Please contact us if you hit this limit.", maxInstancesOfKind, nameCLI)
		}
	}
	confirmVirtualizer(minimumVirtualizer)
	given_name := name
	name = fmt.Sprintf("%s-%s", name, getNewUniqueID())
	var composeFilepath pathlib.Path // TODO: check on the version of the compose file
	var isUrl bool = false
	if existingFile(use) {
		composeFilepath = *pathlib.NewPath(use)
	} else if _, err := url.ParseRequestURI(use); err == nil {
		isUrl = true
		composeFilepath = downloadFile(use, workDir.String()) // downloads to the working directory
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to parse the URL/file: %s.", use)
	}
	compose.SetConfigFile(composeFilepath.String())
	if err := compose.ReadInConfig(); err == nil {
		if isUrl {
			composeFilepath.Remove()
		}
	} else {
		if isUrl {
			composeFilepath.Remove()
		}
		zboth.Fatal().Err(err).Msgf("Invalid formatting for a compose file.")
	}
	// set labels in the docker-compose file
	sections := []string{"services", "volumes", "networks"}
	for _, section := range sections {
		subheadings, _ := getKeysValues(&compose, section) // subheadings are the names of the services, volumes and networks
		for _, k := range subheadings {
			compose.Set(joinKey(section, k, "labels"), map[string]string{"net.chemotion.cli.project": name})
		}
	}
	// set unique name for volumes
	volumes, _ := getKeysValues(&compose, "volumes")
	for _, volume := range volumes {
		n := compose.GetString(joinKey("volumes", volume, "name"))
		compose.Set(joinKey("volumes", volume, "name"), name+"_"+n)
	}
	// set the port
	compose.Set(joinKey("services", "eln", "ports"), []string{fmt.Sprintf("%d:4000", port)})
	zboth.Info().Msgf("Creating a new instance of %s called %s.", nameCLI, name)
	// make folder
	if err := workDir.Join(instancesFolder, name).MkdirAll(); err != nil {
		zboth.Fatal().Err(err).Msgf("Unable to create folder to store instances of %s.", nameCLI)
	}
	os.Chdir("instances/" + name)
	zlog.Debug().Msgf("Changed working directory to: instances/%s", name)
	if err := compose.WriteConfigAs(composeFilename); err != nil {
		os.Chdir("../..")
		zboth.Fatal().Err(err).Msgf("Failed to write the compose file to its repective folder. This is necessary for future use.")
	}
	commandStr := fmt.Sprintf("compose -f %s up --no-start", composeFilename)
	zboth.Info().Msgf("Starting %s with command: %s", toLower(virtualizer), commandStr)
	if success = callVirtualizer(commandStr); !success {
		os.Chdir("../..")
		zboth.Fatal().Err(fmt.Errorf("%s failed", commandStr)).Msgf("Failed to setup %s. Check log. ABORT!", nameCLI)
	}
	os.Chdir("../..")
	zboth.Info().Msgf("Successfully created container the container. New %s port available at %d.", nameCLI, port)
	if firstRun {
		conf.SetConfigFile(workDir.Join(defaultConfigFilepath).String())
		conf.Set("version", versionYAML)
		conf.Set(selector_key, given_name)
	}
	conf.Set(joinKey("instances", given_name, "name"), name)
	conf.Set(joinKey("instances", given_name, "kind"), kind)
	conf.Set(joinKey("instances", given_name, "quiet"), false)
	conf.Set(joinKey("instances", given_name, "debug"), kind == "Development")
	conf.Set(joinKey("instances", given_name, "port"), port)
	if err := conf.WriteConfig(); err == nil {
		zboth.Info().Msgf("Written config file: %s.", conf.ConfigFileUsed())
	} else {
		zboth.Fatal().Err(err).Msg("Failed to write config file. Check log. ABORT!")
	}
	return success
}

// command to install a new container of Chemotion
var newInstanceRootCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.NoArgs,
	Short: "Create a new instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		create := true
		kind := "Production"
		if _chemotion_instance_new_development_ {
			kind = "Development"
		}
		if !currentState.quiet {
			confirmInteractive()
			if selectYesNo("Installation process may download containers (of multiple GBs) and can take some time. Continue", false) {
				if !_chemotion_instance_new_development_ { // i.e. the flag was not set
					fmt.Println("What kind of instance do you want?")
					kind = selectOpt([]string{"Production", "Development"})
				}
				if _chemotion_instance_new_name_ == instanceDefault {
					_chemotion_instance_new_name_ = getString("Please enter name of the instance you want to create")
				}
			} else {
				create = false
				zboth.Info().Msgf("Installation cancelled.")
			}
		}
		if create {
			if success := instanceCreate(_chemotion_instance_new_name_, kind, _chemotion_instance_new_use_); success {
				zboth.Info().Msg("Successfully created the new instance")
			}
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(newInstanceRootCmd)
	newInstanceRootCmd.Flags().StringVar(&_chemotion_instance_new_name_, "name", instanceDefault, "name for the new instance")
	newInstanceRootCmd.Flags().StringVar(&_chemotion_instance_new_use_, "use", composeURL, "URL or filepath to use for creating the instance")
	newInstanceRootCmd.Flags().BoolVar(&_chemotion_instance_new_development_, "development", false, "create a development instance")
}
