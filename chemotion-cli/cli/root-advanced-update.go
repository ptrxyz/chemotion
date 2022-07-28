package cli

import (
	"net/http"
	"strings"

	"github.com/chigopher/pathlib"
	"github.com/spf13/cobra"
)

func getLatestReleaseURL() (url string, err error) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if resp, errGet := client.Get(releaseUnresolvedURL); errGet == nil {
		if loc, errLoc := resp.Location(); errLoc == nil {
			url = loc.String()
		} else {
			err = errLoc
		}
	} else {
		err = errGet
	}
	return strings.Replace(url, "tag", "download", -1), err
}

func selfUpdate() {
	var url string
	if u, err := getLatestReleaseURL(); err != nil {
		zboth.Fatal().Err(err).Msgf("Could not determine address of latest executable.")
		url = u
	}
	oldVersion := pathlib.NewPath(commandForCLI)
	stat, _ := oldVersion.Stat()
	cliFileName := oldVersion.Name()
	url = toSprintf("%s/%s", url, cliFileName)
	newVersion := downloadFile(url, workDir.Join(toSprintf("%s.new", cliFileName)).String())
	if err := newVersion.Chmod(stat.Mode() | 100); err != nil { // make sure that it remains executable for the ErrUseLastResponse
		zboth.Warn().Err(err).Msgf("Could not grant executable permission to the downloaded file. Please do it yourself.")
	}
	if errOld := oldVersion.RenameStr(toSprintf("%s.old", cliFileName)); errOld == nil {
		if errNew := newVersion.RenameStr(cliFileName); errNew == nil {
			zboth.Info().Msgf("Successfully downloaded the new version. Old version is available as %s and is safe to remove.", oldVersion.Name())
		} else {
			zboth.Warn().Err(errNew).Msgf("Successfully downloaded the new version. Please rename it to %s for further use. The old version is available as %s and is safe to remove.", cliFileName, oldVersion.Name())
		}
	} else {
		zboth.Warn().Err(errOld).Msgf("Successfully downloaded the new version but failed to rename the old one. The new version is called %s, please rename it %s. The old version is safe to remove.", newVersion.Name(), cliFileName)
	}

}

var updateSelfAdvancedRootCmd = &cobra.Command{
	Use:   "update",
	Short: "Update this tool itself",
	Run: func(cmd *cobra.Command, args []string) {
		selfUpdate()
	},
}

func init() {
	advancedRootCmd.AddCommand(updateSelfAdvancedRootCmd)
}
