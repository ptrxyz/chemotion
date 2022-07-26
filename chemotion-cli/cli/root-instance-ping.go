package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	client = http.Client{
		Timeout: time.Second * 2,
	}
)

func instancePing(domain string) (code int, err error) {
	req, err := http.NewRequest("HEAD", domain, nil)
	if err == nil {
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			code = resp.StatusCode
		} else {
			code = 0
		}
	} else {
		code = 0
	}
	return
}

// sends the URL of
func getURL(givenName string) string {
	hostProtocol := conf.GetString(joinKey(instancesWord, givenName, "protocol"))
	hostAddress := conf.GetString(joinKey(instancesWord, givenName, "address"))
	hostPort := conf.GetString(joinKey(instancesWord, givenName, "port"))
	if hostProtocol == "" || hostAddress == "" || hostPort == "" {
		zboth.Fatal().Err(fmt.Errorf("key not found")).Msgf("Failed to find parts of URL for %s in %s.", givenName, conf.ConfigFileUsed())
	}
	return fmt.Sprintf("%s://%s:%s", hostProtocol, hostAddress, hostPort)
}

// PingCmd represents the ping command
var pingInstanceRootCmd = &cobra.Command{
	Use:   "ping",
	Args:  cobra.NoArgs,
	Short: "Ping an instance of " + nameCLI,
	Run: func(cmd *cobra.Command, _ []string) {
		var url string
		// Logic
		if ownCall(cmd) && cmd.Flag("url").Changed {
			url = cmd.Flag("url").Value.String()
			if err := addressValidate(url); err != nil {
				zboth.Fatal().Err(err).Msgf("Invalid URL. Please format your URL like this example: `https://chemotion.uni.de`")
			}
		} else {
			url = getURL(currentInstance)
		}
		zboth.Info().Msgf("Ping %s at %s.", currentInstance, url)
		if resp, err := instancePing(url); err == nil {
			zboth.Info().Msgf("Success, response %d received from %s.", resp, url)
		} else {
			zboth.Warn().Err(err).Msgf("Failure, response %d received from %s.", resp, url)
		}
	},
}

func init() {
	//rootCmd.AddCommand(pingCmd)
	instanceRootCmd.AddCommand(pingInstanceRootCmd)
	pingInstanceRootCmd.Flags().String("url", "", "The URL to ping")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
