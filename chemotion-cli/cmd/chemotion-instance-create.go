package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/chigopher/pathlib"
	"github.com/spf13/cobra"
)

var _chemotion_instance_create_development_ bool
var _chemotion_instance_create_use_ string

func instanceCreate(name string, kind string, use string) (success bool) {
	minimumVirtualizer := "17.12" // TODO: set via the compose
	confirmVirtualizer(minimumVirtualizer)
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
	os.Chdir("instances/" + name)
	zlog.Debug().Msgf("Changed working directory to: instances/%s", name)
	commandStr := fmt.Sprintf("compose -f %s create", composeFilepath.Name())
	zboth.Info().Msgf("Starting %s with command: %s", virtualizer, commandStr)
	if success = callVirtualizer(commandStr); !success {
		zboth.Fatal().Err(fmt.Errorf("%s failed", commandStr)).Msgf("Failed to setup %s. Check log. ABORT!", nameCLI)
	}
	os.Chdir("../..")
	// TODO: change the chemotion-cli.yml file
	return success
}

var createInstance = &cobra.Command{
	Use:  "create <name_of_instance>",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var kind string = "Production"
		if _chemotion_instance_create_development_ {
			kind = "Development"
		}
		var actOn string
		if len(args) < 1 {
			confirmInteractive()
			actOn = getString("Please enter name of the instance you want to create")
		} else {
			actOn = args[0]
		}
		zboth.Info().Msgf("We are now going to create an instance called %s.", actOn)
		instanceCreate(actOn, kind, _chemotion_instance_create_use_)
	},
}

func init() {
	instanceCmd.AddCommand(createInstance)
	createInstance.Flags().BoolVar(&_chemotion_instance_create_development_, "development", false, "create a development instance")
	createInstance.Flags().StringVar(&_chemotion_instance_create_use_, "use", composeURL, "URL or filepath to use for creating the instance")
}
