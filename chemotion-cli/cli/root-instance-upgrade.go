package cli

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func pullImages(use string) {
	tempCompose := parseCompose(use)
	services := getSubHeadings(&tempCompose, "services")
	if len(services) == 0 {
		zboth.Warn().Err(toError("no services found")).Msgf("Please check that %s is a valid compose file with named services.", tempCompose.ConfigFileUsed())
	}
	for _, service := range services {
		zboth.Info().Msgf("Pulling image for the service called %s", service)
		if success := callVirtualizer(toSprintf("pull %s", tempCompose.GetString(joinKey("services", service, "image")))); !success {
			zboth.Warn().Err(toError("pull failed")).Msgf("Failed to pull image for the service called %s", service)
		}
	}
}

func instanceUpgrade(givenName, use string) {
	// first read the new compose file
	tempCompose := parseCompose(use)
	name := getInternalName(givenName)
	// back up old compose file
	oldComposeFile := workDir.Join(instancesWord, name, defaultComposeFilename)
	oldCompose := parseCompose(oldComposeFile.String())
	// for some reason (no idea why), labels must be set before port
	setUniqueLabels(&tempCompose, name)
	tempCompose.Set(joinKey("services", "eln", "ports"), oldCompose.GetStringSlice(joinKey("services", "eln", "ports"))) // copy over the ports listing
	// shutdown existing instance's docker
	success := true
	if _, worked, _ := gotoFolder(givenName), callVirtualizer("compose down --remove-orphans"), gotoFolder("workdir"); !worked {
		success = worked
		zboth.Fatal().Err(toError("compose down failed")).Msgf("Failed to stop %s. Check log. ABORT!", givenName)
	}
	if success {
		if _, worked, _ := gotoFolder(givenName), callVirtualizer(toSprintf("volume rm %s_chemotion_app", name)), gotoFolder("workdir"); !worked {
			success = worked
			zboth.Fatal().Err(toError("volume removal failed")).Msgf("Failed to remove old app volume. Check log. ABORT!")
		}
	}
	if success {
		oldComposeFile.Rename(workDir.Join(instancesWord, name, strconv.FormatInt(time.Now().Unix(), 10)+".old."+defaultComposeFilename))
		if err := tempCompose.WriteConfigAs(workDir.Join(instancesWord, name, defaultComposeFilename).String()); err == nil {
			commandStr := toSprintf("compose -f %s up --no-start", defaultComposeFilename)
			zboth.Info().Msgf("Starting %s with command: %s", virtualizer, commandStr)
			if _, worked, _ := gotoFolder(givenName), callVirtualizer(commandStr), gotoFolder("workdir"); !worked {
				zboth.Fatal().Err(toError("%s failed", commandStr)).Msgf("Failed to initialize upgraded %s. Check log. ABORT!", givenName)
			}
		} else {
			zboth.Fatal().Err(toError("compose file write fail")).Msgf("Failed to write the new compose file. The old one is still available as %s", oldComposeFile.Name())
		}
	}
}

func getLatestComposeURL() (url string) {
	var err error
	if url, err = getLatestReleaseURL(); err == nil {
		url = strings.Join([]string{url, defaultComposeFilename}, "/")
	} else {
		zboth.Warn().Err(err).Msgf("Could not determine the address of the latest compose file, using this one: %s.", composeURL)
		url = composeURL
	}
	return
}

var upgradeInstanceRootCmd = &cobra.Command{
	Use:   "upgrade",
	Args:  cobra.NoArgs,
	Short: "Upgrade (the selected) instance of " + nameCLI,
	Run: func(cmd *cobra.Command, _ []string) {
		var pull, backup, upgrade bool = false, false, true
		var use string = ""
		if ownCall(cmd) {
			if cmd.Flag("use").Changed {
				use = cmd.Flag("use").Value.String()
			}
			pull = toBool(cmd.Flag("pull-only").Value.String())
			upgrade = !pull
		}
		if !pull && isInteractive(false) {
			switch selectOpt([]string{"all actions: pull image, backup and upgrade", "preparation: pull image and backup", "upgrade only (if already prepared)", "pull image only", "exit"}, "What do you want to do") {
			case "all actions: pull image, backup and upgrade":
				pull, backup, upgrade = true, true, true
			case "preparation: pull image and backup":
				pull, backup, upgrade = true, true, false
			case "upgrade only (if already prepared)":
				pull, backup, upgrade = false, false, true
			case "pull image only":
				pull, backup, upgrade = true, false, false
			}
		}
		if use == "" {
			use = getLatestComposeURL()
		}
		if pull {
			pullImages(use)
		}
		if backup {
			instanceBackup(currentInstance, "both")
		}
		if upgrade {
			if instanceStatus(currentInstance) == "Up" {
				zboth.Fatal().Err(toError("upgrade fail; instance is up")).Msgf("Cannot upgrade an instance that is currently running. Please turn it off before continuing.")
			}
			instanceUpgrade(currentInstance, use)
		}
	},
}

func init() {
	upgradeInstanceRootCmd.Flags().String("use", composeURL, "URL or filepath of the compose file to use for upgrading")
	upgradeInstanceRootCmd.Flags().Bool("pull-only", false, "Pull image for use in upgrade, don't do the upgrade")
	instanceRootCmd.AddCommand(upgradeInstanceRootCmd)
}
