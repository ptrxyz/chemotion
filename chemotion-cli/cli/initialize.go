package cli

import (
	"os"

	"github.com/rs/zerolog"
)

// Initializes logging. Ignores values in the configuration as configuration is loaded after this initialization.
func initLog() {
	// lowest level reading of the debug and quiet flags
	// alas, it works only with command line flags, otherwise
	// we have to wait for the values to be read in from the config file
	// this low-level reading has to be done because logging begins before reading the config file.
	if zerolog.SetGlobalLevel(zerolog.InfoLevel); elementInSlice("--debug", &os.Args) > 0 || elementInSlice("-d", &os.Args) > 0 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// start logging
	if logFile, err := workDir.Join(logFilename).OpenFile(os.O_APPEND | os.O_CREATE | os.O_WRONLY); err == nil {
		zlog = zerolog.New(logFile).With().Timestamp().Logger()
		if elementInSlice("-q", &os.Args) > 0 || elementInSlice("--quiet", &os.Args) > 0 {
			zboth = zlog // in this case, both the loggers point to the same file and there should be no console output
		} else {
			console := zerolog.ConsoleWriter{Out: os.Stdout}
			console.FormatErrFieldName = func(_ any) string { return "" }  // we don't want error to be shown in the console
			console.FormatErrFieldValue = func(_ any) string { return "" } // PartsExclude doesn't seem to work!
			multi := zerolog.MultiLevelWriter(logFile, console)
			zboth = zerolog.New(multi).With().Timestamp().Logger()
		}
		zlog.Debug().Msgf("%s started. Successfully initialized logging", nameCLI)
	} else {
		minimalConsoleWriter := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		minimalConsoleWriter.Fatal().Err(err).Msg("Can't write log file. ABORT!") // minimal console writer
	}
}

func initFlags() {
	zlog.Debug().Msg("Start: initialize flags")
	// flag 1: instance, i.e. name of the instance to operate upon
	// terminal overrides config-file, default is read from the config file
	rootCmd.PersistentFlags().StringVarP(&currentInstance, "selected-instance", "i", "", toSprintf("select an existing instance of %s when starting", nameCLI))
	// flag 2: config, the configuration file
	// config as a flag cannot be read from the configuration file because that creates a circular dependency, default name is hard-coded
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", defaultConfigFilepath, "path to the configuration file")
	// flag 3: quiet, i.e. should the CLI run in interactive mode
	// terminal overrides config-file, default is false
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, toSprintf("use %s in scripted mode i.e. without an interactive prompt", commandForCLI))
	// flag 4: debug, i.e. should debug messages be logged
	// terminal overrides config-file, default is false
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable logging of debug messages")
	zlog.Debug().Msg("End: initialize flags")
}

// Viper is used to load values from config file. Cobra is the basis of our command line interface.
// This function uses Viper to set flags on Cobra.
// (See how cool this sounds, make sure you pick fun project names!)
func initConf() {
	zlog.Debug().Msg("Start: initialize configuration")
	// check status of the config flag
	// if changed and specified file is not found then exit
	// otherwise set the value of configFile on viper
	// then use the path as determined by viper and set as the value of configFile
	if rootCmd.Flag("config-file").Changed && !existingFile(configFile) {
		// here configFile should be same as rootCmd.Flag("config-file").Value.String()
		zboth.Fatal().Err(toError("specified config file not found")).Msgf("Please ensure that the file you specify using --config/-f flag does exist.")
	}
	conf.SetConfigFile(configFile)
	zlog.Debug().Msg("Attempting to read configuration file")
	// if the flag is not changed, check for the posibility of first run
	if existingFile(configFile) {
		firstRun = false
		// Try and read the configuration file, then unmarshal it
		if err := conf.ReadInConfig(); err == nil {
			if conf.IsSet(selectorWord) && conf.IsSet(instancesWord) {
				if errUnmarshal := conf.UnmarshalKey(selectorWord, &currentInstance); errUnmarshal != nil {
					zboth.Fatal().Err(toError("unmarshal failed")).Msgf("Failed to unmarshal the mandatory key %s in the file: %s.", selectorWord, configFile)
				}
				if !conf.IsSet(joinKey(instancesWord, currentInstance)) {
					zboth.Fatal().Err(toError("unmarshal failed")).Msgf("Failed to find the description for instance `%s` in the file: %s.", currentInstance, configFile)
				}
			} else {
				zboth.Fatal().Err(toError("unmarshal failed")).Msgf("Failed to find the mandatory keys `%s` and `%s` in the file: %s.", selectorWord, instancesWord, configFile)
			}
		} else {
			zboth.Fatal().Err(err).Msgf("Failed to read configuration file: %s. ABORT!", configFile)
		}
	}
	zlog.Debug().Msgf("End: initialize configuration; Config found?: %t; is inside container?: %t", !firstRun, isInContainer)
}

// bind the command line flags to the configuration
func bindFlags() {
	zlog.Debug().Msg("Start: bind flags")
	if err := conf.BindPFlag(selectorWord, rootCmd.Flag("selected-instance")); err != nil {
		zboth.Warn().Err(err).Msgf("Failed to bind flag: %s. Will ignore command line input.", "selected-instance")
	}
	for _, flag := range []string{"debug", "quiet"} {
		if err := conf.BindPFlag(joinKey(stateWord, flag), rootCmd.Flag(flag)); err != nil {
			zboth.Warn().Err(err).Msgf("Failed to bind flag: %s. Will ignore command line input.", flag)
		}
		if !conf.IsSet(joinKey(stateWord, flag)) {
			conf.Set(joinKey(stateWord, flag), false)
		}
	}
	zlog.Debug().Msg("End: bind flags")
}
