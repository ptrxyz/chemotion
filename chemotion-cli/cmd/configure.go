package cmd

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Viper is used to load values from config file. Cobra is the basis of our command line interface.
// This function uses Viper to set flags on Cobra.
// (See how cool this sounds, make sure you pick fun project names!)
func initViper() {
	zlog.Debug().Msg("Start initViper()")
	// initialize viper and check status of the config flag
	zlog.Debug().Msg("Attempting to read configuration file.")
	if existingFile(cmd.Flag("config").Value.String()) {
		viper.SetConfigFile(cmd.Flag("config").Value.String()) // set file path to what is given by user, or the default flag value
	}
	if cmd.Flag("config").Changed && !existingFile(viper.ConfigFileUsed()) {
		zboth.Fatal().Err(fmt.Errorf("specified config file not found")).Msgf("Please ensure that you specify in --config flag does exist.")
	}
	if existingFile(viper.ConfigFileUsed()) {
		firstRun = false
		// Try and read the configuration file, then unmarshal it
		if err := viper.ReadInConfig(); err == nil {
			if err = viper.UnmarshalKey("selected", &currentState.name); err != nil {
				zboth.Fatal().Err(fmt.Errorf("unmarshal failed")).Msgf("Failed to find the mandatory key `instanceDefault` in the file %s.", viper.ConfigFileUsed())
			}
			if errUnmarshal := viper.UnmarshalKey("instances."+currentState.name, &currentState); err == nil {
				zboth.Info().Msgf("Read configuration file: %s.", viper.ConfigFileUsed())
				if currentState.Debug {
					zerolog.SetGlobalLevel(zerolog.DebugLevel) // escalate the debug level if said so by the config file
				} // don't do else because flags have the final say!
			} else {
				zboth.Fatal().Err(errUnmarshal).Msg("Failed to map values from configuration file. ABORT!")
			}
			if !viper.IsSet("instances." + currentState.name) {
				zboth.Fatal().Err(fmt.Errorf("unmarshal failed")).Msgf("Failed to find the values for `%s` instance in the file %s.", currentState.name, viper.ConfigFileUsed())
			}
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to read %s configuration file. ABORT!", configFile)
		}
	}
	zlog.Debug().Msgf("End: initViper(), Config found?: %t, is Inside?: %t", !firstRun, currentState.isInside)
}
