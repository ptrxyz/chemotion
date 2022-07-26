package cli

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func instanceUpgrade(givenName, use string) {
	// first read the new compose file
	tempCompose := getCompose(use)
	name := getInternalName(givenName)
	// back up old compose file
	oldComposeFile := workDir.Join(instancesWord, name, defaultComposeFilename)
	oldCompose := getCompose(oldComposeFile.String())
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

var upgradeInstanceRootCmd = &cobra.Command{
	Use:   "upgrade",
	Args:  cobra.NoArgs,
	Short: "Upgrade (the selected) instance of " + nameCLI,
	Run: func(cmd *cobra.Command, _ []string) {
		if instanceStatus(currentInstance) == "Up" {
			zboth.Fatal().Err(toError("upgrade fail; instance is up")).Msgf("Cannot upgrade an instance that is currently running. Please turn it off before continuing.")
		} else {
			upgrade, use := true, composeURL
			if isInteractive(false) {
				upgrade = selectYesNo("Please be sure to backup before proceeding. Continue", false)
			}
			if ownCall(cmd) {
				use = cmd.Flag("use").Value.String()
			}
			if upgrade {
				instanceUpgrade(currentInstance, use)
			}
		}
	},
}

func init() {
	upgradeInstanceRootCmd.Flags().String("use", composeURL, "URL or filepath of the compose file to use for upgrading")
	instanceRootCmd.AddCommand(upgradeInstanceRootCmd)
}
