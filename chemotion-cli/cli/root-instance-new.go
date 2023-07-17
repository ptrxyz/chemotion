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

// helper to get a compose file
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
		composeFilepath = downloadFile(use, pathlib.NewPath(".").Join(toSprintf("%s.%s", getNewUniqueID(), defaultComposeFilename)).String()) // downloads to where-ever it is called from
	} else {
		if isUrl {
			zboth.Fatal().Err(err).Msgf("Failed to download the file from URL: %s.", use)
		} else {
			zboth.Fatal().Err(err).Msgf("Failed %s for compose not found.", use)
		}
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

// helper to get a fresh (unassigned port)
func getFreshPort(kind string) (port uint64) {
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

// to create a development instance
func instanceCreateDevelopment(details map[string]string) (success bool) {
	zboth.Fatal().Err(toError("not implemented")).Msgf("This feature is currently under development.")
	return false
}

// interaction when creating a new instance
func processInstallAndInstanceCreateCmd(cmd *cobra.Command, details map[string]string) (create bool) {
	askName, askAddress, askDevelopment := true, true, true
	create = true
	details["givenName"] = instanceDefault
	details["accessAddress"] = addressDefault
	details["kind"] = "Production"
	details["use"] = getLatestComposeURL()
	if ownCall(cmd) {
		if cmd.Flag("name").Changed {
			details["givenName"] = cmd.Flag("name").Value.String()
			if err := newInstanceValidate(details["givenName"]); err != nil {
				zboth.Fatal().Err(err).Msgf("Cannot create new instance with name %s: %s", details["givenName"], err.Error())
			}
			askName = false
		}
		if cmd.Flag("address").Changed {
			details["accessAddress"] = cmd.Flag("address").Value.String()
			if err := addressValidate(details["accessAddress"]); err != nil {
				zboth.Fatal().Err(err).Msgf("Cannot accept the address %s: %s", details["accessAddress"], err.Error())
			}
			askAddress = false
		}
		if cmd.Flag("use").Changed {
			details["use"] = cmd.Flag("use").Value.String()
		}
		if cmd.Flag("development") != nil {
			if toBool(cmd.Flag("development").Value.String()) {
				details["kind"] = "Development"
			}
			askDevelopment = !cmd.Flag("use").Changed
		}
	}
	if isInteractive(false) {
		if firstRun || !ownCall(cmd) { // don't ask if the command is run directly i.e. without the menu
			{
				create = selectYesNo("Installation process may download containers (of multiple GBs) and can take some time. Continue", true)
			}
		}
		if create {
			if askName {
				details["givenName"] = getString("Please enter the name of the instance you want to create", newInstanceValidate)
			}
			if askAddress {
				if selectYesNo("Is this instance having its own web-address (e.g. https://chemotion.uni.de or http://chemotion.uni.de:4100)?", false) {
					details["accessAddress"] = getString("Please enter the web-address", addressValidate)
				}
			}
			if askDevelopment && !firstRun {
				if !selectYesNo("Do you want a Production instance", true) {
					details["kind"] = "Development"
				}
			}
		}
	}
	// create new unique name for the instance
	details["name"] = toSprintf("%s-%s", details["givenName"], getNewUniqueID())
	return
}

func createExtendedCompose(name, use string) (extendedCompose viper.Viper) {
	extendedCompose = *viper.New()
	compose := parseCompose(use)
	sections := []string{"services", "volumes", "networks"}
	// set labels on services, volumes and networks for future identification
	for _, section := range sections {
		subheadings := getSubHeadings(&compose, section) // subheadings are the names of the services, volumes and networks
		for _, k := range subheadings {
			extendedCompose.Set(joinKey(section, k, "labels"), map[string]string{"net.chemotion.cli.project": name})
		}
	}
	// set unique name for volumes in the compose file
	volumes := getSubHeadings(&compose, "volumes")
	for _, volume := range volumes {
		n := compose.GetString(joinKey("volumes", volume, "name"))
		if n == "" && volume == "spectra" {
			n = "chemotion_spectra"
		} // because the spectra volume has no name
		if strings.HasPrefix(n, name) { // for compatibility with upgradeThisTool("0.1_to_0.2")
			extendedCompose.Set(joinKey("volumes", volume, "name"), n)
		} else {
			extendedCompose.Set(joinKey("volumes", volume, "name"), name+"_"+n)
		}
	}
	return
}

func instanceCreateProduction(details map[string]string) (success bool) {
	pro, add, port := splitAddress(details["accessAddress"])
	details["protocol"], details["address"] = pro, add
	if port == 0 {
		port = getFreshPort(details["kind"])
		if details["address"] == "localhost" {
			details["accessAddress"] += toSprintf(":%d", port)
		}
	} else {
		if details["address"] == "localhost" {
			zboth.Warn().Err(toError("localhost && port suggested")).Msgf("You suggested a port while running on localhost. We strongly recommend that you use the default schema i.e. do not assign a specific port.")
			if isInteractive(false) {
				if !selectYesNo("Continue still", false) {
					zboth.Info().Msgf("Operation cancelled")
					os.Exit(2)
				}
			}
		}
	}
	details["port"] = strconv.FormatUint(port, 10)
	// download and modify the compose file
	var composeFile pathlib.Path
	if existingFile(details["use"]) {
		if err := copyfile(details["use"], toSprintf("%s.%s", getNewUniqueID(), defaultComposeFilename)); err == nil {
			composeFile = *pathlib.NewPath(toSprintf("%s.%s", getNewUniqueID(), defaultComposeFilename))
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to copy the suggested compose file: %s. This is necessary for future use.")
		}
	} else {
		composeFile = downloadFile(details["use"], toSprintf("%s.%s", getNewUniqueID(), defaultComposeFilename))
	}
	if err := changeKey(composeFile.String(), joinKey("services", "eln", "ports[0]"), toSprintf("%s:4000", details["port"])); err != nil {
		zboth.Fatal().Err(err).Msgf("Failed to update the downloaded compose file. This is necessary for future use.")
	}
	extendedCompose := createExtendedCompose(details["name"], composeFile.String())
	// store values in the conf, the conf file is modified only later
	if firstRun {
		conf.SetConfigFile(workDir.Join(configFile).String())
		conf.Set("version", versionYAML)
		conf.Set(joinKey(stateWord, selectorWord), details["givenName"])
		conf.Set(joinKey(stateWord, "quiet"), false)
		conf.Set(joinKey(stateWord, "debug"), false)
		conf.Set(joinKey(stateWord, "version"), versionCLI)
	}
	conf.Set(joinKey(instancesWord, details["givenName"], "port"), port)
	for _, key := range []string{"name", "kind", "accessAddress"} {
		conf.Set(joinKey(instancesWord, details["givenName"], key), details[key])
	}
	// make folder and move the compose file into it
	zboth.Info().Msgf("Creating a new instance of %s called %s.", nameCLI, details["name"])
	if err := workDir.Join(instancesWord, details["name"]).MkdirAll(); err == nil {
		composeFile.Rename(workDir.Join(instancesWord, details["name"], defaultComposeFilename))
	} else {
		zboth.Fatal().Err(err).Msgf("Unable to create folder to store instances of %s.", nameCLI)
	}
	// write out the extended compose file
	if _, err, _ := gotoFolder(details["givenName"]), extendedCompose.WriteConfigAs(extenedComposeFilename), gotoFolder("workdir"); err == nil {
		zboth.Info().Msgf("Written compose files %s and %s in the above steps.", defaultComposeFilename, extenedComposeFilename)
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to write the extended compose file to its repective folder. This is necessary for future use.")
	}
	if _, success, _ = gotoFolder(details["givenName"]), callVirtualizer(composeCall+"up --no-start"), gotoFolder("workdir"); !success {
		zboth.Fatal().Err(toError("compose up failed")).Msgf("Failed to setup %s. Check log. ABORT!", nameCLI)
	}
	// initialize the env file
	conf.Set(joinKey(instancesWord, details["givenName"], "environment", "URL_HOST"), strings.TrimPrefix(details["accessAddress"], pro+"://"))
	conf.Set(joinKey(instancesWord, details["givenName"], "environment", "URL_PROTOCOL"), pro)
	if err := rewriteConfig(); err != nil {
		zboth.Fatal().Err(err).Msg("Failed to write config file. Check log. ABORT!") // we want a fatal error in this case, `rewriteConfig()` does a Warn error
	}
	return success
}

// command to install a new instance of Chemotion
var newInstanceRootCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.NoArgs,
	Short: "Create a new instance of " + nameCLI,
	Run: func(cmd *cobra.Command, _ []string) {
		details := make(map[string]string)
		create := processInstallAndInstanceCreateCmd(cmd, details)
		if create {
			switch details["kind"] {
			case "Production":
				if success := instanceCreateProduction(details); success {
					zboth.Info().Msgf("Successfully created a new production instance. Once switched on, it can be found at: %s", details["accessAddress"])
				}
			case "Development":
				if success := instanceCreateDevelopment(details); success {
					zboth.Info().Msgf("Successfully created a new development instance.")
				}
			}
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(newInstanceRootCmd)
	newInstanceRootCmd.Flags().StringP("name", "n", instanceDefault, "Name for the new instance")
	newInstanceRootCmd.Flags().String("use", "", "URL or filepath of the compose file to use for creating the instance")
	newInstanceRootCmd.Flags().String("address", addressDefault, "Web-address (or hostname) for accessing the instance")
	newInstanceRootCmd.Flags().Bool("development", false, "Create a development instance")
}
