package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Creates a default Production instance, with set defaults.
func firstRun() {
	confirmVirtualizer(minimumVirtualizer)
	zboth.Info().Msgf("Welcome to your first run of %s. We need to download and install it first.\nThis operation may download containers (~%d GB) and can take some time.", nameCLI, containersSize)
	if !currentState.Quiet {
		if !selectYesNo("Do you want to continue", false) {
			zboth.Info().Msgf("Operation cancelled. %s will exit gracefully.", nameCLI)
			os.Exit(0)
		}
	} else {
		zboth.Warn().Msgf("You chose do first run of chemotion in quiet mode. Will go ahead and install it!")
	}
	success := instanceCreate(defaultInstanceName, "Production", composeURL)
	if success {
		zboth.Info().Msg("Successfully created container for the first run.")
		if err := viper.WriteConfigAs(defaultConfigFileName); err == nil {
			zboth.Info().Msgf("Written config file: %s.", defaultConfigFileName)
		} else {
			zboth.Fatal().Err(err).Msg("Failed to write config file. Check log. ABORT!")
		}
		zboth.Info().Msgf("All done! Now you can do `%s on` and `%s off` to start/stop %s.", chemotionCmd.Use, chemotionCmd.Use, nameCLI)
		zlog.Debug().Msgf("Exiting %s gracefully", nameCLI)
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
			zboth.Fatal().Err(fmt.Errorf("unmarshal failed")).Msgf("Failed to find the mandatory key `chosen` in the file %s.", viper.ConfigFileUsed())
		}
		if errUnmarshal := viper.UnmarshalKey("instances."+currentState.name, &currentState); err == nil {
			zboth.Info().Msgf("Read %s configuration file: %s.", found, viper.ConfigFileUsed())
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
		zboth.Fatal().Err(err).Msgf("Failed to read %s configuration file. ABORT!", found)
	}
	zlog.Debug().Msg("End: initViper()")
}
