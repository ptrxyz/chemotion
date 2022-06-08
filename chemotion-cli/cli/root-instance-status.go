package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var _chemotion_instance_status_all_ bool

func instanceStatus(given_name string) (status string) {
	name := internalName(given_name)
	if res, err := execShell(fmt.Sprintf("docker ps -a --filter \"label=net.chemotion.cli.project=%s\" --format \"{{.Status}}\"", name)); err == nil {
		out := strings.Split(string(res), "\n")
		statuses := []string{}
		for _, line := range out { // determine what are the status messages for all associated containers
			status := strings.Split(line, " ")[0]
			if stringInArray(status, &statuses) == -1 && len(status) != 0 {
				statuses = append(statuses, status)
			}
		}
		status = statuses[0]
		if len(statuses) > 1 {
			status = strings.Join(statuses, " and ")
		}
	} else {
		zboth.Fatal().Err(err).Msgf("Failed to get status of the instance called %s", currentState.name)
		status = ""
	}
	return
}

var statusInstanceRootCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status of an instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		if !_chemotion_instance_status_all_ {
			status := instanceStatus(currentState.name)
			zboth.Info().Msgf("The status of %s is: %s.", currentState.name, status)
		} else {
			for _, instance := range allInstances() {
				status := instanceStatus(instance)
				zboth.Info().Msgf("The status of %s is: %s.", instance, status)
			}
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(statusInstanceRootCmd)
	statusInstanceRootCmd.Flags().BoolVar(&_chemotion_instance_status_all_, "all", false, "show status of all instances")
}
