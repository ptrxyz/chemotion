package cmd

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
	if stringInArray("--debug", &os.Args) > 0 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	currentState.Quiet = false
	if stringInArray("-q", &os.Args) > 0 || stringInArray("--quiet", &os.Args) > 0 {
		currentState.Quiet = true
	}
	// start logging
	if logFile, err := workDir.Join(logFileName).OpenFile(os.O_APPEND | os.O_CREATE | os.O_WRONLY); err == nil {
		zlog = zerolog.New(logFile).With().Timestamp().Logger()
		if currentState.Quiet {
			zboth = zlog // in this case, both the loggers point to the same file and there should be no console output
		} else {
			console := zerolog.ConsoleWriter{Out: os.Stdout}
			console.FormatErrFieldName = func(i interface{}) string { return "" }  // we don't want error to be shown in the console
			console.FormatErrFieldValue = func(i interface{}) string { return "" } // PartsExclude doesn't seem to work!
			multi := zerolog.MultiLevelWriter(logFile, console)
			zboth = zerolog.New(multi).With().Timestamp().Logger()
		}
		zlog.Debug().Msg("Successfully initialized logging.")
		logRunningOn()
	} else {
		minimalConsoleWriter := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		minimalConsoleWriter.Fatal().Err(err).Msg("Can't write log file. ABORT!") // minimal console writer
	}
}

// helpers: logging

// debug level logging of where we are running at the moment
func logRunningOn() {
	if currentState.isInside {
		zlog.Debug().Msgf("Running inside container `%s`.", currentState.name)
	} else {
		zlog.Debug().Msgf("Running `%s` on a host machine.", currentState.name)
	}
}
