package cmd

import (
	"github.com/cavaliergopher/grab/v3"
	"github.com/rs/zerolog"

	"github.com/spf13/viper"
)

const defaultConfig = "chemotion-cli.yml"

type Config struct {
	// NB: only variables that begin with a capital letter are unmarshalled
	// see https://stackoverflow.com/a/64919032
	Instance string `mapstructure:"instance"`
	Quiet    bool
	Debug    bool
}

var (
	configFile string
	configURL  string = "https://github.com/harivyasi/chemotion/tree/chemotion-cli/chemotion-cli/configs/" + defaultConfig
	config     Config
)

// Viper is used to load values from config file. Cobra is the basis of our command line interface.
// This function uses Viper to set flags on Cobra.
// (See how cool this sounds, make sure you pick fun project names!)
func initViper() {
	// initialize viper and check status of the config flag
	zlog.Debug().Str("container", containerName).Str("instance", config.Instance).Msg("Attempting to read configuration file.")
	filePath := chemotionCmd.Flag("config").Value.String()
	// specified: if a particular filepath was specified by using the flag i.e. changed by the user
	specified := chemotionCmd.Flag("config").Changed
	// determine how to go ahead depending on the status of the flag
	var kind string // one of the following: specified, existing, downloaded
	if specified {
		kind = "specified"
		viper.SetConfigFile(filePath) // set by file path
	} else {
		if existingFile(filePath) {
			kind = "existing"
		} else {
			kind = "downloaded"
			// this function executes before the flags are read in from the terminal
			// i.e. it ignores the `quiet` flag and logs out everything to terminal
			zboth.Warn().Msg("Could not find the configuration file. Downloading from repository.")
			resp, errDownload := grab.Get(".", configURL)
			if errDownload == nil {
				filePath = resp.Filename
				zlog.Info().Str("container", containerName).Str("instance", config.Instance).Msg("Downloaded configuration file saved as: " + configFile)
			} else {
				zboth.Fatal().Err(errDownload).Msg("Failed to download configuration file from " + configURL + ". ABORT!")
			}
		}
		viper.SetConfigName(filePath) // set by file name
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	// Try and read the configuration file, then unmarshal it
	if err := viper.ReadInConfig(); err == nil {
		if errUnmarshal := viper.Unmarshal(&config); err == nil {
			zlog.Info().Str("container", containerName).Str("instance", config.Instance).Msg("Read " + kind + " configuration file: " + viper.ConfigFileUsed())
			if !config.Debug {
				logLevel = zerolog.InfoLevel
			}
		} else {
			zboth.Fatal().Err(errUnmarshal).Msg("Failed to map values from configuration file. ABORT!")
		}
	} else {
		zboth.Fatal().Err(err).Msg("Failed to read " + kind + " configuration file. ABORT!")
	}
}
