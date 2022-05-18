package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// to create a default Production instance
func firstRun() {
	zboth.Info().Msgf("Welcome to your first run of %s.", projectName)
	currentState.Kind = "Production"
	success := instanceCreate(currentState)
	if success {
		zboth.Info().Msg("Successfully created container for the first run.")
		if err := viper.WriteConfigAs(defaultConfigFileName); err == nil {
			zboth.Info().Msgf("Written config file: %s", defaultConfigFileName)
		} else {
			zboth.Fatal().Err(err).Msg("Failed to write config file. Check log. ABORT!")
		}
		zboth.Info().Msgf("All done! Now you can do `%s on` and `%s off` to start/stop %s.", chemotionCmd.Use, chemotionCmd.Use, projectName)
		zlog.Info().Msgf("Exiting %s gracefully", projectName)
		os.Exit(0)
	}
}

// Viper is used to load values from config file. Cobra is the basis of our command line interface.
// This function uses Viper to set flags on Cobra.
// (See how cool this sounds, make sure you pick fun project names!)
func initViper() {
	zlog.Debug().Msg("Start initViper()")
	// initialize viper and check status of the config flag
	zlog.Debug().Msg("Attempting to read configuration file.")
	// determine how to go ahead depending on the status of the flag
	var found string // one of the following: specified, existing, created
	if chemotionCmd.Flag("config").Changed {
		found = "specified"
		viper.SetConfigFile(chemotionCmd.Flag("config").Value.String()) // set file path to what is given by user
	} else {
		if existingFile(defaultConfigFileName) {
			found = "existing"
		} else {
			found = "created"
			firstRun()
		}
		viper.SetConfigFile(defaultConfigFileName)
	}
	// Try and read the configuration file, then unmarshal it
	if err := viper.ReadInConfig(); err == nil {
		if err = viper.UnmarshalKey("chosen", &currentState.name); err != nil {
			zboth.Fatal().Err(fmt.Errorf("unmarshal failed")).Msgf("Failed to find the key `chosen` in the file %s.", viper.ConfigFileUsed())
		}
		if !viper.IsSet(currentState.name) {
			zboth.Fatal().Err(fmt.Errorf("unmarshal failed")).Msgf("Failed to find the values for `%s` instance in the file %s.", currentState.name, viper.ConfigFileUsed())
		}
		if errUnmarshal := viper.UnmarshalKey(currentState.name, &currentState); err == nil {
			zlog.Info().Str("instance", currentState.name).Msg("Read " + found + " configuration file: " + viper.ConfigFileUsed())
			if currentState.Debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel) // escalate the debug level if said so by the config file
			} // don't do else because flags have the final say!
		} else {
			zboth.Fatal().Err(errUnmarshal).Msg("Failed to map values from configuration file. ABORT!")
		}
	} else {
		zboth.Fatal().Err(err).Msg("Failed to read " + found + " configuration file. ABORT!")
	}
	zlog.Debug().Msg("End: initViper()")
}
