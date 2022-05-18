package cmd

import (
	"os"

	"github.com/rs/zerolog"
)

// Initializes logging. Ignores configured values as loading of configuration is done after this initialization.
func initLog() {
	if stringInArray("--debug", os.Args) > 0 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// lowest level reading of the debug flag
	// alas, it works only with explicit flagging, otherwise
	// we have to wait for the flag to be read in from the config file
	if logFile, err := fs.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePerm); err == nil {
		zlog = zerolog.New(logFile).With().Timestamp().Logger()
		if currentState.Quiet {
			zboth = zlog // in this case, both the loggers point to the same file and there is no console output
		} else {
			console := zerolog.ConsoleWriter{Out: os.Stdout}
			console.FormatErrFieldName = func(i interface{}) string { return "" }  // we don't want error to be shown in the console
			console.FormatErrFieldValue = func(i interface{}) string { return "" } // PartsExclude doesn't seem to work!
			multi := zerolog.MultiLevelWriter(logFile, console)
			zboth = zerolog.New(multi).With().Timestamp().Logger()
		}
		zlog.Info().Msg("Successfully initialized logging.")
		logRunningOn()
	} else {
		minimalConsoleWriter := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		minimalConsoleWriter.Fatal().Err(err).Msg("Can't write log file. ABORT!") // minimal console writer
	}
}

// logging helpers

// log (debug level) where we are running at the moment
func logRunningOn() {
	if currentState.isInside {
		zlog.Debug().Msg("Running inside a container.")
	} else {
		zlog.Debug().Msg("Running on a host machine.")
	}
}
