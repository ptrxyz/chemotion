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
	if conf.GetString((joinKey("instances", name))) == "" {
		zboth.Warn().Msgf("An instance with name %s already exists", name)
		return false
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
	zboth.Info().Msgf("Creating a new instance of %s called %s.", nameCLI, name)
	// make folder
	if err := workDir.Join(instancesFolder, name).MkdirAll(); err != nil {
		zboth.Fatal().Err(err).Msgf("Unable to create folder to store instances of %s.", nameCLI)
	}
	if isUrl { // move file if required
		if err := composeFilepath.Rename(workDir.Join(instancesFolder, name, composeFilepath.Name())); err != nil {
			zboth.Fatal().Err(err).Msgf("Unable to move file to its respective folder. This is necessary for future use.")
		}
	} else { // copy the file into the folder for future use
		zlog.Info().Msgf("Copying specified file to `%s` folder for future use.", instancesFolder)
		if err := copyTextFile(&composeFilepath, workDir.Join(instancesFolder, name, composeFilepath.Name())); err != nil {
			zboth.Fatal().Err(err).Msgf("Unable to copy file to its respective folder. This is necessary for future use.")
		}
	}
	allServices := getValueInViper(&composeConf, "services")
	for _, k := range allServices {
		setValueInViper(&composeConf, k+".name", name, "prefix")
	}
	os.Chdir("instances/" + name)
	zlog.Debug().Msgf("Changed working directory to: instances/%s", name)
	commandStr := fmt.Sprintf("compose -f %s create", composeFilepath.Name())
	zboth.Info().Msgf("Starting %s with command: %s", toLower(virtualizer), commandStr)
	if success = callVirtualizer(commandStr); !success {
		zboth.Fatal().Err(fmt.Errorf("%s failed", commandStr)).Msgf("Failed to setup %s. Check log. ABORT!", nameCLI)
	}
	os.Chdir("../..")
	zboth.Info().Msg("Successfully created container the container.")
	if firstRun {
		conf.SetConfigFile(workDir.Join(defaultConfigFilepath).String())
		conf.Set("version", versionYAML)
		conf.Set(selector_key, given_name)
	}
	conf.Set(joinKey("instances", given_name, "name"), name)
	conf.Set(joinKey("instances", given_name, "kind"), kind)
	conf.Set(joinKey("instances", given_name, "quiet"), false)
	conf.Set(joinKey("instances", given_name, "debug"), kind == "Development")
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
