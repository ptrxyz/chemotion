package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
)

const (
	logFileName = "chemotion-cli.log"
)

var (
	zlog     zerolog.Logger
	zboth    zerolog.Logger
	logLevel zerolog.Level = zerolog.DebugLevel // start in debug, eventually value from user
	fs       afero.Fs      = afero.NewOsFs()    // pointer to the filesystem
)

// Initializes logging. Ignores configured values as loading of configuration is done after this initialization.
func initLog() {
	zerolog.SetGlobalLevel(logLevel)
	var (
		logFile afero.File
		err     error
	)
	if logFile, err = fs.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		minimalConsoleWriter := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
		minimalConsoleWriter.Fatal().Err(err).Msg("Can't write log file. ABORT!") // minimal console writer
	}
	zlog = zerolog.New(logFile).With().Timestamp().Logger()
	if config.Quiet {
		zboth = zlog
	} else {
		console := zerolog.ConsoleWriter{Out: os.Stdout}
		console.FormatErrFieldName = func(i interface{}) string { return "" }  // we don't want error to be shown in the console
		console.FormatErrFieldValue = func(i interface{}) string { return "" } // PartsExclude doesn't seem to work!
		multi := zerolog.MultiLevelWriter(logFile, console)
		zboth = zerolog.New(multi).With().Timestamp().Logger()
	}
	zlog.Debug().Str("container", containerName).Str("instance", config.Instance).Msg("Successfully initialized logging.")
}
