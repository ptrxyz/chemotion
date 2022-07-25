package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// check if the CLI is running interactively; if no and fail, then exit. Wrapper around conf.GetBool(joinKey(stateWord,"quiet")).
func isInteractive(fail bool) (interactive bool) {
	interactive = true
	if conf.GetBool(joinKey(stateWord, "quiet")) {
		if interactive = false; fail {
			zboth.Fatal().Err(fmt.Errorf("incomplete in quiet mode")).Msgf("%s is in quiet mode. Give all arguments to specify the desired action; use '--help' flag for more. ABORT!", nameCLI)
		}
	}
	if isInContainer {
		if interactive = false; fail {
			zboth.Fatal().Err(fmt.Errorf("inside container in interactive mode")).Msgf("%s CLI is not meant to executed interactively from within a container. Use the `-q` flag. ABORT!", nameCLI)
		}
	}
	return
}

// check if an element is in an array of type(element), if yes, return the 1st index, else -1.
func elementInSlice[T uint | int | float64 | string](elem T, slice *[]T) int {
	for index, element := range *slice {
		if element == elem {
			return index
		}
	}
	return -1
}

// remove element from slice, CAUTION: order not preserved
func removeElementInSlice[T uint | int | float64 | string](index int, slice []T) []T {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

// generate a new UID (of the form xxxxxxxx) as a string
func getNewUniqueID() string {
	id, _ := uuid.NewRandom()
	return strings.Split(id.String(), "-")[0]
}

// to manage config files as loaded into Viper
func getSubHeadings(configuration *viper.Viper, key string) (subheadings []string) {
	for k := range configuration.GetStringMapString(key) {
		subheadings = append(subheadings, k)
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

// toBool
func toBool(s string) (value bool) {
	if toLower(s) == "true" {
		value = true
	} else if toLower(s) == "false" {
		value = false
	} else {
		err := fmt.Errorf("cannot convert %s to bool", s)
		zboth.Fatal().Err(err).Msgf(err.Error())
	}
	return
}

// determine if the command was called on its own (true) or access via a menu (false)
func ownCall(cmd *cobra.Command) bool {
	return len(cmd.Commands()) == 0 // a command is accessed on its own if there are no child commands
}

// to get all existing instances as determined by the configuration file
func allInstances() (instances []string) {
	instances = getSubHeadings(&conf, instancesWord)
	return
}

// to get all existing used ports
func allPorts() (ports []uint) {
	existingInstances := allInstances()
	for _, instance := range existingInstances {
		ports = append(ports, uint(conf.GetUint32(joinKey(instancesWord, instance, "port"))))
	}
	return
}

// get internal name for an instance
func getInternalName(givenName string) (name string) {
	name = conf.GetString(joinKey(instancesWord, givenName, "name"))
	if name == "" {
		zboth.Fatal().Err(fmt.Errorf("instance not found")).Msgf("No such instance: %s", givenName)
	}
	return
}

// get column associated with `ps` output for a given instance of chemotion
func getColumn(givenName, column string) (values []string) {
	name := getInternalName(givenName)
	if res, err := execShell(fmt.Sprintf("%s ps -a --filter \"label=net.chemotion.cli.project=%s\" --format \"{{.%s}}\"", toLower(virtualizer), name, column)); err == nil {
		values = strings.Split(string(res), "\n")
	} else {
		values = []string{}
	}
	return
}

// get services associated with a given `instance` of Chemotion
func getServices(givenName string) (services []string) {
	name, out := getInternalName(givenName), getColumn(givenName, "Names")
	for _, line := range out { // determine what are the status messages for all associated containers
		l := strings.TrimSpace(line) // use only the first word
		if len(l) > 0 {
			l = strings.TrimPrefix(l, fmt.Sprintf("%s-", name))
			l = strings.TrimSuffix(l, fmt.Sprintf("-%d", rollNum))
			services = append(services, l)
		}
	}
	return
}

// split address into subcomponents
func splitAddress(full string) (protocol string, address string, port uint) {
	if err := addressValidate(full); err != nil {
		zboth.Fatal().Err(err).Msgf("Given address %s is invalid.", full)
	}
	protocol, address, _ = strings.Cut(full, ":")
	address = strings.TrimPrefix(address, "//")
	address, portStr, _ := strings.Cut(address, ":")
	if port = 0; portStr != "" {
		p, _ := strconv.Atoi(portStr)
		port = uint(p)
	}
	return
}

// TODO

// Start shell for user
// var shellSystemRootCmd = &cobra.Command{
// 	Use:        "shell",
// 	SuggestFor: []string{"she"},
// 	Args:       cobra.NoArgs,
// 	Run: func(cmd *cobra.Command, args []string) {
//
//
// 		fmt.Println("We are now going to start shell")
// 	},
// }

// Start a rails shell for user
// var railsSystemRootCmd = &cobra.Command{
// 	Use:        "rails",
// 	SuggestFor: []string{"rai"},
// 	Args:       cobra.NoArgs,
// 	Run: func(cmd *cobra.Command, args []string) {
//
//
// 		fmt.Println("We are now going to start Rails shell")
// 	},
// }

// example starter

// var uninstallAdvancedRootCmd = &cobra.Command{
// 	Use:   "uninstall",
// 	Args:  cobra.NoArgs,
// 	Short: fmt.Sprintf("Uninstall %s completely.", nameCLI),
// 	Run: func(cmd *cobra.Command, args []string) {
//
//
// 		confirmInteractive()
// 	},
// }
