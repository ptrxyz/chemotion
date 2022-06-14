package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var _root_instance_log_all_ bool
var _root_instance_log_service_ string
var _root_instance_log_rollNum_ = 1            // the index number assigned by virtualizer
var _root_instance_log_defaultService_ = "eln" // the service that is used as the default for printing log"
var _root_instance_log_v_details_ bool
var _root_instance_log_v_follow_ bool
var _root_instance_log_v_timestamps_ bool
var _root_instance_log_v_since_ string
var _root_instance_log_v_until_ string
var _root_instance_log_v_tail_ string

func instanceLog(givenName, service string) {
	name := internalName(givenName)
	out := getColumn(givenName, "Names")
	var services, logOf []string
	for _, line := range out { // determine what are the status messages for all associated containers
		l := strings.TrimSpace(line) // use only the first word
		if len(l) > 0 {
			l = strings.TrimPrefix(l, fmt.Sprintf("%s-", name))
			l = strings.TrimSuffix(l, fmt.Sprintf("-%d", _root_instance_log_rollNum_))
			services = append(services, l)
		}
	}
	if _root_instance_log_all_ {
		logOf = services
	} else {
		if stringInArray(service, &services) > -1 {
			logOf = []string{service}
		} else {
			zboth.Warn().Err(fmt.Errorf("named service not found")).Msgf("No service called %s found associated with the instance called %s.", service, givenName)
		}
	}
	for _, service = range logOf {
		zboth.Info().Msgf("Printing logs for the instance-service called %s-%s.", givenName, service)
		args := fmt.Sprintf("--tail %s", _root_instance_log_v_tail_)
		if _root_instance_log_v_details_ {
			args += " --details"
		}
		if _root_instance_log_v_timestamps_ {
			args += " --timestamps"
		}
		if _root_instance_log_v_since_ != "" {
			args += fmt.Sprintf("--since %s", _root_instance_log_v_until_)
		}
		if _root_instance_log_v_until_ != "" {
			args += fmt.Sprintf("--until %s", _root_instance_log_v_until_)
		}
		if _root_instance_log_v_follow_ {
			args += " --follow"
			callVirtualizer(fmt.Sprintf("logs %s %s-%s-%d", args, name, service, _root_instance_log_rollNum_))
		} else {
			if res, err := execShell(fmt.Sprintf("%s logs %s %s-%s-%d", toLower(virtualizer), args, name, service, _root_instance_log_rollNum_)); err == nil {
				if n, errPrint := fmt.Println(string(res)); errPrint == nil {
					zboth.Debug().Msgf("Printed logs to screen that were %d lines long", n)
				} else {
					zboth.Warn().Err(errPrint).Msgf("Error while printing logs for the instance-container called %s-%s-%d.", name, service, _root_instance_log_rollNum_)
				}
			} else {
				zboth.Fatal().Err(err).Msgf("Failed to get logs for the instance-container called %s-%s-%d.", name, service, _root_instance_log_rollNum_)
			}
		}
	}
}

var logInstanceRootCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"logs"},
	Args:    cobra.NoArgs,
	Short:   "Get logs of an instance of " + nameCLI,
	Run: func(cmd *cobra.Command, args []string) {
		logWhere()
		confirmInstalled()
		if currentState.quiet {
			zboth.Warn().Err(fmt.Errorf("illegal operation")).Msgf("Logs can't be printed in quiet mode.")
		} else {
			if _root_instance_log_service_ == "" {
				_root_instance_log_all_ = true
			}
			_root_instance_log_service_ = toLower(_root_instance_log_service_)
			instanceLog(currentState.name, _root_instance_log_service_)
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(logInstanceRootCmd)
	logInstanceRootCmd.Flags().StringVar(&_root_instance_log_service_, "service", _root_instance_log_defaultService_, "show the log of a given service")
	logInstanceRootCmd.Flags().BoolVar(&_root_instance_log_all_, "all", false, "show the logs of all services")
	logInstanceRootCmd.Flags().BoolVar(&_root_instance_log_v_details_, "details", false, fmt.Sprintf("--details flag as received by %s logs command", virtualizer))
	logInstanceRootCmd.Flags().BoolVar(&_root_instance_log_v_follow_, "follow", false, fmt.Sprintf("--follow flag as received by %s logs command", virtualizer))
	logInstanceRootCmd.Flags().BoolVarP(&_root_instance_log_v_timestamps_, "timestamps", "t", false, fmt.Sprintf("--timestamps flag as received by %s logs command", virtualizer))
	logInstanceRootCmd.Flags().StringVar(&_root_instance_log_v_since_, "since", "", fmt.Sprintf("--since flag as received by %s logs command", virtualizer))
	logInstanceRootCmd.Flags().StringVar(&_root_instance_log_v_until_, "until", "", fmt.Sprintf("--until flag as received by %s logs command", virtualizer))
	logInstanceRootCmd.Flags().StringVarP(&_root_instance_log_v_tail_, "tail", "n", "all", fmt.Sprintf("--tail flag as received by %s logs command", virtualizer))
}