package cli

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/chigopher/pathlib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// helper function to get a compose file
func parseCompose(use string) (compose viper.Viper) {
	var (
		composeFilepath pathlib.Path
		isUrl           bool
	)
	// TODO: check on the version of the compose file
	if existingFile(use) {
		composeFilepath = *pathlib.NewPath(use)
	} else if _, err := url.ParseRequestURI(use); err == nil {
		isUrl = true
		composeFilepath = downloadFile(use, workDir.String()) // downloads to the working directory
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to parse the URL/file: %s.", use)
	}
	// parse the compose file
	compose = *viper.New()
	compose.SetConfigFile(composeFilepath.String())
	err := compose.ReadInConfig()
	if isUrl {
		composeFilepath.Remove()
	}
	if err != nil {
		zboth.Fatal().Err(err).Msgf("Invalid formatting for a compose file.")
	}
	return
}

// set unique labels and volume names in the compose file
func setUniqueLabels(compose *viper.Viper, name string) {
	sections := []string{"services", "volumes", "networks"}
	for _, section := range sections {
		subheadings := getSubHeadings(compose, section) // subheadings are the names of the services, volumes and networks
		for _, k := range subheadings {
			compose.Set(joinKey(section, k, "labels"), map[string]string{"net.chemotion.cli.project": name})
		}
	}
	// set unique name for volumes in the compose file
	volumes := getSubHeadings(compose, "volumes")
	for _, volume := range volumes {
		n := compose.GetString(joinKey("volumes", volume, "name"))
		compose.Set(joinKey("volumes", volume, "name"), name+"_"+n)

	}
}

// helper function to get a fresh (unassigned port)
func getFreshPort(kind string) (port uint) {
	if firstRun {
		port = firstPort
	} else {
		existingPorts := allPorts()
		if kind == "Production" {
			for i := firstPort + 101; i <= maxInstancesOfKind+(firstPort+101); i++ {
				if elementInSlice(i, &existingPorts) == -1 {
					port = i
					break
				}
			}
		} else if kind == "Development" {
			for i := firstPort + 201; i <= maxInstancesOfKind+(firstPort+201); i++ {
				if elementInSlice(i, &existingPorts) == -1 {
					port = i
					break
				}
			}
		}
		if port == (firstPort+101)+maxInstancesOfKind || port == (firstPort+201)+maxInstancesOfKind {
			zboth.Fatal().Err(toError("max instances")).Msgf("A maximum of %d instances of %s are allowed. Please contact us if you hit this limit.", maxInstancesOfKind, nameCLI)
		}
	}
	return
}

func readEnv(filepath string) (env *viper.Viper) {
	env = viper.New()
	env.SetConfigType("env")
	if filepath == "" {
		env.SetConfigFile(".env")
	} else {
		env.SetConfigFile(filepath)
		if err := env.ReadInConfig(); err != nil {
			zboth.Fatal().Err(err).Msgf("Failed to parse the supplied .env file.")
		}
	}
	return
}

func instanceCreateDevelopment(cmd *cobra.Command) (success bool) {
	zboth.Fatal().Err(toError("not implemented")).Msgf("This feature is currently under development.")
	if err := newInstanceValidate(cmd.Flag("name").Value.String()); err != nil {
		zboth.Fatal().Err(err).Msgf("Given instance name is invalid because %s.", err.Error())
	}
	return false
}

func instanceCreateProduction(cmd *cobra.Command) (success bool) {
	if err := newInstanceValidate(cmd.Flag("name").Value.String()); err != nil {
		zboth.Fatal().Err(err).Msgf("Given instance name is invalid because %s.", err.Error())
	}
	var (
		port uint
		// expose   uint
		protocol string
		address  string
		env      *viper.Viper
	)
	givenAddress := cmd.Flag("address").Value.String()
	givenName := cmd.Flag("name").Value.String()
	env = readEnv(cmd.Flag("env").Value.String())
	if env.InConfig("URL_PROTOCOL") && env.InConfig("URL_HOST") {
		if cmd.Flag("address").Changed {
			zboth.Warn().Msgf("It seems you have `address` set in .env file as well as via the --address flag. The value in the given .env file will be overwritten.")
		} else {
			givenAddress = env.GetString("URL_PROTOCOL") + "://" + env.GetString("URL_HOST")
		}
	}
	protocol, address, port = splitAddress(givenAddress)
	if address != "localhost" && (protocol == "http" && port == 443) || (protocol == "https" && port == 80) {
		zboth.Warn().Err(toError("port mismatch")).Msgf("You have chosen port %d for protocol %s. This is generally a very bad idea.", port, protocol)
		if isInteractive(false) {
			if !selectYesNo("Continue still", false) {
				zboth.Info().Msgf("Operation cancelled")
				os.Exit(2)
			}
		}
	}
	if port == 0 { // i.e. a port was not suggested by the user
		if address == "localhost" {
			port = getFreshPort("Production")
			givenAddress += ":" + strconv.Itoa(int(port))
		} else {
			if protocol == "http" {
				port = 80
			} else {
				port = 443
			}
		}
	} else {
		if address == "localhost" {
			zboth.Warn().Err(toError("localhost && port suggested")).Msgf("You suggested a port while running on localhost. We strongly recommend that you use the default schema i.e. do not assign a specific port.")
			if isInteractive(false) {
				if !selectYesNo("Continue still", false) {
					zboth.Info().Msgf("Operation cancelled")
					os.Exit(2)
				}
			}
		}
	}
	// create new unique name for the instance
	name := toSprintf("%s-%s", givenName, getNewUniqueID())
	// store values in the conf, the conf file is modified only later
	if firstRun {
		conf.SetConfigFile(workDir.Join(defaultConfigFilepath).String())
		conf.Set("version", versionYAML)
		conf.Set(joinKey(stateWord, selectorWord), givenName)
		conf.Set(joinKey(stateWord, "quiet"), false)
		conf.Set(joinKey(stateWord, "debug"), false)
	}
	conf.Set(joinKey(instancesWord, givenName, "name"), name)
	conf.Set(joinKey(instancesWord, givenName, "kind"), "Production")
	conf.Set(joinKey(instancesWord, givenName, "protocol"), protocol)
	conf.Set(joinKey(instancesWord, givenName, "address"), address)
	conf.Set(joinKey(instancesWord, givenName, "port"), port)
	// get the compose file for the instance
	compose := parseCompose(cmd.Flag("use").Value.String())
	// for some reason (no idea why), labels must be set before port
	setUniqueLabels(&compose, name)
	// set the port in the compose file
	compose.Set(joinKey("services", "eln", "ports"), []string{toSprintf("%d:4000", port)})
	zboth.Info().Msgf("Creating a new instance of %s called %s.", nameCLI, name)
	// make folder
	if err := workDir.Join(instancesWord, name).MkdirAll(); err != nil {
		zboth.Fatal().Err(err).Msgf("Unable to create folder to store instances of %s.", nameCLI)
	}
	if _, err, _ := gotoFolder(givenName), compose.WriteConfigAs(defaultComposeFilename), gotoFolder("workdir"); err == nil {
		zboth.Info().Msgf("Written compose file %s in the above step.", compose.ConfigFileUsed())
		commandStr := toSprintf("compose -f %s up --no-start", defaultComposeFilename)
		zboth.Info().Msgf("Starting %s with command: %s", virtualizer, commandStr)
		if _, worked, _ := gotoFolder(givenName), callVirtualizer(commandStr), gotoFolder("workdir"); !worked {
			success = worked
			zboth.Fatal().Err(toError("%s failed", commandStr)).Msgf("Failed to setup %s. Check log. ABORT!", nameCLI)
		}
	} else {
		success = false
		zboth.Fatal().Err(err).Msgf("Failed to write the compose file to its repective folder. This is necessary for future use.")
	}
	// write env file into the container
	envFile := workDir.Join(instancesWord, name, ".env")
	env.SetConfigFile(envFile.String())
	env.Set("URL_HOST", strings.TrimPrefix(givenAddress, protocol+"://"))
	env.Set("URL_PROTOCOL", protocol)
	if err := env.WriteConfig(); err == nil {
		modifyContainer(givenName, "mkdir -p", "shared/pullin", "")
		if worked := modifyContainer(givenName, "cp", ".env", "shared/pullin/."); !worked {
			success = worked
			zboth.Warn().Msgf("Failed to write .env file in `%s/shared/pullin`", name)
		}
	} else {
		zboth.Warn().Err(err).Msgf("Failed to write .env file")
	}
	envFile.Remove()
	zboth.Info().Msgf("Successfully created the instance called %s. New %s port available at %d.", givenName, nameCLI, port)
	// now modify the config file
	if err := rewriteConfig(); err != nil {
		zboth.Fatal().Err(err).Msg("Failed to write config file. Check log. ABORT!")
	}
	return success
}

func newInstanceInteraction(cmd *cobra.Command) (create bool) {
	create = true
	if firstRun || !ownCall(cmd) { // don't ask if the command is run directly i.e. without the menu
		create = selectYesNo("Installation process may download containers (of multiple GBs) and can take some time. Continue", true)
	}
	if create {
		if ownCall(cmd) && !cmd.Flag("name").Changed { // i.e user has not changed it by passing an argument
			if err := cmd.Flag("name").Value.Set(getString("Please enter the name of the instance you want to create", newInstanceValidate)); err != nil {
				zboth.Warn().Err(err).Msgf("Failed to allocate given value. It will be ignored.")
			}
		}
		if ownCall(cmd) && (!cmd.Flag("env").Changed && !cmd.Flag("address").Changed) { // i.e user has not changed it by passing an argument
			if selectYesNo("Is this instance running on a web-server?", false) {
				if err := cmd.Flag("address").Value.Set(getString("Please enter the web-address e.g. https://chemotion.uni.de:125", addressValidate)); err != nil {
					zboth.Warn().Err(err).Msgf("Failed to allocate given value. It will be ignored.")
				}
			}
		}
	} else {
		zboth.Info().Msgf("Installation cancelled.")
	}
	return
}

// command to install a new container of Chemotion
var newInstanceRootCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.NoArgs,
	Short: "Create a new instance of " + nameCLI,
	Run: func(cmd *cobra.Command, _ []string) {
		create := true
		if isInteractive(true) {
			create = newInstanceInteraction(cmd)
			if create && ownCall(cmd) && !cmd.Flag("development").Changed { // i.e. the flag was not set
				cmd.Flag("development").Value.Set(strconv.FormatBool(!selectYesNo("Do you want a Production instance", true)))
			}
		}
		if create {
			if toBool(cmd.Flag("development").Value.String()) {
				if success := instanceCreateDevelopment(cmd); success {
					zboth.Info().Msg("Successfully created a new development instance.")
				}
			} else {
				if success := instanceCreateProduction(cmd); success {
					zboth.Info().Msg("Successfully created a new production instance.")
				}
			}
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(newInstanceRootCmd)
	newInstanceRootCmd.Flags().String("name", instanceDefault, "Name for the new instance")
	newInstanceRootCmd.Flags().String("use", composeURL, "URL or filepath of the compose file to use for creating the instance")
	newInstanceRootCmd.Flags().String("address", addressDefault, "Web-address (or hostname) for accessing the instance")
	newInstanceRootCmd.Flags().String("env", "", ".env file for the new instance")
	newInstanceRootCmd.Flags().Uint("expose", 0, "port that is exposed by the instance to access it")
	newInstanceRootCmd.Flags().Bool("development", false, "Create a development instance")
}
