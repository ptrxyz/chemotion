package cli

import (
	"net/http"

	"github.com/spf13/cobra"
)

func instancePing(domain string) (code int, err error) {
	var resp *http.Response
	if resp, err = http.Get(domain); err == nil {
		code = resp.StatusCode
	} else {
		code = 0
	}
	return
}

// sends the URL of a given chemotion instance
func getURL(givenName string) string {
	hostProtocol := conf.GetString(joinKey(instancesWord, givenName, "protocol"))
	hostAddress := conf.GetString(joinKey(instancesWord, givenName, "address"))
	hostPort := conf.GetString(joinKey(instancesWord, givenName, "port"))
	if hostProtocol == "" || hostAddress == "" || hostPort == "" {
		zboth.Fatal().Err(toError("key not found")).Msgf("Failed to find parts of URL for %s in %s.", givenName, conf.ConfigFileUsed())
	}
	return toSprintf("%s://%s:%s", hostProtocol, hostAddress, hostPort)
}

// PingCmd represents the ping command
var pingInstanceRootCmd = &cobra.Command{
	Use:   "ping",
	Args:  cobra.NoArgs,
	Short: "Ping an instance of " + nameCLI,
	Run: func(cmd *cobra.Command, _ []string) {
		url := getURL(currentInstance)
		zboth.Info().Msgf("Ping %s at %s.", currentInstance, url)
		if resp, err := instancePing(url); err == nil {
			zboth.Info().Msgf("Success, response %d received from %s.", resp, url)
		} else {
			zboth.Warn().Err(err).Msgf("Failure, response %d received from %s.", resp, url)
		}
	},
}

func init() {
	instanceRootCmd.AddCommand(pingInstanceRootCmd)
}
